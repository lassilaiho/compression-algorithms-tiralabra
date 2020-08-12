package lz77

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bits"
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
)

func TestWindowBuffer(t *testing.T) {
	window := newWindowBuffer(4)
	t.Run("Append", func(t *testing.T) {
		window.append([]byte{4, 9, 1})
		for i, b := range []byte{0, 4, 9, 1} {
			if window.get(i) != b {
				t.Fatalf("expected byte %d to be %d, found %d",
					i, b, window.get(i))
			}
		}
		window.append([]byte{3, 2})
		for i, b := range []byte{9, 1, 3, 2} {
			if window.get(i) != b {
				t.Fatalf("expected byte %d to be %d, found %d",
					i, b, window.get(i))
			}
		}
	})
	t.Run("FindLongestPrefix", func(t *testing.T) {
		expected := reference{length: 2, distance: 3}
		ref := window.findLongestPrefix([]byte{1, 3})
		if ref != expected {
			t.Fatalf("expected %v, found %v", expected, ref)
		}
		expected = reference{length: 0, distance: 0}
		ref = window.findLongestPrefix([]byte{0, 1})
		if ref != expected {
			t.Fatalf("expected %v, found %v", expected, ref)
		}
		expected = reference{length: 2, distance: 3}
		ref = window.findLongestPrefix([]byte{1, 3, 6})
		if ref != expected {
			t.Fatalf("expected %v, found %v", expected, ref)
		}
	})
	t.Run("ExpandReference", func(t *testing.T) {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		window.expandReference(w, reference{
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
		if found != expected {
			t.Fatalf("expected %v, found %v", expected, found)
		}
	})
	var encoded bytes.Buffer
	w := bits.NewWriter(&encoded)
	w.WriteUint16(ref.asUint16())
	w.Flush()
	t.Run("Decode", func(t *testing.T) {
		decoded, _ := decodeReference(bits.NewReader(&encoded))
		if decoded != ref {
			t.Fatalf("expected %v, found %v", ref, decoded)
		}
	})
}

func TestEncodingAndDecoding(t *testing.T) {
	input := "aösdkfjaöslkdfjaösldkjfaösldkjföalsdkjflaskjdhfakjsdflkdsajhfaksdjhflsakdjhf"
	var encoded bytes.Buffer
	if err := Encode(strings.NewReader(input), &encoded); err != nil {
		t.Fatal(err)
	}
	var decoded strings.Builder
	if err := Decode(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	result := decoded.String()
	if result != input {
		t.Fatalf("expected %s, found %s", input, result)
	}
}

func BenchmarkEncode(b *testing.B) {
	input, err := ioutil.ReadFile("../test/kalevala.txt")
	if err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Reset(input)
		var buf bytes.Buffer
		Encode(r, &buf)
	}
}
