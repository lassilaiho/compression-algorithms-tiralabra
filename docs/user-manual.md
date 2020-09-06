# User manual

## Building the project

You can download pre-built binaries from the [releases
page](https://github.com/lassilaiho/compression-algorithms-tiralabra/releases/tag/final-version)
or you can build the programs yourself.

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
where \<input> is the input file and \<output> is the file the output is written
to. The default action is to compress \<input> and write the output to
\<output>. Passing the `-d` flag switches the program to decompression mode. In
decompression mode \<input> must be a file compressed using the same program.
The decompressed file is written to \<output>. Both programs support the `-help`
flag which prints usage information.
