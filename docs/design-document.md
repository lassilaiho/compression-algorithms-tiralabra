# Design document

The objective of the project is to implement two data compression algorithms,
Huffman coding and LZ77, and compare their compression and decompression
performance, both in terms of time and memory usage, as well as compression
ratio. The project will be implemented in Go.

## Huffman coding

When compressing data using Huffman coding, the code table used for mapping
input bytes to code values must be constructed first. It is constructed from the
estimated or measured probabilities of every possible byte value in the input.
The best possible compression ratio is achieved by measuring the actual
frequencies of byte values in the input, so that's how I will implement it. The
code table will be prepended to the compressed data to allow decompression
without having to store the code table anywhere else. The actual compression and
decompression are rather simple to implement, basically just mapping byte values
to codes and vice versa.

Wikipedia describes two algorithms for finding the codes. The time complexity of
the first one is O(*n* log *n*) and the time complexity of the other one is
O(*n*), where *n* is the number of distinct byte values found in the input. Both
have space complexity of O(*n*). The algorithm with O(*n*) time complexity has
the additional restriction that the input must be sorted by the frequency of
each byte value. Because the measured byte values must be sorted, the time
complexity of code table construction is O(*n* log *n*) regardless of the choice
of algorithm. I plan to implement the first algorithm, since it seems easier to
implement. The time complexity of the actual compression and decompression is
O(*m*) where *m* is the size of the input. The time complexity of calculating
the byte value frequencies in the input is O(*m*) where *m* is the size of the
input.

Putting these together, the time complexity I try to achieve is O(*n* log *n* +
*m*) and space complexity is O(*n*) where *n* is the number of distinct byte
values in the input and *m* is the size of the input.

Implementing Huffman coding requires a minimum heap for the algorithm finding
the codes as well as a binary tree structure used for mapping codes to
decompressed byte values. Some simple additional data structures, such as
growable arrays and bit arrays, are probably also needed.

## LZ77

LZ77 works by converting the input data into a series of references to previous
occurences of the data, thereby avoiding multiple copies of the same data. This
is achieved by keeping track of recently encountered data, and replacing
occurrences with pointers to earlier occurrences. The size of the remembered
data is bounded by a constant independent of the size of the input. Therefore,
searching for matching byte strings has a constant maximum time independent of
the size of the input. Decompression is performed by copying already encountered
byte strings in place of pointers. The time complexity of compression and
decompression is therefore O(*m*) where *m* is the size of the input. The space
complexity of compression and decompression is O(1), since the size of the
remembered data is constant.

A simple implementation doesn't require complex data structures, growable arrays
and bit strings are sufficient. The searching of recent data can probably be
optimized from a simple implementation by using more efficient data structures,
such as hash tables. Different implementations must be benchmarked to find the
best solution.

## High level view of project structure

Each algorithm implementation will have a command line program which can be used
to compress and decompress arbitrary files. Input and output file paths will be
given as command line arguments. The algorithms will be implemented in separate
packages to make the project modular.

### Package structure

- `github.com/lassilaiho/compression-algorithms-tiralabra`
  - `cmd`
    - `huffman` - Huffman coding program
    - `lz77` - LZ77 program
  - `huffman` - Huffman coding implementation
  - `lz77` - LZ77 implementation

There may be additional packages for algorithms and data structures shared
between the two compression algorithms.

## Sources

- https://en.wikipedia.org/wiki/Huffman_coding
- https://en.wikipedia.org/wiki/LZ77_and_LZ78
- https://www.cs.helsinki.fi/u/tpkarkka/opetus/12k/dct/lecture07.pdf
