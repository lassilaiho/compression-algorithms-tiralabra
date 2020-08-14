# Comparison of compression algorithms

This is a project for the course "Aineopintojen harjoitustyö: tietorakenteet ja
algoritmit" at the University of Helsinki.

## Building the project

The compression programs can be built by running `make all` at the root of the
repository.

`make lint` runs `go vet` and `golint` on the project (requires golint to be
installed and found in the path).

`make test` executes unit tests and writes a test coverage report to
`cover.out`. The report can be viewed by running
```
go tool cover -html=cover.out
```

An online test coverage report can be found [here](https://codecov.io/gh/lassilaiho/compression-algorithms-tiralabra).

## Documentation

[Design document](docs/design-document.md)

[Implementation document](docs/implementation.md)

[Testing document](docs/testing.md)

## Weekly reports

- [Week 1](docs/weekly-report-1.md)
- [Week 2](docs/weekly-report-2.md)
- [Week 3](docs/weekly-report-3.md)
