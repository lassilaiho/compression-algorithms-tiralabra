package bits

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	tu "github.com/lassilaiho/compression-algorithms-tiralabra/util/testutil"
)

func (l *List) Equals(other List) bool {
	if l.len != other.len {
		return false
	}
	for i := 0; i < l.len; i++ {
		if l.Get(i) != other.Get(i) {
			return false
		}
	}
	return true
}

func TestBitReaderReadBit(t *testing.T) {
	input := []byte{0b00010110, 0b11010010, 0b11010010}
	output := []byte{
		0, 0, 0, 1, 0, 1, 1, 0,
		1, 1, 0, 1, 0, 0, 1, 0,
		1, 1, 0, 1, 0, 0, 1, 0,
	}
	r := NewReader(bytes.NewBuffer(input))
	var bit bool
	var err error
	for i := 0; i < len(output); i++ {
		bit, err = r.ReadBit()
		tu.ExpectNil(t, err)
		correct := output[i] != 0
		if bit != correct {
			t.Fatalf("expected bit %d to be %v, found %v", i, correct, bit)
		}
	}
	bit, err = r.ReadBit()
	tu.ExpectEOF(t, err)
}

func TestBitReaderReadByte(t *testing.T) {
	input := []byte{0b00010110, 0b01010010, 0b11010010}
	r := NewReader(bytes.NewBuffer(input))
	for i := 0; i < len(input); i++ {
		byt, err := r.ReadByte()
		tu.ExpectNil(t, err)
		if byt != input[i] {
			t.Fatalf("expected byte %d to be %d, found %v", i, input[i], byt)
		}
	}
	_, err := r.ReadBit()
	tu.ExpectEOF(t, err)
}

func TestBitReaderReadInt64(t *testing.T) {
	correct := int64(192479821742174211)
	input := make([]byte, unsafe.Sizeof(correct))
	binary.LittleEndian.PutUint64(input, uint64(correct))
	r := NewReader(bytes.NewBuffer(input))
	found, err := r.ReadInt64()
	tu.ExpectNil(t, err)
	tu.Check(t, correct, found)
	_, err = r.ReadInt64()
	tu.ExpectEOF(t, err)
}

func TestBitReaderReadUint16(t *testing.T) {
	correct := uint16(19174)
	input := make([]byte, unsafe.Sizeof(correct))
	binary.LittleEndian.PutUint16(input, correct)
	r := NewReader(bytes.NewBuffer(input))
	found, err := r.ReadUint16()
	tu.ExpectNil(t, err)
	tu.Check(t, correct, found)
	_, err = r.ReadUint16()
	tu.ExpectEOF(t, err)
}

func TestBitWriterWriteBit(t *testing.T) {
	input := []byte{
		0, 0, 0, 1, 0, 1, 1, 0,
		1, 1, 0, 1, 0, 0, 1, 0,
		1, 1, 0, 1, 0, 0, 1, 0,
	}
	correctOutput := []byte{0b00010110, 0b11010010, 0b11010010}
	var output bytes.Buffer
	w := NewWriter(&output)
	for _, bit := range input {
		tu.ExpectNil(t, w.WriteBit(bit != 0))
	}
	tu.ExpectNil(t, w.Flush())
	tu.Check(t, len(correctOutput), output.Len())
	outputSlice := output.Bytes()
	for i, correct := range correctOutput {
		if outputSlice[i] != correct {
			t.Fatalf("expected byte %d to be %d, found %d",
				i, correct, outputSlice[i])
		}
	}
}

func TestBitWriterWriteBits(t *testing.T) {
	correctOutput := []byte{0b00010110, 0b11010010, 0b11010010}
	input := NewList(correctOutput)
	var output bytes.Buffer
	w := NewWriter(&output)
	tu.ExpectNil(t, w.WriteBits(&input))
	tu.ExpectNil(t, w.Flush())
	t.Log(output.Bytes())
	tu.Check(t, len(correctOutput), output.Len())
	outputSlice := output.Bytes()
	for i, correct := range correctOutput {
		if outputSlice[i] != correct {
			t.Fatalf("expected byte %d to be %d, found %d",
				i, correct, outputSlice[i])
		}
	}
}

func TestBitWriterWriteInt64(t *testing.T) {
	correct := int64(192479821742174211)
	var output bytes.Buffer
	w := NewWriter(&output)
	tu.ExpectNil(t, w.WriteInt64(correct))
	tu.ExpectNil(t, w.Flush())
	tu.Check(t, correct, int64(binary.LittleEndian.Uint64(output.Bytes())))
}

func TestBitWriterWriteUint16(t *testing.T) {
	correct := uint16(7914)
	var output bytes.Buffer
	w := NewWriter(&output)
	tu.ExpectNil(t, w.WriteUint16(correct))
	tu.ExpectNil(t, w.Flush())
	tu.Check(t, correct, binary.LittleEndian.Uint16(output.Bytes()))
}

func TestBitListSet(t *testing.T) {
	bits := NewList([]byte{0, 0})
	bits.Set(3, true)
	tu.Check(t, "0001000000000000", bits.String())
	bits.Set(4, true)
	tu.Check(t, "0001100000000000", bits.String())
	bits.Set(3, false)
	tu.Check(t, "0000100000000000", bits.String())
	bits.Set(10, false)
	tu.Check(t, "0000100000000000", bits.String())
	bits.Set(11, true)
	tu.Check(t, "0000100000010000", bits.String())
}
