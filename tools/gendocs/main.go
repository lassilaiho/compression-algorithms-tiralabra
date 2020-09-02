package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var showHelp bool

func init() {
	flag.BoolVar(&showHelp, "help", false, "print help message")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			"usage: ", os.Args[0], " [flags] <template_dir> <data_dir>\n",
			"\n",
			"Processes templates in directory <template_dir> with data\n",
			"in directory <data_dir>.",
		)
		flag.PrintDefaults()
	}
	flag.Parse()
}

type perfData struct {
	Huffman, LZ77 [][]string

	HuffmanTable, LZ77Table string
}

func (d *perfData) genTables() {
	d.HuffmanTable = formatTable(d.Huffman)
	d.LZ77Table = formatTable(d.LZ77)
}

func repeat(s string, n int) []string {
	xs := make([]string, n)
	for i := range xs {
		xs[i] = s
	}
	return xs
}

func formatTableRow(b *strings.Builder, row []string) {
	for _, cell := range row {
		b.WriteString("| ")
		b.WriteString(cell)
		b.WriteByte(' ')
	}
	b.WriteString("|\n")
}

func formatTable(data [][]string) string {
	if len(data) == 0 {
		return ""
	}
	var b strings.Builder
	formatTableRow(&b, data[0])
	formatTableRow(&b, repeat("----", len(data[0])))
	for _, row := range data[1:] {
		formatTableRow(&b, row)
	}
	return b.String()
}

func readData(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return csv.NewReader(f).ReadAll()
}

func run() error {
	if showHelp {
		flag.Usage()
		return nil
	}
	if flag.NArg() != 2 {
		return errors.New("expected 2 argument")
	}
	tmpldir := flag.Arg(0)
	datadir := flag.Arg(1)
	tmpl, err := template.ParseGlob(filepath.Join(tmpldir, "*.template.md"))
	if err != nil {
		return err
	}
	var data perfData
	data.Huffman, err = readData(filepath.Join(datadir, "huffman-stats.csv"))
	if err != nil {
		return err
	}
	data.LZ77, err = readData(filepath.Join(datadir, "lz77-stats.csv"))
	if err != nil {
		return err
	}
	data.genTables()
	for _, t := range tmpl.Templates() {
		name := strings.TrimSuffix(t.Name(), ".template.md") + ".md"
		f, err := os.Create(filepath.Join(tmpldir, name))
		if err != nil {
			return err
		}
		defer f.Close()
		if err := t.Execute(f, &data); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
