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
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/slices"
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
				units = slices.AppendUint16(units, uint16(next))
				window.appendByte(next)
			} else {
				unitHeader.Set(i, true)
				units = slices.AppendUint16(units, ref.asUint16())
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
	// Contains recent bytes.
	buf []byte
	// An index to buf which determines where the window logically starts.
	start int
	// A dictionary used to speed up prefix matching performance.
	dict dictionary
	// A monotonically increasing counter representing the current position in
	// the data stream.
	pos int64
}

// newWindowBuffer returns a windowBuffer with the specified size.
func newWindowBuffer(size int) *windowBuffer {
	return &windowBuffer{
		buf:  make([]byte, size),
		dict: dictionary{},
		pos:  int64(size),
	}
}

// append copies bytes in data to the end of the window while discarding an
// equal amount of bytes from the beginning of the window.
func (w *windowBuffer) append(data []byte) {
	// Remove discarded byte sequences from dictionary.
	w.pos += int64(len(data))
	for i := 0; i < len(data); i++ {
		w.dict.removeLesserThan(w.dictKey(i), w.pos-int64(len(w.buf))+int64(i))
	}

	// Copy data to buffer.
	copied := slices.CopyBytes(w.buf[w.start:], data)
	if copied < len(data) {
		slices.CopyBytes(w.buf, data[copied:])
	}
	w.start = (w.start + len(data)) % len(w.buf)

	// Add new byte sequences to dictionary.
	pos := w.pos - dictKeySize - int64(len(data)) + 1
	bufIndex := len(w.buf) - dictKeySize - len(data) + 1
	if bufIndex < 0 {
		pos -= int64(bufIndex)
		bufIndex = 0
	}
	for i := 0; i < len(data); i++ {
		w.dict.add(w.dictKey(bufIndex+i), pos+int64(i))
	}
}

// appendByte is similar to append but for a single byte.
func (w *windowBuffer) appendByte(b byte) {
	w.append([]byte{b})
}

// findLongestPrefix returns a reference to the longest prefix of input found in
// the current window. A zeroed reference is returned if no prefix is found.
func (w *windowBuffer) findLongestPrefix(input []byte) reference {
	if len(input) < 2 {
		return reference{}
	}
	start := 0
	length := 0
	value := w.dict.get(dictKey{input[0], input[1]})
	for value != nil {
		i := int(value.value) - (int(w.pos) - len(w.buf))
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
		value = value.next
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

// dictKey returns the dictionary key corresponding to logical position i in the
// window.
func (w *windowBuffer) dictKey(i int) dictKey {
	key := dictKey{}
	for j := 0; j < len(key); j++ {
		key[j] = w.get(i + j)
	}
	return key
}

const dictKeySize = 2

type dictKey [dictKeySize]byte

type dictEntry struct {
	first *dictValue
	last  *dictValue
}

type dictValue struct {
	value int64
	next  *dictValue
}

// dictionary maps keys into byte sequence positions in a data stream.
//
// A key is formed from the initial bytes of the sequence. The number of bytes
// used is specified using the constant dictKeySize. Sequences shorter than this
// can't be stored in the dictionary.
//
// Each key maps to a single entry. An entry is a linked list of values. The
// values in a entry must be added in ascending order. The addition order is not
// enforced. It is the caller's responsibility to add values in the correct
// order.
type dictionary map[dictKey]*dictEntry

// add adds value to the end of the entry corresponding to key.
func (d dictionary) add(key dictKey, value int64) {
	entry := d[key]
	if entry == nil {
		entry = &dictEntry{}
		d[key] = entry
	}
	if entry.first == nil {
		entry.first = &dictValue{value: value}
		entry.last = entry.first
	} else {
		entry.last.next = &dictValue{value: value}
		entry.last = entry.last.next
	}
}

// get returns the first value corresponding to key or nil if the entry
// corresponding to key has no values.
func (d dictionary) get(key dictKey) *dictValue {
	entry := d[key]
	if entry == nil {
		return nil
	}
	return entry.first
}

// removeLesserThan removes all values lesser to value from the entry
// corresponding to key.
func (d dictionary) removeLesserThan(key dictKey, value int64) {
	entry := d[key]
	if entry == nil {
		return
	}
	dictValue := entry.first
	for dictValue != nil && dictValue.value < value {
		dictValue = dictValue.next
	}
	entry.first = dictValue
	if dictValue == nil {
		entry.last = nil
	}
}
