// gendocs generates documentation from templates. The templates are MarkDown
// files ending in .template.md and are processed as templates using the
// standard library package "text/templates".
package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var showHelp bool

// Define command line parameters.
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

// perfData is used as the context for templates.
type perfData struct {
	// Tables as MarkDown formatted strings
	HuffmanTable, LZ77Table string

	// Paths to test result files
	huffmanFile           string
	lz77File              string
	huffmanComplexityFile string
	lz77ComplexityFile    string

	graphDir   string // Directory where generated graphs are stored
	linkPrefix string // Prefix used for MarkDown links
}

// loadPerfData loads performance data from dataDir and prepares it for use with
// templates. Generated graphs are stored in graphDir and links are prefixed
// with linkPrefix.
func loadPerfData(dataDir, graphDir, linkPrefix string) (*perfData, error) {
	data := &perfData{
		huffmanFile:           filepath.Join(dataDir, "huffman-stats.csv"),
		lz77File:              filepath.Join(dataDir, "lz77-stats.csv"),
		huffmanComplexityFile: filepath.Join(dataDir, "huffman-complexity-stats.csv"),
		lz77ComplexityFile:    filepath.Join(dataDir, "lz77-complexity-stats.csv"),
		graphDir:              graphDir,
		linkPrefix:            linkPrefix,
	}
	huffmanData, err := readData(data.huffmanFile)
	if err != nil {
		return nil, err
	}
	lz77Data, err := readData(data.lz77File)
	if err != nil {
		return nil, err
	}
	data.HuffmanTable = formatTable(huffmanData)
	data.LZ77Table = formatTable(lz77Data)
	return data, nil
}

// Gnuplot generates a graph using gnuplot. The file name is graphName suffixed
// with the file format. commands is the script used to generate the graph. Some
// default settings and variables are set automatically. Gnuplot returns a
// MarkDown image reference on success and an error value on failure.
func (d *perfData) Gnuplot(graphName, commands string) (string, error) {
	graphFileName := graphName + ".png"
	graphFilePath := filepath.Join(d.graphDir, graphFileName)
	var script bytes.Buffer
	fmt.Fprintf(&script, `
set datafile separator comma
set grid
set output '%s'
set term png size 700,480
`, graphFilePath)
	script.WriteString(commands)

	cmd := exec.Command("gnuplot",
		"-e", "Huffman='"+d.huffmanFile+"'",
		"-e", "LZ77='"+d.lz77File+"'",
		"-e", "HuffmanComplexity='"+d.huffmanComplexityFile+"'",
		"-e", "LZ77Complexity='"+d.lz77ComplexityFile+"'",
		"-",
	)
	cmd.Stdin = &script
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run gnuplot: %s", string(output))
	}
	return imageRef(graphName, path.Join(d.linkPrefix, graphFileName)), nil
}

// Graphviz generates a graph using Graphviz. The file name is graphName
// suffixed with the file format. commands is the script used to generate the
// graph. Graphviz returns a MarkDown image reference on success and an error
// value on failure.
func (d *perfData) Graphviz(graphName, commands string) (string, error) {
	graphFileName := graphName + ".png"
	graphFilePath := filepath.Join(d.graphDir, graphFileName)
	cmd := exec.Command("dot", "-Tpng", "-o", graphFilePath)
	cmd.Stdin = strings.NewReader(commands)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run dot: %s", string(output))
	}
	return imageRef(graphName, path.Join(d.linkPrefix, graphFileName)), nil
}

// generateDocument generates a document using template t and data from d. The
// file is written to outFile.
func (d *perfData) generateDocument(t *template.Template, outFile string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, `<!--
    This file was automatically generated from "%s".
    Any changes will be overwritten when the file is regenerated.
-->
`, t.Name())
	if err != nil {
		return err
	}
	return t.Execute(f, d)
}

// imageRef returns a MarkDown image reference with name and path.
func imageRef(name, path string) string {
	return "![" + name + "](" + path + ")"
}

// repeat returns a slice with n copies of s.
func repeat(s string, n int) []string {
	xs := make([]string, n)
	for i := range xs {
		xs[i] = s
	}
	return xs
}

// formatTableRow formats row as MarkDown and writes the result to b.
func formatTableRow(b *strings.Builder, row []string) {
	for _, cell := range row {
		b.WriteString("| ")
		b.WriteString(cell)
		b.WriteByte(' ')
	}
	b.WriteString("|\n")
}

// formatTable formats data as a MarkDown table. If data has at least one row,
// the first row is considered to be a header row.
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

// readData reads CSV data from file.
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
		return errors.New("expected 2 arguments")
	}
	tmpldir := flag.Arg(0)
	dataDir := flag.Arg(1)
	tmpl, err := template.ParseGlob(filepath.Join(tmpldir, "*.template.md"))
	if err != nil {
		return err
	}
	imageDir := filepath.Join(tmpldir, "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return err
	}
	data, err := loadPerfData(dataDir, imageDir, "images")
	if err != nil {
		return err
	}
	for _, t := range tmpl.Templates() {
		name := strings.TrimSuffix(t.Name(), ".template.md") + ".md"
		err := data.generateDocument(t, filepath.Join(tmpldir, name))
		if err != nil {
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
