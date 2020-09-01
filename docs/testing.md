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

| File | Average compression execution time (s) | Avarege compression peak memory usage (B) | Average decompression execution time (s) | Avarege decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.01 | 2093056 | 0.01 | 2013593 | 152089 | 87789 | 0.42277876769523104 |
| asyoulik.txt | 0.01 | 2097971 | 0.01 | 2013593 | 125179 | 75899 | 0.39367625560197794 |
| cp.html | 0 | 2097971 | 0 | 2020966 | 24603 | 16314 | 0.33691013291062066 |
| fields.c | 0 | 2102067 | 0 | 2000486 | 11150 | 7147 | 0.35901345291479825 |
| grammar.lsp | 0 | 2112716 | 0 | 2011955 | 3721 | 2273 | 0.389142703574308 |
| kalevala.txt | 0.05 | 2106163 | 0.05 | 2011955 | 658369 | 371751 | 0.43534552811569194 |
| kennedy.xls | 0.052 | 2162688 | 0.05 | 2026700 | 1029744 | 462860 | 0.5505096412312187 |
| lcet10.txt | 0.032 | 2114355 | 0.03 | 2009497 | 426754 | 250677 | 0.4125960155030767 |
| plrabn12.txt | 0.04 | 2120908 | 0.04 | 1999667 | 481861 | 275694 | 0.42785575093232275 |
| ptt5 | 0.02 | 2103705 | 0.01 | 2016870 | 513216 | 106758 | 0.7919823232323232 |
| random-data | 7.53 | 2201190 | 8.267999999 | 2034892 | 52428800 | 52429128 | -0.000006256103515678291 |
| regensburg.jpg | 5.088 | 2186444 | 5.586 | 2046361 | 30771782 | 30760865 | 0.00035477308398979 |
| sum | 0.016 | 2136473 | 0 | 2011136 | 38240 | 25972 | 0.32081589958158996 |
| xargs.1 | 0 | 2118451 | 0 | 2005401 | 4227 | 2702 | 0.3607759640406908 |

<h3 align="center">Huffman coding performance test results</h3>

Huffman coding is less demanding in terms of computation time and memory usage
compared to LZ77. Memory usage is about constant.

Huffman coding compression is considerably faster than LZ77. All test files
compressed faster with Huffman coding than with LZ77. Compressing data using
Huffman coding is essentially just performing array lookups whereas LZ77 must
search the window for matches to the current input position. Decompression
speeds in contrast are almost equal for both algorithms, with LZ77 being
slightly slower to decompress.

Both algorithms have a constant upper bound for memory usage for compression and
decompression, although memory usage varies based on the contents of the file.
LZ77 uses more memory than Huffman coding because of the window.

LZ77 has a better compression ratio than Huffman coding for most of the test
files. However, the two large binary files, `random-data` and `regensburg.jpg`,
compress significantly worse with LZ77 than Huffman coding, altough they don't
compress very well with Huffman coding either. In fact, compressing
`random-data` with either algorithm actually increases the file size. This is
caused by not getting enough space savings from repeated data to outweight the
additional space required by metadata needed for decompression.

| File | Average compression execution time (s) | Avarege compression peak memory usage (B) | Average decompression execution time (s) | Avarege decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.05 | 5029888 | 0.01 | 2111078 | 152089 | 75273 | 0.5050726877025952 |
| asyoulik.txt | 0.046 | 4536729 | 0.01 | 2088960 | 125179 | 67405 | 0.46153108748272476 |
| cp.html | 0.004 | 2775449 | 0 | 2100428 | 24603 | 11423 | 0.5357070275982604 |
| fields.c | 0 | 2413363 | 0 | 2088960 | 11150 | 4004 | 0.6408968609865471 |
| grammar.lsp | 0 | 2192179 | 0 | 2110259 | 3721 | 1581 | 0.5751142166084386 |
| kalevala.txt | 0.236 | 7448985 | 0.05 | 2104524 | 658369 | 300326 | 0.5438333214352438 |
| kennedy.xls | 0.328 | 12558336 | 0.06 | 2112716 | 1029744 | 349789 | 0.6603146024643018 |
| lcet10.txt | 0.156 | 7221248 | 0.032 | 2093875 | 426754 | 205940 | 0.5174269016810622 |
| plrabn12.txt | 0.196 | 7249920 | 0.042 | 2079948 | 481861 | 270760 | 0.4380952183305974 |
| ptt5 | 1.862 | 7340851 | 0.024 | 2084864 | 513216 | 116575 | 0.7728539250529991 |
| random-data | 58.924 | 18959564 | 9.786 | 2085683 | 52428800 | 58595928 | -0.11762863159179693 |
| regensburg.jpg | 34.451999999 | 18388582 | 5.78 | 2107801 | 30771782 | 34197038 | -0.11131159059946549 |
| sum | 0.02 | 3485696 | 0 | 2087321 | 38240 | 18335 | 0.5205282426778243 |
| xargs.1 | 0 | 2179891 | 0 | 2086502 | 4227 | 2165 | 0.48781641826354394 |

<h3 align="center">LZ77 performance test results</h3>
