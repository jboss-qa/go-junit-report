// Package gtr defines a standard test report format and provides convenience
// methods to create and convert reports.
package gtr

import (
	"strings"
	"time"
)

// Result is the result of a test or benchmark.
type Result int

const (
	Unknown Result = iota
	Pass
	Fail
	Skip
)

func (r Result) String() string {
	switch r {
	case Unknown:
		return "UNKNOWN"
	case Pass:
		return "PASS"
	case Fail:
		return "FAIL"
	case Skip:
		return "SKIP"
	default:
		panic("invalid Result")
	}
}

// Report contains the build, test and/or benchmark results of a collection of
// packages.
type Report struct {
	Packages []Package
}

// IsSuccessful returns true if none of the packages in this report have build
// or runtime errors and all tests passed without failures or were skipped.
func (r *Report) IsSuccessful() bool {
	for _, pkg := range r.Packages {
		if pkg.BuildError.Name != "" || pkg.RunError.Name != "" {
			return false
		}
		for _, t := range pkg.Tests {
			if t.Result != Pass && t.Result != Skip {
				return false
			}
		}
	}
	return true
}

// Package contains build, test and/or benchmark results for a single package.
type Package struct {
	Name       string
	Timestamp  time.Time
	Duration   time.Duration
	Coverage   float64
	Output     []string
	Properties map[string]string

	Tests      []Test
	Benchmarks []Benchmark

	BuildError Error
	RunError   Error
}

// SetProperty stores a key/value property in the current package. If a
// property with the given key already exists, its old value will be
// overwritten with the given value.
func (p *Package) SetProperty(key, value string) {
	if p.Properties == nil {
		p.Properties = make(map[string]string)
	}
	p.Properties[key] = value
}

// Test contains the results of a single test.
type Test struct {
	Name     string
	Duration time.Duration
	Result   Result
	Level    int
	Output   []string
}

// Benchmark contains the results of a single benchmark.
type Benchmark struct {
	Name        string
	Result      Result
	Output      []string
	Iterations  int64
	NsPerOp     float64
	MBPerSec    float64
	BytesPerOp  int64
	AllocsPerOp int64
}

// Error contains details of a build or runtime error.
type Error struct {
	Name     string
	Duration time.Duration
	Cause    string
	Output   []string
}

// TrimPrefixSpaces trims the leading whitespace of the given line using the
// indentation level of the test. Printing logs in a Go test is typically
// prepended by blocks of 4 spaces to align it with the rest of the test
// output. TrimPrefixSpaces intends to only trim the whitespace added by the Go
// test command, without inadvertently trimming whitespace added by the test
// author.
func TrimPrefixSpaces(line string, indent int) string {
	// We only want to trim the whitespace prefix if it was part of the test
	// output. Test output is usually prefixed by a series of 4-space indents,
	// so we'll check for that to decide whether this output was likely to be
	// from a test.
	prefixLen := strings.IndexFunc(line, func(r rune) bool { return r != ' ' })
	if prefixLen%4 == 0 {
		// Use the subtest level to trim a consistenly sized prefix from the
		// output lines.
		for i := 0; i <= indent; i++ {
			line = strings.TrimPrefix(line, "    ")
		}
	}
	return strings.TrimPrefix(line, "\t")
}
