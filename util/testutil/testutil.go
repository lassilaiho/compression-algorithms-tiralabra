// Package testutil contains utility functions useful in unit tests.
package testutil

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

// ReadFile reads the file at filename and returns the contents of the file.
// Panics if reading fails.
func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}

// ExpectNil fails the test if a != nil.
func ExpectNil(t *testing.T, a interface{}) {
	if a != nil {
		t.Fatalf("expected nil, got %v", a)
	}
}

// ExpectEOF fails the test if err's error chain doesn't contain io.EOF.
func ExpectEOF(t *testing.T, err error) {
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

// Check fails if expected and found are not equal.
func Check(t *testing.T, expected, found interface{}) {
	if expected != found {
		t.Fatalf("expected %v, found %v", expected, found)
	}
}
