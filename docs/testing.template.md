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

There are two sets of test files: `test/files` and
`test/files/complexity-analysis`. The first includes files of various types and
sizes for comparing the overall performance of the algorithms. The latter file
set is used specifically for time and space complexity analysis in [the
implementation document](implementation.md). The file set includes only text
files with different sizes to more easily see the effect of file size on running
time and memory usage, because different types of files have very different
constant factors for time and space complexity.

The file set `test/files` comprises text and binary files in various formats and
sizes. There is also one file filled with random data read from `/dev/urandom`.

The following information was gathered for both algorithms: compression and
decompression speed, peak memory usage during compression and decompression as
well as compressed and uncompressed file size for each test file.

Running `make perf-report` generates performance reports for both programs. The
command requires GNU time to be available in the path. Data gathered from
`test/files` is  written in CSV format to `huffman-stas.csv` and
`lz77-stats.csv` for Huffman coding and LZ77, respectively. Data gathered from
`test/files/complexity-analysis` is written to `huffman-complexity-stats.csv`
and `lz77-complexity-stats.csv`.

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

### Huffman coding performance test results

{{ .HuffmanTable }}

### LZ77 performance test results

{{ .LZ77Table }}

## Sources

Test files are from the following websites:
  - http://sun.aei.polsl.pl/~sdeor/index.php?page=silesia
  - https://corpus.canterbury.ac.nz/descriptions/#cantrbry
  - http://www.gutenberg.org/ebooks/7000
  - https://pixabay.com/photos/regensburg-high-resolution-panorama-4423626/
