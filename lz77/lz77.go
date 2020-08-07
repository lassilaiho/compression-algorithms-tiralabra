// Package lz77 implements LZ77 encoding and decoding.
//
// Encode encodes data into triplets (l, d, n) where l is the length of the
// sequence of bytes this triplet refers to, d is the number of characters
// behind the current position the sequence starts and n is the next byte after
// the sequence ends. If length and distance are 0 the triplet doesn't refer to
// anything and is considered a literal byte n. l and d are encoded together as
// a little-endian 16 bit unit. n is a single byte.
package lz77

import (
	"encoding/binary"
	"io"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
)

// These constants specify how the 16 bits of a
// reference are distributed between length and distance.
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
	dst := bufio.NewWriter(output)
	window := newWindowBuffer(windowBufferSize)

	for {
		lookahead, err := src.Peek(lookaheadBufferSize)
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(lookahead) == 0 {
				break
			}
		}
		var next byte
		ref := window.findLongestPrefix(lookahead)
		window.append(lookahead[:ref.length])
		if err := ref.encode(dst); err != nil {
			return err
		}
		if ref.length == 0 {
			next = lookahead[0]
			if _, err := src.Discard(1); err != nil {
				panic(err)
			}
		} else {
			if _, err := src.Discard(int(ref.length)); err != nil {
				panic(err)
			}
			if next, err = src.ReadByte(); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
		if err := dst.WriteByte(next); err != nil {
			return err
		}
		window.appendByte(next)
	}
	return dst.Flush()
}

// Decode reads LZ77 encoded data from input, decodes it and writes the decoded
// data to output.
func Decode(input io.Reader, output io.Writer) error {
	src := bufio.NewReaderSize(input, lookaheadBufferSize)
	dst := bufio.NewWriter(output)
	window := newWindowBuffer(windowBufferSize)

	for {
		ref, err := decodeReference(src)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if ref.length != 0 {
			if err := window.expandReference(dst, ref); err != nil {
				return err
			}
		}
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
	return dst.Flush()
}

// reference is a reference to an earlier byte sequence in the current window
// buffer.
type reference struct {
	length, distance uint16
}

// encode encodes r to w as a single uint16 with little-endian byte order.
func (r reference) encode(w io.Writer) error {
	ref := (r.length << refDistBits) | r.distance
	return binary.Write(w, binary.LittleEndian, ref)
}

// decodeReference decodes a single reference from r.
func decodeReference(r io.Reader) (reference, error) {
	var ref uint16
	if err := binary.Read(r, binary.LittleEndian, &ref); err != nil {
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
