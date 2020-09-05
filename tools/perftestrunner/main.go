package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	command        string
	inputDir       string
	workDir        string
	showHelp       bool
	iterationCount int
)

func init() {
	flag.StringVar(&command, "cmd", "", "(required) command to test")
	flag.StringVar(&inputDir, "dir", "", "(required) directory to read input files from")
	flag.BoolVar(&showHelp, "help", false, "print help message")
	flag.IntVar(&iterationCount, "iters", 5, "iteration count")
	flag.StringVar(&workDir, "workdir", ".", "working directory for temporary files")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			"usage: ", os.Args[0], " [flags]\n",
			"\n",
			"run <cmd> <iters> number of times for each file in <files>\n",
			"and its measure running time, memory usage and compression ratio\n",
			"\n",
		)
		flag.PrintDefaults()
	}
	flag.Parse()
}

type runResult struct {
	averageTime        time.Duration
	averageMaxMemUsage int
}

func (r *runResult) setAverage(iterations int) {
	r.averageTime /= time.Duration(iterations)
	r.averageMaxMemUsage /= iterations
}

func (r *runResult) averageTimeString() string {
	return strconv.FormatFloat(float64(r.averageTime)/float64(time.Second), 'f', 3, 64)
}

type testResult struct {
	file             string
	compression      runResult
	decompression    runResult
	compressedSize   int
	uncompressedSize int
}

func (r *testResult) spaceSavings() float64 {
	return 100 * (1 - float64(r.compressedSize)/float64(r.uncompressedSize))
}

type cmdTimeError struct {
	output string
	err    error
}

func (err *cmdTimeError) Error() string {
	return err.err.Error() + ", output: " + err.output
}

func timeCommand(result *runResult, cmd ...string) error {
	c := exec.Command("time", append([]string{"-f", "%e;%M"}, cmd...)...)
	output, err := c.CombinedOutput()
	if err != nil {
		return &cmdTimeError{
			output: string(output),
			err:    err,
		}
	}
	parts := bytes.Split(output, []byte{';'})
	if len(parts) < 2 {
		return fmt.Errorf("unrecognized command output: %s", string(output))
	}
	execTime, err := strconv.ParseFloat(string(parts[0]), 64)
	if err != nil {
		return err
	}
	memUsage, err := strconv.Atoi(string(bytes.Trim(parts[1], "\n")))
	if err != nil {
		return err
	}
	result.averageTime += time.Duration(execTime * float64(time.Second))
	result.averageMaxMemUsage += memUsage * 1024
	return nil
}

func runTest(command, inputFile, outputFile string, iterations int) (*testResult, error) {
	compressionOut := outputFile + ".compressed"
	decompressionOut := outputFile + ".decompressed"
	stat, err := os.Stat(inputFile)
	if err != nil {
		return nil, err
	}
	t := &testResult{
		file:             filepath.Base(inputFile),
		uncompressedSize: int(stat.Size()),
	}
	for i := 0; i < iterations; i++ {
		err := timeCommand(&t.compression, command, inputFile, compressionOut)
		if err != nil {
			return nil, err
		}
	}
	stat, err = os.Stat(compressionOut)
	if err != nil {
		return nil, err
	}
	t.compressedSize = int(stat.Size())
	t.compression.setAverage(iterations)
	for i := 0; i < iterations; i++ {
		err := timeCommand(&t.decompression, command, "-d", compressionOut, decompressionOut)
		if err != nil {
			return nil, err
		}
	}
	t.decompression.setAverage(iterations)
	return t, nil
}

func printResults(results []*testResult) {
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{
		"File",
		"Average compression execution time (s)",
		"Average compression peak memory usage (B)",
		"Average decompression execution time (s)",
		"Average decompression peak memory usage (B)",
		"Uncompressed size (B)",
		"Compressed size (B)",
		"Space savings (%)",
	})
	for _, result := range results {
		if result == nil {
			continue
		}
		w.Write([]string{
			result.file,
			result.compression.averageTimeString(),
			strconv.Itoa(result.compression.averageMaxMemUsage),
			result.decompression.averageTimeString(),
			strconv.Itoa(result.decompression.averageMaxMemUsage),
			strconv.Itoa(result.uncompressedSize),
			strconv.Itoa(result.compressedSize),
			strconv.FormatFloat(result.spaceSavings(), 'f', 2, 64),
		})
	}
	w.Flush()
}

func run() error {
	if showHelp {
		flag.Usage()
		return nil
	}
	if command == "" {
		return errors.New("-cmd is required")
	}
	if inputDir == "" {
		return errors.New("-dir is required")
	}
	fmt.Fprintln(os.Stderr, "Testing ", command)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}
	inputFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}
	results := make([]*testResult, len(inputFiles))
	for i, inputFile := range inputFiles {
		if inputFile.IsDir() {
			continue
		}
		fmt.Fprintln(os.Stderr, "  Processing", inputFile.Name())

		inputPath := filepath.Join(inputDir, inputFile.Name())
		outputPath := filepath.Join(workDir, inputFile.Name())
		results[i], err = runTest(command, inputPath, outputPath, iterationCount)
		if err != nil {
			return err
		}
	}
	printResults(results)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
