# Comparison of compression algorithms

This is a project for the course "Aineopintojen harjoitustyö: tietorakenteet ja
algoritmit" at the University of Helsinki.

## Building the project

The compression programs can be built by running `make all` at the root of the
repository.

`make test` executes unit tests and writes a test coverage report to
`cover.out`. The report can be viewed by running
```
go tool cover -html=cover.out
```

An online test coverage report can be found [here](codecov).

## Documentation

[Design document](docs/design-document.md)

## Weekly reports

- [Week 1](docs/weekly-report-1.md)
- [Week 2](docs/weekly-report-2.md)


<!-- Links -->
[codecov]: https://codecov.io/gh/lassilaiho/compression-algorithms-tiralabra
