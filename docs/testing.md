# Testing document

The project uses automated testing for everything except the command line
programs, which were tested manually. Test coverage is good but not 100 %. Most of
the lack of test coverage is due to not testing for IO error conditions. The
error handling code is very simple as it basically just halts the algorithm
immeaditely so I think there is not a pressing need to write unit tests for it.

An online test coverage report can be found
[here](https://codecov.io/gh/lassilaiho/compression-algorithms-tiralabra).
