<!--
    This file was automatically generated from "testing.template.md".
    Any changes will be overwritten when the file is regenerated.
-->
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

| File | Average compression execution time (s) | Average compression peak memory usage (B) | Average decompression execution time (s) | Average decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.012 | 5830246 | 0.010 | 5764710 | 152089 | 87789 | 42.28 |
| asyoulik.txt | 0.010 | 5833523 | 0.010 | 5763072 | 125179 | 75899 | 39.37 |
| cp.html | 0.000 | 6666649 | 0.000 | 5762252 | 24603 | 16314 | 33.69 |
| dickens | 0.760 | 5831884 | 0.758 | 5763891 | 10192446 | 5826054 | 42.84 |
| fields.c | 0.004 | 5830246 | 0.000 | 5763072 | 11150 | 7147 | 35.90 |
| grammar.lsp | 0.000 | 5832704 | 0.000 | 5763891 | 3721 | 2273 | 38.91 |
| kalevala.txt | 0.048 | 5833523 | 0.040 | 6181683 | 658369 | 371751 | 43.53 |
| kennedy.xls | 0.050 | 5832704 | 0.042 | 5764710 | 1029744 | 462860 | 55.05 |
| lcet10.txt | 0.030 | 5833523 | 0.030 | 5765529 | 426754 | 250677 | 41.26 |
| mozilla | 4.540 | 5835161 | 4.738 | 5766348 | 51220480 | 39977192 | 21.95 |
| mr | 0.672 | 5832704 | 0.640 | 6183321 | 9970564 | 4623203 | 53.63 |
| nci | 1.588 | 6665011 | 1.334 | 6178406 | 33553445 | 10223991 | 69.53 |
| ooffice | 0.614 | 5834342 | 0.646 | 5764710 | 6152192 | 5124540 | 16.70 |
| osdb | 1.074 | 6248038 | 1.128 | 6603571 | 10085684 | 8342202 | 17.29 |
| plrabn12.txt | 0.036 | 6251315 | 0.032 | 5763072 | 481861 | 275694 | 42.79 |
| ptt5 | 0.018 | 5834342 | 0.010 | 5764710 | 513216 | 106758 | 79.20 |
| random-data | 1.214 | 6248857 | 1.276 | 5779456 | 10485760 | 10486088 | -0.00 |
| regensburg.jpg | 3.852 | 5835980 | 4.194 | 6180044 | 30771782 | 30760865 | 0.04 |
| reymont | 0.488 | 5832704 | 0.486 | 5766348 | 6627202 | 4031623 | 39.17 |
| samba | 1.826 | 5832704 | 1.870 | 6184960 | 21606400 | 16546872 | 23.42 |
| sao | 0.850 | 5833523 | 0.904 | 5781094 | 7251944 | 6843596 | 5.63 |
| sum | 0.002 | 6248038 | 0.000 | 5780275 | 38240 | 25972 | 32.08 |
| webster | 3.276 | 6246400 | 3.272 | 5763072 | 41458703 | 25929010 | 37.46 |
| x-ray | 0.824 | 6666649 | 0.846 | 5764710 | 8474240 | 7021658 | 17.14 |
| xargs.1 | 0.004 | 5831065 | 0.000 | 5762252 | 4227 | 2702 | 36.08 |
| xml | 0.412 | 6249676 | 0.414 | 5763072 | 5345280 | 3711199 | 30.57 |


### LZ77 performance test results

| File | Average compression execution time (s) | Average compression peak memory usage (B) | Average decompression execution time (s) | Average decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.066 | 8057651 | 0.010 | 5763072 | 152089 | 75273 | 50.51 |
| asyoulik.txt | 0.052 | 8015872 | 0.010 | 6178406 | 125179 | 67405 | 46.15 |
| cp.html | 0.010 | 6231654 | 0.000 | 6183321 | 24603 | 11423 | 53.57 |
| dickens | 4.368 | 10509516 | 0.838 | 5771264 | 10192446 | 5391894 | 47.10 |
| fields.c | 0.002 | 5785190 | 0.000 | 6184140 | 11150 | 4004 | 64.09 |
| grammar.lsp | 0.000 | 6191513 | 0.000 | 5769625 | 3721 | 1581 | 57.51 |
| kalevala.txt | 0.256 | 10412851 | 0.048 | 6619136 | 658369 | 300326 | 54.38 |
| kennedy.xls | 0.364 | 10548838 | 0.050 | 5781913 | 1029744 | 349789 | 66.03 |
| lcet10.txt | 0.174 | 10330112 | 0.030 | 5786009 | 426754 | 205940 | 51.74 |
| mozilla | 36.218 | 22565683 | 3.980 | 5784371 | 51220480 | 24782036 | 51.62 |
| mr | 8.032 | 10734796 | 0.772 | 5781913 | 9970564 | 5364448 | 46.20 |
| nci | 17.732 | 11123916 | 1.514 | 5801574 | 33553445 | 7565462 | 77.45 |
| ooffice | 4.034 | 21764505 | 0.594 | 6197248 | 6152192 | 3845486 | 37.49 |
| osdb | 7.332 | 21988966 | 1.122 | 5799936 | 10085684 | 7474200 | 25.89 |
| plrabn12.txt | 0.220 | 10412851 | 0.040 | 5782732 | 481861 | 270760 | 43.81 |
| ptt5 | 1.468 | 10388275 | 0.020 | 5781913 | 513216 | 116575 | 77.29 |
| random-data | 13.490 | 21649817 | 1.716 | 5783552 | 10485760 | 11719209 | -11.76 |
| regensburg.jpg | 38.472 | 20733952 | 5.008 | 6198067 | 30771782 | 34197038 | -11.13 |
| reymont | 2.632 | 13103104 | 0.452 | 6200524 | 6627202 | 2652153 | 59.98 |
| samba | 11.256 | 24024678 | 1.440 | 5783552 | 21606400 | 8422677 | 61.02 |
| sao | 6.614 | 19588710 | 0.914 | 6201344 | 7251944 | 6107625 | 15.78 |
| sum | 0.022 | 5855641 | 0.000 | 5783552 | 38240 | 18335 | 52.05 |
| webster | 16.212 | 10540646 | 3.042 | 5783552 | 41458703 | 17736878 | 57.22 |
| x-ray | 5.916 | 10532454 | 1.074 | 5781094 | 8474240 | 8193531 | 3.31 |
| xargs.1 | 0.000 | 6226739 | 0.000 | 5781094 | 4227 | 2165 | 48.78 |
| xml | 1.724 | 10528358 | 0.278 | 5781094 | 5345280 | 1398194 | 73.84 |


## Sources

Test files are from the following websites:
  - http://sun.aei.polsl.pl/~sdeor/index.php?page=silesia
  - https://corpus.canterbury.ac.nz/descriptions/#cantrbry
  - http://www.gutenberg.org/ebooks/7000
  - https://pixabay.com/photos/regensburg-high-resolution-panorama-4423626/
