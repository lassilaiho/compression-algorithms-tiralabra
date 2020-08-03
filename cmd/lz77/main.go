package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lassilaiho/compression-algorithms-tiralabra/lz77"
)

var decompress bool
var showHelp bool

func init() {
	flag.BoolVar(&decompress, "d", false, "decompress instead of compressing")
	flag.BoolVar(&showHelp, "help", false, "print help message")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr,
			"usage:", os.Args[0], "[flags] <input file> <output file>")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr,
			"compress <input file> and write the output to <output file>")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func printErrln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func run() error {
	if showHelp {
		flag.Usage()
		return nil
	}
	if flag.NArg() != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", flag.NArg())
	}
	inputFile, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}
	defer inputFile.Close()
	outputFile, err := os.Create(flag.Arg(1))
	if err != nil {
		return err
	}
	defer outputFile.Close()
	if decompress {
		return lz77.Decode(inputFile, outputFile)
	}
	return lz77.Encode(inputFile, outputFile)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
