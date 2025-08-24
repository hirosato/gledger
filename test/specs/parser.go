package specs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TestSpec represents a single test case from a ledger test file
type TestSpec struct {
	Name        string   // Test name (e.g., "accounts", "balance --flat")
	Command     string   // Full command to run
	InputData   []string // Lines of input ledger data (before "test" directive)
	ExpectedOut []string // Expected output lines
	ExpectedErr []string // Expected error output (if any)
	SourceFile  string   // Path to the source test file
	LineNumber  int      // Line number where test starts
}

// TestFile represents a complete test file containing multiple test specs
type TestFile struct {
	Path      string     // Path to the test file
	InputData []string   // Shared input data for all tests in file
	Tests     []TestSpec // Individual test cases
}

// Parser handles parsing of ledger test files
type Parser struct {
	verbose bool
}

// NewParser creates a new test spec parser
func NewParser(verbose bool) *Parser {
	return &Parser{
		verbose: verbose,
	}
}

// ParseFile parses a single ledger test file
func (p *Parser) ParseFile(path string) (*TestFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open test file: %w", err)
	}
	defer file.Close()

	return p.parse(file, path)
}

// parse performs the actual parsing of the test file
func (p *Parser) parse(r io.Reader, sourcePath string) (*TestFile, error) {
	scanner := bufio.NewScanner(r)
	result := &TestFile{
		Path:  sourcePath,
		Tests: []TestSpec{},
	}

	var currentTest *TestSpec
	var lineNum int
	inTest := false
	collectingInput := true

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for test directive
		if strings.HasPrefix(line, "test ") {
			// Save any previous test
			if currentTest != nil {
				result.Tests = append(result.Tests, *currentTest)
			}

			// Start new test
			testName := strings.TrimPrefix(line, "test ")
			currentTest = &TestSpec{
				Name:        testName,
				Command:     p.parseCommand(testName),
				InputData:   result.InputData, // Share input data
				ExpectedOut: []string{},
				ExpectedErr: []string{},
				SourceFile:  sourcePath,
				LineNumber:  lineNum,
			}
			inTest = true
			collectingInput = false
			continue
		}

		// Check for end of test
		if line == "end test" {
			if currentTest != nil {
				result.Tests = append(result.Tests, *currentTest)
				currentTest = nil
			}
			inTest = false
			continue
		}

		// Collect lines based on current state
		if inTest && currentTest != nil {
			// Inside a test block - collect expected output
			currentTest.ExpectedOut = append(currentTest.ExpectedOut, line)
		} else if collectingInput && !inTest {
			// Before any test directive - collect input data (skip empty trailing lines)
			if line != "" || len(result.InputData) > 0 {
				result.InputData = append(result.InputData, line)
			}
		}
	}

	// Handle any remaining test
	if currentTest != nil {
		result.Tests = append(result.Tests, *currentTest)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading test file: %w", err)
	}

	return result, nil
}

// parseCommand extracts the actual ledger command from the test name
func (p *Parser) parseCommand(testName string) string {
	// The test name often IS the command (e.g., "accounts", "balance --flat")
	// We'll prepend "ledger" and add any necessary flags
	parts := strings.Fields(testName)
	if len(parts) == 0 {
		return ""
	}

	// Build the command
	cmd := parts[0]
	args := strings.Join(parts[1:], " ")
	
	// Return the command and arguments
	if args != "" {
		return fmt.Sprintf("%s %s", cmd, args)
	}
	return cmd
}

// ParseDirectory parses all test files in a directory
func (p *Parser) ParseDirectory(dir string) ([]*TestFile, error) {
	var testFiles []*TestFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .test files
		if !strings.HasSuffix(path, ".test") {
			return nil
		}

		if p.verbose {
			fmt.Printf("Parsing test file: %s\n", path)
		}

		testFile, err := p.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		testFiles = append(testFiles, testFile)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return testFiles, nil
}

// GetTestCount returns the total number of tests across all parsed files
func GetTestCount(files []*TestFile) int {
	count := 0
	for _, f := range files {
		count += len(f.Tests)
	}
	return count
}

// GetTestsByCommand groups tests by their command name
func GetTestsByCommand(files []*TestFile) map[string][]TestSpec {
	result := make(map[string][]TestSpec)
	
	for _, file := range files {
		for _, test := range file.Tests {
			parts := strings.Fields(test.Command)
			if len(parts) > 0 {
				cmd := parts[0]
				result[cmd] = append(result[cmd], test)
			}
		}
	}
	
	return result
}