package testutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// TestSpec represents a single test case from ledger test files
type TestSpec struct {
	Name           string
	InputData      string
	Command        string
	ExpectedOutput string
	ShouldFail     bool
	FileName       string
	LineNumber     int
}

// SpecParser parses ledger test files into TestSpec structures
type SpecParser struct {
	reader *bufio.Scanner
	file   string
	line   int
}

// NewSpecParser creates a new parser for ledger test files
func NewSpecParser(r io.Reader, filename string) *SpecParser {
	return &SpecParser{
		reader: bufio.NewScanner(r),
		file:   filename,
		line:   0,
	}
}

// ParseFile parses a ledger test file and returns all test specs
func ParseFile(filename string) ([]*TestSpec, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open test file: %w", err)
	}
	defer file.Close()

	parser := NewSpecParser(file, filename)
	return parser.ParseAll()
}

// ParseAll parses all test specs from the reader
func (p *SpecParser) ParseAll() ([]*TestSpec, error) {
	var specs []*TestSpec
	var currentInput strings.Builder
	inInput := true

	for p.reader.Scan() {
		p.line++
		line := p.reader.Text()

		// Check for test block start
		if strings.HasPrefix(line, "test ") {
			// Save any pending input data
			inputData := strings.TrimSpace(currentInput.String())
			currentInput.Reset()
			inInput = false

			// Parse test block
			spec, err := p.parseTestBlock(line, inputData)
			if err != nil {
				return nil, fmt.Errorf("error parsing test at line %d: %w", p.line, err)
			}
			if spec != nil {
				specs = append(specs, spec)
			}
		} else if inInput {
			// Accumulate input data before first test
			currentInput.WriteString(line)
			currentInput.WriteString("\n")
		}
	}

	if err := p.reader.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return specs, nil
}

// parseTestBlock parses a single test block
func (p *SpecParser) parseTestBlock(startLine string, inputData string) (*TestSpec, error) {
	// Extract test name/command
	testCmd := strings.TrimPrefix(startLine, "test ")
	
	spec := &TestSpec{
		Name:       testCmd,
		Command:    testCmd,
		InputData:  inputData,
		FileName:   p.file,
		LineNumber: p.line,
	}

	// Read expected output until "end test" or "__ERROR__"
	var output strings.Builder
	for p.reader.Scan() {
		p.line++
		line := p.reader.Text()

		if line == "end test" {
			spec.ExpectedOutput = strings.TrimSpace(output.String())
			return spec, nil
		}

		if line == "__ERROR__" {
			spec.ShouldFail = true
			continue
		}

		output.WriteString(line)
		output.WriteString("\n")
	}

	return nil, fmt.Errorf("unexpected end of file in test block")
}

// CompareOutput compares actual output with expected output
func CompareOutput(expected, actual string) bool {
	// Normalize line endings and trim spaces
	expected = normalizeOutput(expected)
	actual = normalizeOutput(actual)
	return expected == actual
}

// normalizeOutput normalizes output for comparison
func normalizeOutput(s string) string {
	// Trim trailing whitespace from each line
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	
	// Join and trim overall
	result := strings.Join(lines, "\n")
	return strings.TrimSpace(result)
}