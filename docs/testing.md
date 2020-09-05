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

| File | Average compression execution time (s) | Average compression peak memory usage (B) | Average decompression execution time (s) | Average decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.014 | 2106982 | 0.012 | 2010316 | 152089 | 87789 | 0.42277876769523104 |
| asyoulik.txt | 0.01 | 2094694 | 0.01 | 2001305 | 125179 | 75899 | 0.39367625560197794 |
| cp.html | 0 | 2118451 | 0 | 2011136 | 24603 | 16314 | 0.33691013291062066 |
| fields.c | 0 | 2108620 | 0 | 2008678 | 11150 | 7147 | 0.35901345291479825 |
| grammar.lsp | 0 | 2093875 | 0 | 1998028 | 3721 | 2273 | 0.389142703574308 |
| kalevala.txt | 0.052 | 2101248 | 0.05 | 2002124 | 658369 | 371751 | 0.43534552811569194 |
| kennedy.xls | 0.06 | 2138112 | 0.056 | 2051276 | 1029744 | 462860 | 0.5505096412312187 |
| lcet10.txt | 0.04 | 2113536 | 0.036 | 1992294 | 426754 | 250677 | 0.4125960155030767 |
| plrabn12.txt | 0.04 | 2096332 | 0.042 | 2002124 | 481861 | 275694 | 0.42785575093232275 |
| ptt5 | 0.02 | 2119270 | 0.014 | 2022604 | 513216 | 106758 | 0.7919823232323232 |
| sum | 0 | 2154496 | 0 | 2056192 | 38240 | 25972 | 0.32081589958158996 |
| xargs.1 | 0 | 2101248 | 0 | 1993932 | 4227 | 2702 | 0.3607759640406908 |


<h3 align="center">Huffman coding performance test results</h3>

| File | Average compression execution time (s) | Average compression peak memory usage (B) | Average decompression execution time (s) | Average decompression peak memory usage (B) | Uncompressed size (B) | Compressed size (B) | Space savings (%) |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| alice29.txt | 0.06 | 5055283 | 0.01 | 2097152 | 152089 | 75273 | 0.5050726877025952 |
| asyoulik.txt | 0.052 | 4429414 | 0.01 | 2089779 | 125179 | 67405 | 0.46153108748272476 |
| cp.html | 0.01 | 2822144 | 0 | 2096332 | 24603 | 11423 | 0.5357070275982604 |
| fields.c | 0 | 2373222 | 0 | 2080768 | 11150 | 4004 | 0.6408968609865471 |
| grammar.lsp | 0 | 2201190 | 0 | 2090598 | 3721 | 1581 | 0.5751142166084386 |
| kalevala.txt | 0.252 | 7498956 | 0.054 | 2143027 | 658369 | 300326 | 0.5438333214352438 |
| kennedy.xls | 0.36 | 7853670 | 0.06 | 2088140 | 1029744 | 349789 | 0.6603146024643018 |
| lcet10.txt | 0.174 | 7032012 | 0.03 | 2082406 | 426754 | 205940 | 0.5174269016810622 |
| plrabn12.txt | 0.232 | 7344128 | 0.046 | 2079948 | 481861 | 270760 | 0.4380952183305974 |
| ptt5 | 1.914 | 7440793 | 0.02 | 2120089 | 513216 | 116575 | 0.7728539250529991 |
| sum | 0.022 | 3504537 | 0 | 2085683 | 38240 | 18335 | 0.5205282426778243 |
| xargs.1 | 0 | 2203648 | 0 | 2066022 | 4227 | 2165 | 0.48781641826354394 |


<h3 align="center">LZ77 performance test results</h3>
