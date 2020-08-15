package huffman

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bits"
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
	tu "github.com/lassilaiho/compression-algorithms-tiralabra/util/testutil"
)

const (
	testKalevala = "../test/kalevala.txt"
)

func printTree(node *codeTreeNode, indent string) {
	if node.left == nil {
		fmt.Println(indent, string([]byte{node.symbol}))
	} else {
		fmt.Println(indent + "X")
		printTree(node.left, indent+"0 ")
		printTree(node.right, indent+"1 ")
	}
}

func checkTrees(t *testing.T, expected, found *codeTreeNode) {
	t.Helper()
	if expected.left == nil {
		if found.left == nil {
			tu.Check(t, expected.symbol, found.symbol)
		} else {
			t.Fatal("expected leaf node, found internal node")
		}
	} else if found.left == nil {
		t.Fatal("expected internal node, found leaf node")
	} else {
		checkTrees(t, expected.left, found.left)
		checkTrees(t, expected.right, found.right)
	}
}

func TestEncoding(t *testing.T) {
	input := "45621354622615342165326143453614216346214"
	var freqs frequencyTable
	t.Run("CountFrequencies", func(t *testing.T) {
		tu.ExpectNil(t, countFrequencies(
			bufio.NewReader(strings.NewReader(input)), &freqs))
		tu.Check(t, 7, int(freqs['1']))
		tu.Check(t, 7, int(freqs['2']))
		tu.Check(t, 6, int(freqs['3']))
		tu.Check(t, 8, int(freqs['4']))
		tu.Check(t, 5, int(freqs['5']))
		tu.Check(t, 8, int(freqs['6']))
		tu.Check(t, 41, int(freqs.byteCount()))
	})
	var codeTree *codeTreeNode
	var codeTable *codeTable
	t.Run("NewCodeTable", func(t *testing.T) {
		codeTree = buildCodeTree(&freqs)
		codeTable = newCodeTable(codeTree)

		t.Log("1:", codeTable['1'].String())
		t.Log("2:", codeTable['2'].String())
		t.Log("3:", codeTable['3'].String())
		t.Log("4:", codeTable['4'].String())
		t.Log("5:", codeTable['5'].String())
		t.Log("6:", codeTable['6'].String())

		tu.Check(t, "111", codeTable['1'].String())
		tu.Check(t, "110", codeTable['2'].String())
		tu.Check(t, "101", codeTable['3'].String())
		tu.Check(t, "00", codeTable['4'].String())
		tu.Check(t, "100", codeTable['5'].String())
		tu.Check(t, "01", codeTable['6'].String())
	})
	t.Run("CodeTreeEncode", func(t *testing.T) {
		var output bytes.Buffer
		writer := bits.NewWriter(&output)
		tu.ExpectNil(t, codeTree.encodeTo(writer))
		writer.Flush()
		bitList := bits.NewList(output.Bytes())
		tu.Check(t,
			"0010011010010011011000100110101100110011010011001010011000100000",
			bitList.String())
	})
	t.Run("CodeTableEncode", func(t *testing.T) {
		var output bytes.Buffer
		tu.ExpectNil(t, codeTable.Encode(
			bufio.NewReader(strings.NewReader(input)),
			bits.NewWriter(&output)))
		bitList := bits.NewList(output.Bytes())
		tu.Check(t,
			"0010001110111101100000111011001111100101001101110110010111001111001010010010101111001101110110100011101110000000",
			bitList.String())
	})
}

func TestDecoding(t *testing.T) {
	data := "45621354622615342165326143453614216346214"
	var encodedData bytes.Buffer
	tu.ExpectNil(t, Encode(strings.NewReader(data), &encodedData))
	input := bits.NewReader(&encodedData)

	var codeTree *codeTreeNode
	var err error
	t.Run("DecodeTree", func(t *testing.T) {
		expected := &codeTreeNode{
			left: &codeTreeNode{
				left:  &codeTreeNode{symbol: '4'},
				right: &codeTreeNode{symbol: '6'},
			},
			right: &codeTreeNode{
				left: &codeTreeNode{
					left:  &codeTreeNode{symbol: '5'},
					right: &codeTreeNode{symbol: '3'},
				},
				right: &codeTreeNode{
					left:  &codeTreeNode{symbol: '2'},
					right: &codeTreeNode{symbol: '1'},
				},
			},
		}
		codeTree, err = decodeCodeTree(input)
		tu.ExpectNil(t, err)
		checkTrees(t, expected, codeTree)
	})
	var byteCount int64
	t.Run("DecodeByteCount", func(t *testing.T) {
		byteCount, err = input.ReadInt64()
		tu.ExpectNil(t, err)
		tu.Check(t, int64(len(data)), byteCount)
	})
	t.Run("DecodeData", func(t *testing.T) {
		for i := int64(0); i < byteCount; i++ {
			byt, err := codeTree.readCode(input)
			tu.ExpectNil(t, err)
			if byt != data[i] {
				t.Fatalf("expected byte %d to be %d, found %d", i, data[i], byt)
			}
		}
	})
}

func TestDecode(t *testing.T) {
	cases := []struct {
		desc string
		data []byte
	}{
		{
			desc: "RandomNumbers",
			data: []byte("45621354622615342165326143453614216346214"),
		},
		{
			desc: "Kalevala",
			data: tu.ReadFile(testKalevala),
		},
	}
	var buf bytes.Buffer
	var output bytes.Buffer
	for _, c := range cases {
		buf.Reset()
		output.Reset()
		t.Run(c.desc, func(t *testing.T) {
			tu.ExpectNil(t, Encode(bytes.NewReader(c.data), &buf))
			tu.ExpectNil(t, Decode(&buf, &output))
			if !bytes.Equal(c.data, output.Bytes()) {
				if len(c.data) < 200 {
					t.Fatalf("expected %v, found %v", string(c.data), output.String())
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
