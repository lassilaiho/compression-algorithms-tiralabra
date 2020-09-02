# Weekly report 2

This week I implemented Huffman coding using standard library algorithms. The
implementation is unit tested and the test coverage is 84.8 % according to `go
test`. I also set up GitHub Actions to automatically test each commit and report
test coverage to codecov.io. I worked on the project for about 11 hours this
week.

## Questions

Which parts of the standard library are allowed? For example, are IO utilities,
such as package "bufio" and common IO interfaces in package "io" allowed when
used as part of the algorithm implementation? Error handling utilities in
packages "error" and "fmt"? Built-in hash maps are probably not allowed. What
about slices and associated built-in functions such as append and copy?

Is code generation allowed for generating multiple copies of the same algorithm
or data structure for different types? Escpecially if slices and associated
functions are not allowed, it would quickly become tedious to copy and paste
multiple almost identical implementations.

## Next week

Next week I will implement LZ77 using the standard library and possibly start
implementing replacements for the parts of the standard library I use in the
project.
