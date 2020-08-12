/*
Package lz77 implements LZ77 encoding and decoding.

Encode outputs data in blocks. A block starts with an 8-bit header and is
followed by at most eight data units. At the end of the data stream there may be
less than 8 units following a header. In this case the bits in the header
without corresponding data units are meaningless.

Each bit in the header specifies the type of the corresponding unit following
the header. A 0-bit means the corresponding unit is a literal byte. A 1-bit
means the corresponding unit is a reference to a previous location in the data
stream.

A reference is a pair (l, d) where l is the length of the referred byte sequence
and d is the starting point of the sequence as an offset from the current
position. References are encoded as unsigned 16-bit little-endian integers. The
length part of the integer is 4 bits and the distance part takes the remaining
12 bits.
*/
package lz77

import (
	"io"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bits"
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
)

// These constants specify how the 16 bits of a reference are distributed
// between length and distance.
const (
	refLenBits  = 4
	refDistBits = 12
)

// These constants specify the sizes of the lookahead and window buffers.
const (
	lookaheadBufferSize = (1 << refLenBits) - 1
	windowBufferSize    = (1 << refDistBits) - 1
)

// Encode reads data from input, encodes it using LZ77 and writes the result to
// output.
func Encode(input io.Reader, output io.Writer) error {
	src := bufio.NewReaderSize(input, lookaheadBufferSize)
	dst := bits.NewWriter(output)
	window := newWindowBuffer(windowBufferSize)
	headerBuf := make([]byte, 1)
	units := make([]uint16, 0, 8)

	for {
		headerBuf[0] = 0
		unitHeader := bits.NewList(headerBuf)
		units = units[:0]
		for i := 0; i < cap(units); i++ {
			lookahead, err := src.Peek(lookaheadBufferSize)
			if err != nil {
				if err != io.EOF {
					return err
				}
				if len(lookahead) == 0 {
					break
				}
			}
			ref := window.findLongestPrefix(lookahead)
			if ref.length == 0 {
				unitHeader.Set(i, false)
				next, err := src.ReadByte()
				if err != nil {
					return err
				}
				units = append(units, uint16(next))
				window.appendByte(next)
			} else {
				unitHeader.Set(i, true)
				units = append(units, ref.asUint16())
				window.append(lookahead[:ref.length])
				if _, err := src.Discard(int(ref.length)); err != nil {
					panic(err)
				}
			}
		}
		if len(units) == 0 {
			break
		}
		if err := dst.WriteBits(&unitHeader); err != nil {
			return err
		}
		for i := 0; i < len(units); i++ {
			if unitHeader.Get(i) {
				if err := dst.WriteUint16(units[i]); err != nil {
					return err
				}
			} else {
				if err := dst.WriteByte(byte(units[i])); err != nil {
					return err
				}
			}
		}
	}
	return dst.Flush()
}

// Decode reads LZ77 encoded data from input, decodes it and writes the decoded
// data to output.
func Decode(input io.Reader, output io.Writer) (err error) {
	src := bits.NewReader(input)
	dst := bufio.NewWriter(output)
	window := newWindowBuffer(windowBufferSize)
	headerBuf := make([]byte, 1)

	for {
		headerBuf[0], err = src.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		header := bits.NewList(headerBuf)
		for i := 0; i < header.Len(); i++ {
			if header.Get(i) {
				ref, err := decodeReference(src)
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				if err := window.expandReference(dst, ref); err != nil {
					return err
				}
			} else {
				next, err := src.ReadByte()
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				if err := dst.WriteByte(next); err != nil {
					return err
				}
				window.appendByte(next)
			}
		}
	}
	return dst.Flush()
}

// reference is a reference to an earlier byte sequence in the current window
// buffer.
type reference struct {
	length, distance uint16
}

// asUint16 combines length and reference into a single uint16 value.
func (r reference) asUint16() uint16 {
	return (r.length << refDistBits) | r.distance
}

// decodeReference decodes a single reference from r.
func decodeReference(r *bits.Reader) (reference, error) {
	ref, err := r.ReadUint16()
	if err != nil {
		return reference{}, err
	}
	return reference{
		length:   ref >> refDistBits,
		distance: ref & (^uint16(0) >> refLenBits),
	}, nil
}

// windowBuffer is a sliding window that keeps track of recent processed bytes
// to allow replacing future duplicate byte sequences with references.
type windowBuffer struct {
	// Contains recent bytes
	buf []byte
	// An index to buf which determines where the window logically starts.
	start int
}

// newWindowBuffer returns a windowBuffer with the specified size.
func newWindowBuffer(size int) *windowBuffer {
	return &windowBuffer{buf: make([]byte, size)}
}

// append copies bytes in data to the end of the window while discarding an
// equal amount of bytes from the beginning of the window.
func (w *windowBuffer) append(data []byte) {
	copied := copy(w.buf[w.start:], data)
	if copied < len(data) {
		copy(w.buf, data[copied:])
	}
	w.start = (w.start + len(data)) % len(w.buf)
}

// appendByte is similar to append but for a single byte.
func (w *windowBuffer) appendByte(b byte) {
	w.buf[w.start] = b
	w.start = (w.start + 1) % len(w.buf)
}

// findLongestPrefix returns a reference to the longest prefix of input found in
// the current window. A zeroed reference is returned if no prefix is found.
func (w *windowBuffer) findLongestPrefix(input []byte) reference {
	start := 0
	length := 0
	for i := 0; i < len(w.buf); i++ {
		j := 0
		for ; j < len(input) && i+j < len(w.buf); j++ {
			if w.get(i+j) != input[j] {
				break
			}
		}
		if j > length {
			start = i
			length = j
		}
	}
	if length == 0 {
		return reference{}
	}
	return reference{
		length:   uint16(length),
		distance: uint16(len(w.buf) - start),
	}
}

// expandReference expands ref by writing the corresponding byte sequence in the
// window to out and the end of the window.
func (w *windowBuffer) expandReference(out *bufio.Writer, ref reference) error {
	start := len(w.buf) - int(ref.distance)
	for i := 0; i < int(ref.length); i++ {
		byt := w.get(start)
		if err := out.WriteByte(byt); err != nil {
			return err
		}
		w.appendByte(byt)
	}
	return nil
}

// get returns the byte at logical index i in the window.
func (w *windowBuffer) get(i int) byte {
	return w.buf[(w.start+i)%len(w.buf)]
}
