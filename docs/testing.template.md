# Testing document

## Unit and integration testing

The project uses automated testing for everything except the command line
programs, which were tested manually. Test coverage is good but not 100 %. Most
of the lack of test coverage is due to not testing for IO error conditions. The
error handling code is very simple as it basically just halts the algorithm
immeaditely so I think there is not a pressing need to write unit tests for it.
Unit tests can be run by running `make test`.

An online test coverage report can be found
[here](https://codecov.io/gh/lassilaiho/compression-algorithms-tiralabra).

## Performance testing

Both algorithms were tested on a set of test files in `test/files`. The
collection of files comprises text and binary files in various formats and
sizes. There is also one file filled with random data read from `/dev/urandom`.

The following information was gathered for both algorithms: compression and
decompression speed, peak memory usage during compression and decompression as
well as compressed and uncompressed file size for each test file.

Running `make perf-report` generates performance reports for both programs. The
command requires GNU time to be available in the path. The output is written in
CSV format to `huffman-stas.csv` and `lz77-stats.csv` for Huffman coding and
LZ77, respectively.

By default each test file is compressed and decompressed five times and the
measured runtimes are averaged to improve accuracy. The number of iterations can
be specified by running `make perf-report ITERATIONS=x` where x is the number of
iterations to perform.

## Algorithm comparison

Huffman coding compression is considerably faster than LZ77. All test files
compressed faster with Huffman coding than with LZ77. Compressing data using
Huffman coding is essentially just performing array lookups whereas LZ77 must
search the window for matches to the current input position. Decompression
speeds in contrast are almost equal for both algorithms, with LZ77 being
slightly slower to decompress.

Both algorithms have a constant upper bound for memory usage for compression and
decompression, although memory usage varies based on the contents of the file.
LZ77 uses a lot more memory in the worst case than Huffman coding because of the
window.

LZ77 has a better compression ratio than Huffman coding for most of the test
files. However, the two large binary files, `random-data` and `regensburg.jpg`,
compress significantly worse with LZ77 than Huffman coding, altough they don't
compress very well with Huffman coding either. In fact, compressing
`random-data` with either algorithm actually increases the file size. This is
caused by not getting enough space savings from repeated data to outweight the
additional space required by metadata needed for decompression.

{{ .HuffmanTable }}

<h3 align="center">Huffman coding performance test results</h3>

{{ .LZ77Table }}

<h3 align="center">LZ77 performance test results</h3>
