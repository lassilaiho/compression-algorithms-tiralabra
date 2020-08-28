# User manual

## Building the project

Building the programs requires Go version 1.14 or later and GNU Make. Command
line programs can be built by running `make all` at the root of the repository.

## Command line programs

There are two command line programs, huffmancmd and lz77cmd. Huffmancmd
compresses and decompresses files using Huffman coding. Lz77cmd uses an
implementation of LZ77 compression algorithm.

Both programs have a uniform user interface:
```
program [flags] <input> <output>
```
where \<input> is the input file, \<output> is the file the output is written
to. The default action is to compress \<input> and write the output to
\<output>. Passing the `-d` flag switches the program to decompression mode. In
decompression mode \<input> must be a file compressed using the same program.
The decompressed file is written to \<output>. Both programs support the `-help`
flag which prints usage information to standard output.

## Performance report

Running `make perf-report` generates a performance report for both programs. The
output is written in CSV format to `huffman-stas.csv` and `lz77-stats.csv` for
Huffman coding and LZ77, respectively. The report includes average compression
and decompression time, peak memory usage as well as original and compressed
file sizes for each test file in `test/files` directory.

By default each test file is compressed and decompressed five times and the
measured runtimes are averaged to improve accuracy. The number of iterations can
be specified by running `make perf-report ITERATIONS=x` where x is the number of
iterations to perform.
