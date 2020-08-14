// Package bufio implements parts of the standard library package "bufio".
package bufio

import (
	"io"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/slices"
)

// defaultBufSize is the default size in bytes of the buffer for Reader and
// Writer.
const defaultBufSize = 4096

// Reader buffers reads from an io.Reader.
type Reader struct {
	rd  io.Reader // The underlying io.Reader
	buf []byte    // Buffer for storing bytes

	// next and end are indices to buf that denote the range of valid buffered
	// bytes. Bytes in range [next, end) are considered as not having been read
	// yet.
	next, end int

	// The most recent error encountered when reading from rd.
	err error
}

// NewReader returns a Reader that reads from rd.
func NewReader(rd io.Reader) *Reader {
	return NewReaderSize(rd, defaultBufSize)
}

// NewReaderSize returns a Reader with the specified buffer size that reads from
// rd.
func NewReaderSize(rd io.Reader, size int) *Reader {
	return &Reader{
		rd:  rd,
		buf: make([]byte, size),
	}
}

// Read reads data into p. total is the amount of bytes read. If it is less than
// len(p), a non-nil error is also returned to explain why the read failed.
func (r *Reader) Read(p []byte) (total int, err error) {
	goal := len(p)
	for {
		n := slices.CopyBytes(p, r.buf[r.next:r.end])
		total += n
		r.next += n
		if total == goal {
			return total, nil
		}
		if r.err != nil {
			err, r.err = r.err, nil
			return total, err
		}
		p = p[n:]
		r.end, r.err = r.rd.Read(r.buf)
		r.next = 0
	}
}

// ReadByte reads a single byte from r. If no byte is available, a non-nil error
// is returned.
func (r *Reader) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := r.Read(buf[:])
	return buf[0], err
}

// Reset discards all buffered data and resets r to read from rd.
func (r *Reader) Reset(rd io.Reader) {
	r.rd = rd
	for i := range r.buf {
		r.buf[i] = 0
	}
	r.next = 0
	r.end = 0
	r.err = nil
}

// Discard discards the next n bytes. total is the number of bytes discarded. If
// total < n, a non-nil error is returned.
func (r *Reader) Discard(n int) (total int, err error) {
	for {
		if n <= r.end-r.next {
			r.next += n
			return total + n, nil
		}
		if r.err != nil {
			err, r.err = r.err, nil
			return total, err
		}
		n -= r.end - r.next
		r.end, r.err = r.rd.Read(r.buf)
		r.next = 0
	}
}

// Peek returns the next n bytes without advancing the Reader. The bytes are
// valid until the next read. If Peek returns fewer than n bytes, a non-nil
// error is returned.
func (r *Reader) Peek(n int) (buf []byte, err error) {
	if n <= r.end-r.next {
		return r.buf[r.next:r.end], nil
	}
	if r.err != nil {
		err, r.err = r.err, nil
		return r.buf[r.next:r.end], err
	}
	if r.next != 0 {
		slices.CopyBytes(r.buf, r.buf[r.next:r.end])
		r.end -= r.next
		r.next = 0
	}
	var x int
	x, r.err = r.rd.Read(r.buf[r.end:])
	r.end += x
	if n <= r.end {
		return r.buf[:r.end], nil
	}
	err, r.err = r.err, nil
	return r.buf[:r.end], err
}

// Writer implements buffering for an io.Writer.
type Writer struct {
	wr  io.Writer // The underlying io.Reader
	buf []byte    // Buffer for storing bytes

	// next and end are indices to buf that denote the range of valid buffered
	// bytes. Bytes in range [next, end) are considered written but not yet
	// flushed to the underlying io.Writer.
	next, end int
}

// NewWriter returns a Writer that writes to wr.
func NewWriter(wr io.Writer) *Writer {
	return &Writer{
		wr:  wr,
		buf: make([]byte, defaultBufSize),
	}
}

// Write writes data from p. total is the number of bytes written. If total <
// len(p), a non-nil error is returned.
func (w *Writer) Write(p []byte) (total int, err error) {
	goal := len(p)
	for {
		n := slices.CopyBytes(w.buf[w.end:], p)
		w.end += n
		total += n
		if total == goal {
			return total, nil
		}
		p = p[n:]
		n, err := w.wr.Write(w.buf[w.next:w.end])
		w.next += n
		if err != nil {
			return total, err
		}
		w.next = 0
		w.end = 0
	}
}

// Flush writes all buffered data to the underlying io.Writer. A non-nil error
// is returned if an error occurs.
func (w *Writer) Flush() error {
	n, err := w.wr.Write(w.buf[w.next:w.end])
	w.next += n
	if err != nil {
		return err
	}
	w.next = 0
	w.end = 0
	return nil
}

// WriteByte writes a single byte to w.
func (w *Writer) WriteByte(b byte) error {
	buf := [1]byte{b}
	_, err := w.Write(buf[:])
	return err
}
