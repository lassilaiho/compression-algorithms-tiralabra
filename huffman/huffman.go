/*
Package huffman implements the Huffman coding algorithm. Data can be encoded
and decoded using Encode and Decode, respectively.

The output of Encode is formatted as follows:

	encoded code tree
	size of uncompressed data as a little endian int64 value
	encoded data
	possible zero bits to pad the result to full bytes
*/
package huffman

import (
	"io"

	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bits"
	"github.com/lassilaiho/compression-algorithms-tiralabra/util/bufio"
)

// Encode encodes all data from input using Huffman coding and writes the result
// to output.
func Encode(input io.ReadSeeker, output io.Writer) error {
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return err
	}
	src := bufio.NewReader(input)
	var freqs frequencyTable
	if err := countFrequencies(src, &freqs); err != nil {
		return err
	}
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return err
	}
	src.Reset(input)
	dst := bits.NewWriter(output)
	tree := buildCodeTree(&freqs)
	if tree == nil {
		return io.EOF
	}
	table := newCodeTable(tree)
	if err := tree.encodeTo(dst); err != nil {
		return err
	}
	if err := dst.WriteInt64(freqs.byteCount()); err != nil {
		return err
	}
	return table.Encode(src, dst)
}

// Decode decodes data encoded using Encode from input and writes the unencoded
// data to output.
func Decode(input io.Reader, output io.Writer) error {
	src := bits.NewReader(input)
	dst := bufio.NewWriter(output)
	codeTree, err := decodeCodeTree(src)
	if err != nil {
		return err
	}
	byteCount, err := src.ReadInt64()
	if err != nil {
		return err
	}
	for ; byteCount > 0; byteCount-- {
		byt, err := codeTree.readCode(src)
		if err != nil {
			return err
		}
		if err := dst.WriteByte(byt); err != nil {
			return err
		}
	}
	return dst.Flush()
}

// codeTable maps byte values to Huffman codes.
type codeTable [256]bits.List

// Encode encodes all data in src using table and writes the result to dst.
func (table *codeTable) Encode(src *bufio.Reader, dst *bits.Writer) error {
	for {
		b, err := src.ReadByte()
		if err != nil {
			if err == io.EOF {
				return dst.Flush()
			}
			return err
		}
		if err := dst.WriteBits(&table[b]); err != nil {
			return err
		}
	}
}

// newCodeTable constructs the codeTable corresponding codeTree.
func newCodeTable(codeTree *codeTreeNode) *codeTable {
	table := &codeTable{}
	code := bits.List{}
	buildCodeTable(table, &code, codeTree)
	return table
}

func buildCodeTable(table *codeTable, code *bits.List, codeTree *codeTreeNode) {
	if codeTree.left == nil {
		table[codeTree.symbol] = code.Copy()
		return
	}
	code.Append(false)
	buildCodeTable(table, code, codeTree.left)
	code.Set(code.Len()-1, true)
	buildCodeTable(table, code, codeTree.right)
	code.Shrink(1)
}

type frequencyTable [256]int64

// countFrequencies counts the occurrences of each byte value in input and
// writes the results to freqs. Entries for byte values not encountered in input
// are not zeroed in freqs.
func countFrequencies(input *bufio.Reader, freqs *frequencyTable) error {
	for {
		b, err := input.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		freqs[b]++
	}
}

// byteCount returns the total size in bytes of the data represnted by t.
func (t *frequencyTable) byteCount() int64 {
	var n int64
	for i := 0; i < len(t); i++ {
		n += t[i]
	}
	return n
}

// codeTreeNode is a node in a code tree. The tree is always a complete binary
// tree.
type codeTreeNode struct {
	left, right *codeTreeNode // both nil iff the node is a leaf node
	symbol      byte          // meaningless for non-leaf nodes
}

// buildCodeTree builds a code tree using freqs.
func buildCodeTree(freqs *frequencyTable) *codeTreeNode {
	queue := priorityQueue{}
	for symbol := 0; symbol < len(freqs); symbol++ {
		freq := freqs[symbol]
		if freq > 0 {
			queue.Append(&queueItem{
				node: &codeTreeNode{
					symbol: byte(symbol),
				},
				frequency: freq,
			})
		}
	}
	queue.Init()
	if queue.Len() == 0 {
		return nil
	}
	for queue.Len() >= 2 {
		left := queue.Pop()
		right := queue.Pop()
		queue.Push(&queueItem{
			node: &codeTreeNode{
				left:  left.node,
				right: right.node,
			},
			frequency: left.frequency + right.frequency,
		})
	}
	return queue.Pop().node
}

// encodeTo encodes tree and writes the result to out.
func (tree *codeTreeNode) encodeTo(out *bits.Writer) error {
	if tree.left == nil {
		if err := out.WriteBit(true); err != nil {
			return err
		}
		return out.WriteByte(tree.symbol)
	}
	if err := out.WriteBit(false); err != nil {
		return err
	}
	if err := tree.left.encodeTo(out); err != nil {
		return err
	}
	return tree.right.encodeTo(out)
}

// decodeCodeTree decodes a code tree from src that was previously encoded using
// encodeTo.
func decodeCodeTree(src *bits.Reader) (*codeTreeNode, error) {
	bit, err := src.ReadBit()
	if err != nil {
		return nil, err
	}
	if bit {
		symbol, err := src.ReadByte()
		if err != nil {
			return nil, err
		}
		return &codeTreeNode{symbol: symbol}, nil
	}
	left, err := decodeCodeTree(src)
	if err != nil {
		return nil, err
	}
	right, err := decodeCodeTree(src)
	if err != nil {
		return nil, err
	}
	return &codeTreeNode{left: left, right: right}, nil
}

// readCode reads a code from src and returns the corresponding byte value.
func (tree *codeTreeNode) readCode(src *bits.Reader) (byte, error) {
	if tree.left == nil {
		return tree.symbol, nil
	}
	bit, err := src.ReadBit()
	if err != nil {
		return 0, err
	}
	if bit {
		return tree.right.readCode(src)
	}
	return tree.left.readCode(src)
}
