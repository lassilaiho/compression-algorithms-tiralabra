package huffman

import (
	"io"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
)

// bitList is a growable packed list of bits. The zero value is an empty list
// ready for use.
type bitList struct {
	len int
	buf []byte
}

// newBitList returns a list that uses buf for it's initial contents. The
// returned list takes ownership of buf and it shouldn't be used after passing
// it to this function.
func newBitList(buf []byte) bitList {
	return bitList{
		len: 8 * len(buf),
		buf: buf,
	}
}

func (l *bitList) Append(bit bool) {
	if l.len/8 >= len(l.buf) {
		l.buf = append(l.buf, 0)
	}
	l.len++
	l.Set(l.len-1, bit)
}

func (l *bitList) Get(i int) bool {
	return (l.buf[i/8]>>(7-i%8))&1 != 0
}

func (l *bitList) Set(i int, bit bool) {
	if bit {
		l.buf[i/8] |= 1 << byte(7-i%8)
	} else {
		l.buf[i/8] &= ^(1 << byte(7-i%8))
	}
}

func (l *bitList) Len() int {
	return l.len
}

// Shrink shirnks the length of l by n. n must be in range [0, l.Len()). Shrink
// only reduces the length of l. No memory is freed.
func (l *bitList) Shrink(n int) {
	l.len -= n
}

func (l *bitList) Copy() bitList {
	copied := bitList{
		len: l.len,
		buf: make([]byte, len(l.buf)),
	}
	copy(copied.buf, l.buf)
	return copied
}

// bitWriter is used to write individual bits into an io.Writer.
//
// All writes to a bitWriter are buffered. Calling Flush writes all buffered
// data to the underlying io.Writer along with possible additional trailing zero
// bits to round the data to full bytes.
type bitWriter struct {
	w   *bufio.Writer
	i   byte
	buf byte
}

// newBitWriter returns a bitWriter that writes to w.
func newBitWriter(w io.Writer) *bitWriter {
	return &bitWriter{w: bufio.NewWriter(w)}
}

func (w *bitWriter) WriteBit(b bool) error {
	if b {
		w.buf |= 1 << (7 - w.i)
	}
	if w.i == 7 {
		if err := w.w.WriteByte(w.buf); err != nil {
			return err
		}
		w.i = 0
		w.buf = 0
	} else {
		w.i++
	}
	return nil
}

func (w *bitWriter) WriteBits(bits *bitList) error {
	for i := 0; i < bits.Len(); i++ {
		if err := w.WriteBit(bits.Get(i)); err != nil {
			return err
		}
	}
	return nil
}

func (w *bitWriter) WriteByte(n byte) error {
	bits := newBitList([]byte{n})
	return w.WriteBits(&bits)
}

// WriteInt64 writes n to the writer using little-endian byte order.
func (w *bitWriter) WriteInt64(n int64) error {
	bits := newBitList([]byte{
		byte(uint64(n)),
		byte(uint64(n) >> 8),
		byte(uint64(n) >> 16),
		byte(uint64(n) >> 24),
		byte(uint64(n) >> 32),
		byte(uint64(n) >> 40),
		byte(uint64(n) >> 48),
		byte(uint64(n) >> 56),
	})
	return w.WriteBits(&bits)
}

func (w *bitWriter) Flush() error {
	if w.i > 0 {
		if err := w.w.WriteByte(w.buf); err != nil {
			return err
		}
		w.i = 0
		w.buf = 0
	}
	return w.w.Flush()
}

// bitReader is used to read individual bits from an io.Reader.
type bitReader struct {
	w   *bufio.Reader
	i   byte
	buf byte
}

// newBitReader returns a bitReader that reads from r.
func newBitReader(r io.Reader) *bitReader {
	return &bitReader{
		w: bufio.NewReader(r),
		i: 8,
	}
}

func (w *bitReader) ReadBit() (bit bool, err error) {
	if w.i == 8 {
		w.buf, err = w.w.ReadByte()
		if err != nil {
			return false, err
		}
		w.i = 0
	}
	bit = (w.buf>>(7-w.i))&1 != 0
	w.i++
	return bit, nil
}

func (w *bitReader) readBitPanicing() byte {
	bit, err := w.ReadBit()
	if err != nil {
		panic(readBitErr(err))
	}
	if bit {
		return 1
	}
	return 0
}

func (w *bitReader) ReadByte() (byt byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(readBitErr); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	return w.readBytePanicking(), nil
}

func (w *bitReader) readBytePanicking() byte {
	return w.readBitPanicing()<<7 |
		w.readBitPanicing()<<6 |
		w.readBitPanicing()<<5 |
		w.readBitPanicing()<<4 |
		w.readBitPanicing()<<3 |
		w.readBitPanicing()<<2 |
		w.readBitPanicing()<<1 |
		w.readBitPanicing()
}

// ReadInt64 reads an int64 value in little-endian byte order.
func (w *bitReader) ReadInt64() (n int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(readBitErr); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()
	x := uint64(w.readBytePanicking()) |
		uint64(w.readBytePanicking())<<8 |
		uint64(w.readBytePanicking())<<16 |
		uint64(w.readBytePanicking())<<24 |
		uint64(w.readBytePanicking())<<32 |
		uint64(w.readBytePanicking())<<40 |
		uint64(w.readBytePanicking())<<48 |
		uint64(w.readBytePanicking())<<56
	return int64(x), nil
}

// readBitErr is used to differentiate panics caused by panicking variants of
// read and write methods on bitReader and bitWriter.
type readBitErr error
