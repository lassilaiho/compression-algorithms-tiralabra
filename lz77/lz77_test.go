package lz77

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bits"
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
	tu "github.com/lassilaiho/compression-algorithms-tiralabra/util/testutil"
)

const (
	testKalevala = "../test/files/kalevala.txt"
)

func (d *dictionary) forEach(t *testing.T, f func(*dictEntry)) {
	t.Helper()
	for _, bucket := range d.buckets {
		for _, entry := range bucket {
			f(entry)
		}
	}
}

func newDictOf(entries ...*dictEntry) *dictionary {
	dict := newDictionary()
	for _, e := range entries {
		dict.setEntry(e)
	}
	return dict
}

func checkDictValue(t *testing.T, found, expected *dictValue) {
	t.Helper()
	if expected == nil && found == nil {
		return
	} else if expected != nil && found != nil {
		tu.Check(t, expected.value, found.value)
		checkDictValue(t, found.next, expected.next)
	} else {
		t.Fatalf("expected %v, found %v", expected, found)
	}
}

func checkDict(t *testing.T, found, expected *dictionary) {
	t.Helper()
	expected.forEach(t, func(eentry *dictEntry) {
		t.Helper()
		key := eentry.key
		fentry := found.get(key)
		if eentry == nil && fentry == nil {
			return
		} else if eentry != nil && fentry != nil {
			checkDictValue(t, found.get(key), eentry.first)
		} else {
			t.Fatalf("expected %v, found %v", eentry, fentry)
		}
	})
}

func TestEncoderWindowBuffer(t *testing.T) {
	window := newEncoderWindowBuffer(4)
	t.Run("Append", func(t *testing.T) {
		window.append([]byte{4, 9, 1})
		for i, b := range []byte{0, 4, 9, 1} {
			if window.get(i) != b {
				t.Fatalf("expected byte %d to be %d, found %d",
					i, b, window.get(i))
			}
		}
		tu.Check(t, int64(7), window.pos)
		checkDict(t, window.dict, newDictOf(
			&dictEntry{
				key:   dictKey{0, 4},
				first: &dictValue{value: 3},
			},
			&dictEntry{
				key:   dictKey{4, 9},
				first: &dictValue{value: 4},
			},
			&dictEntry{
				key:   dictKey{9, 1},
				first: &dictValue{value: 5},
			},
		))

		window.append([]byte{3, 2})
		for i, b := range []byte{9, 1, 3, 2} {
			if window.get(i) != b {
				t.Fatalf("expected byte %d to be %d, found %d",
					i, b, window.get(i))
			}
		}
		tu.Check(t, int64(9), window.pos)
		checkDict(t, window.dict, newDictOf(
			&dictEntry{
				key:   dictKey{9, 1},
				first: &dictValue{value: 5},
			},
			&dictEntry{
				key:   dictKey{1, 3},
				first: &dictValue{value: 6},
			},
			&dictEntry{
				key:   dictKey{3, 2},
				first: &dictValue{value: 7},
			},
		))
	})
	t.Run("FindLongestPrefix", func(t *testing.T) {
		expected := reference{length: 2, distance: 3}
		ref := window.findLongestPrefix([]byte{1, 3})
		if ref != expected {
			for i := range window.win.buf {
				fmt.Print(window.get(i), " ")
			}
			fmt.Println()
			t.Fatalf("expected %v, found %v", expected, ref)
		}
		expected = reference{length: 0, distance: 0}
		ref = window.findLongestPrefix([]byte{0, 1})
		tu.Check(t, expected, ref)
		expected = reference{length: 2, distance: 3}
		ref = window.findLongestPrefix([]byte{1, 3, 6})
		tu.Check(t, expected, ref)
	})
	t.Run("ExpandReference", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		window.win.expandReference(w, reference{
			length:   3,
			distance: 4,
		})
		w.Flush()
		expected := []byte{9, 1, 3}
		found := buf.Bytes()
		if !bytes.Equal(expected, found) {
			t.Fatalf("expected %v, found %v", expected, found)
		}
		for i, b := range []byte{2, 9, 1, 3} {
			if window.get(i) != b {
				t.Fatalf("expected byte %d to be %d, found %d",
					i, b, window.get(i))
			}
		}
	})
}

func TestReference(t *testing.T) {
	ref := reference{
		length:   0b1000,
		distance: 0b100111100110,
	}
	t.Run("AsUint16", func(t *testing.T) {
		expected := uint16(0b10001001_11100110)
		found := ref.asUint16()
		tu.Check(t, expected, found)
	})
	var encoded bytes.Buffer
	w := bits.NewWriter(&encoded)
	w.WriteUint16(ref.asUint16())
	w.Flush()
	t.Run("Decode", func(t *testing.T) {
		decoded, _ := decodeReference(bits.NewReader(&encoded))
		tu.Check(t, ref, decoded)
	})
}

func TestEncodingAndDecoding(t *testing.T) {
	cases := []struct {
		desc string
		data []byte
	}{
		{
			desc: "Random",
			data: []byte("aösdkfjaöslkdfjaösldkjfaösldkjföalsdkjflaskjdhfakjsdflkdsajhfaksdjhflsakdjhf"),
		},
		{
			desc: "Kalevala",
			data: tu.ReadFile(testKalevala),
		},
	}
	var encoded bytes.Buffer
	var decoded bytes.Buffer
	for _, c := range cases {
		encoded.Reset()
		decoded.Reset()
		t.Run(c.desc, func(t *testing.T) {
			if err := Encode(bytes.NewReader(c.data), &encoded); err != nil {
				t.Fatal(err)
			}
			if err := Decode(&encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(c.data, decoded.Bytes()) {
				if len(c.data) < 200 {
					t.Fatalf("expected %v, found %v", string(c.data), decoded.String())
				} else {
					t.FailNow()
				}
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	input := tu.ReadFile(testKalevala)
	r := bytes.NewReader(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Reset(input)
		var buf bytes.Buffer
		Encode(r, &buf)
	}
}
