package huffman

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func bitListOfBits(bits []byte) bitList {
	list := newBitList([]byte{})
	for _, bit := range bits {
		list.Append(bit != 0)
	}
	return list
}

func printTree(t *testing.T, node *codeTreeNode, indent string) {
	if node.left == nil {
		fmt.Println(indent, string([]byte{node.symbol}))
	} else {
		fmt.Println(indent + "X")
		printTree(t, node.left, indent+"0 ")
		printTree(t, node.right, indent+"1 ")
	}
}

func checkTrees(t *testing.T, expected, found *codeTreeNode) {
	if expected.left == nil {
		if found.left == nil {
			check(t, expected.symbol, found.symbol)
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
		expectNil(t, countFrequencies(
			bufio.NewReader(strings.NewReader(input)), &freqs))
		check(t, 7, int(freqs['1']))
		check(t, 7, int(freqs['2']))
		check(t, 6, int(freqs['3']))
		check(t, 8, int(freqs['4']))
		check(t, 5, int(freqs['5']))
		check(t, 8, int(freqs['6']))
		check(t, 41, int(freqs.byteCount()))
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

		check(t, true, codeTable['1'].Equals(bitListOfBits([]byte{1, 1, 1})))
		check(t, true, codeTable['2'].Equals(bitListOfBits([]byte{1, 1, 0})))
		check(t, true, codeTable['3'].Equals(bitListOfBits([]byte{1, 0, 1})))
		check(t, true, codeTable['4'].Equals(bitListOfBits([]byte{0, 0})))
		check(t, true, codeTable['5'].Equals(bitListOfBits([]byte{1, 0, 0})))
		check(t, true, codeTable['6'].Equals(bitListOfBits([]byte{0, 1})))
	})
	t.Run("CodeTreeEncode", func(t *testing.T) {
		var output bytes.Buffer
		writer := newBitWriter(&output)
		expectNil(t, codeTree.encodeTo(writer))
		writer.Flush()
		bits := newBitList(output.Bytes())
		check(t,
			"0010011010010011011000100110101100110011010011001010011000100000",
			bits.String())
	})
	t.Run("CodeTableEncode", func(t *testing.T) {
		var output bytes.Buffer
		expectNil(t, codeTable.Encode(
			bufio.NewReader(strings.NewReader(input)),
			newBitWriter(&output)))
		bits := newBitList(output.Bytes())
		check(t,
			"0010001110111101100000111011001111100101001101110110010111001111001010010010101111001101110110100011101110000000",
			bits.String())
	})
}

func TestDecoding(t *testing.T) {
	data := "45621354622615342165326143453614216346214"
	var encodedData bytes.Buffer
	expectNil(t, Encode(strings.NewReader(data), &encodedData))
	input := newBitReader(&encodedData)

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
		expectNil(t, err)
		checkTrees(t, expected, codeTree)
	})
	var byteCount int64
	t.Run("DecodeByteCount", func(t *testing.T) {
		byteCount, err = input.ReadInt64()
		expectNil(t, err)
		check(t, int64(len(data)), byteCount)
	})
	t.Run("DecodeData", func(t *testing.T) {
		for i := int64(0); i < byteCount; i++ {
			byt, err := codeTree.readCode(input)
			expectNil(t, err)
			if byt != data[i] {
				t.Fatalf("expected byte %d to be %d, found %d", i, data[i], byt)
			}
		}
	})
}

func TestDecode(t *testing.T) {
	input := "45621354622615342165326143453614216346214"
	var buf bytes.Buffer
	var output strings.Builder
	expectNil(t, Encode(strings.NewReader(input), &buf))
	expectNil(t, Decode(&buf, &output))
	check(t, input, output.String())
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
