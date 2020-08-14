# Implementation document

## Pakcage structure

- `github.com/lassilaiho/compression-algorithms-tiralabra`
  - `cmd`
    - `huffman` - Command line interface for Huffman coding
    - `lz77` - Command line interface for LZ77
  - `huffman` - Huffman coding implementation
  - `lz77` - LZ77 implementation
  - `test/cmd` - Test program for benchmarking the algorithms
  - `util`
    - `bits` - Utilities for reading and writing bit streams
    - `bufio` - Utilities for buffered IO
    - `slices` - Utilities for manipulating slices

## Sources

- https://en.wikipedia.org/wiki/Huffman_coding
- https://en.wikipedia.org/wiki/LZ77_and_LZ78
- https://www.cs.helsinki.fi/u/tpkarkka/opetus/12k/dct/lecture07.pdf
- Sadakane, Kunihiko & Imai, Hiroshi. (2000). Improving the Speed of LZ77
  Compression by Hashing and Suffix Sorting. IEICE Transactions on Fundamentals
  of Electronics Communications and Computer Sciences. E83A. 
