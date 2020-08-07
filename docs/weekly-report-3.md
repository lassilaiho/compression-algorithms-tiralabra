# Weekly report 3

This week I implemented the first version of LZ77 encoding and decoding. The
current encoding implementation is quite slow. I'll optimize it next week. I
also began replacing built-in and standard library algorithms and data
structures with my own implementations. I worked on the project for about 15
hours this week.

## Questions

Currently all of my encoding and decoding functions take as parameters standard
library interfaces for reading and writing bytes. In tests I use standard
library implementations for strings and arrays of bytes. In the command line
programs I use the standard library file IO type `os.File`. This is to reduce
memory usage by not reading the whole file into memory at once. Is this allowed?

## Next week

Next week I'll optimize the LZ77 encoding implementation.
