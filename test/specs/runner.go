package specs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TestResult represents the outcome of running a single test
type TestResult struct {
	Spec        TestSpec
	Passed      bool
	ActualOut   []string
	ActualErr   []string
	Diff        string // Difference between expected and actual
	Error       error  // Any execution error
}

// Runner executes test specs and compares results
type Runner struct {
	ledgerPath  string // Path to ledger executable
	gledgerPath string // Path to gledger executable
	tempDir     string // Temporary directory for test files
	verbose     bool
}

// NewRunner creates a new test runner
func NewRunner(ledgerPath, gledgerPath string, verbose bool) (*Runner, error) {
	// Create temporary directory for test files
	tempDir, err := ioutil.TempDir("", "gledger-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	return &Runner{
		ledgerPath:  ledgerPath,
		gledgerPath: gledgerPath,
		tempDir:     tempDir,
		verbose:     verbose,
	}, nil
}

// Cleanup removes temporary files
func (r *Runner) Cleanup() error {
	return os.RemoveAll(r.tempDir)
}

// RunTest executes a single test spec with gledger
func (r *Runner) RunTest(spec TestSpec) TestResult {
	result := TestResult{
		Spec: spec,
	}

	// Create temporary input file
	inputFile, err := r.createInputFile(spec.InputData)
	if err != nil {
		result.Error = fmt.Errorf("failed to create input file: %w", err)
		return result
	}
	defer os.Remove(inputFile)

	// Run gledger with the test command
	output, err := r.runCommand(r.gledgerPath, spec.Command, inputFile)
	if err != nil && !isExpectedError(spec) {
		result.Error = fmt.Errorf("command execution failed: %w", err)
		return result
	}

	// Split output into lines
	result.ActualOut = splitLines(output)

	// Compare output with expected
	result.Passed = compareOutput(result.ActualOut, spec.ExpectedOut)
	if !result.Passed {
		result.Diff = generateDiff(spec.ExpectedOut, result.ActualOut)
	}

	return result
}

// RunComparison runs the same test with both ledger and gledger and compares
func (r *Runner) RunComparison(spec TestSpec) TestResult {
	result := TestResult{
		Spec: spec,
	}

	// Create temporary input file
	inputFile, err := r.createInputFile(spec.InputData)
	if err != nil {
		result.Error = fmt.Errorf("failed to create input file: %w", err)
		return result
	}
	defer os.Remove(inputFile)

	// Run original ledger
	ledgerOutput, ledgerErr := r.runCommand(r.ledgerPath, spec.Command, inputFile)
	
	// Run gledger
	gledgerOutput, gledgerErr := r.runCommand(r.gledgerPath, spec.Command, inputFile)

	// Check if errors match
	if (ledgerErr != nil) != (gledgerErr != nil) {
		result.Error = fmt.Errorf("error mismatch: ledger=%v, gledger=%v", ledgerErr, gledgerErr)
		return result
	}

	// Compare outputs
	ledgerLines := splitLines(ledgerOutput)
	gledgerLines := splitLines(gledgerOutput)

	result.ActualOut = gledgerLines
	result.Passed = compareOutput(ledgerLines, gledgerLines)
	
	if !result.Passed {
		result.Diff = generateDiff(ledgerLines, gledgerLines)
	}

	return result
}

// RunTestFile runs all tests in a test file
func (r *Runner) RunTestFile(file *TestFile) []TestResult {
	results := make([]TestResult, 0, len(file.Tests))
	
	for _, spec := range file.Tests {
		if r.verbose {
			fmt.Printf("Running test: %s\n", spec.Name)
		}
		result := r.RunTest(spec)
		results = append(results, result)
	}
	
	return results
}

// createInputFile creates a temporary file with the test input data
func (r *Runner) createInputFile(lines []string) (string, error) {
	file, err := ioutil.TempFile(r.tempDir, "test-*.ledger")
	if err != nil {
		return "", err
	}
	defer file.Close()

	content := strings.Join(lines, "\n")
	if _, err := file.WriteString(content); err != nil {
		os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

// runCommand executes a ledger/gledger command and returns the output
func (r *Runner) runCommand(executable, command, inputFile string) (string, error) {
	// Parse command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	// Build full command arguments
	args := append([]string{"-f", inputFile}, parts...)

	// Execute command
	cmd := exec.Command(executable, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// For ledger commands, we typically want stdout even if there's an error
	// (e.g., validation errors might be expected)
	output := stdout.String()
	
	if err != nil && stderr.Len() > 0 {
		// Include stderr in error for debugging
		return output, fmt.Errorf("%w: %s", err, stderr.String())
	}

	return output, err
}

// Helper functions

func isExpectedError(spec TestSpec) bool {
	// Check if the test expects an error (has expected error output)
	return len(spec.ExpectedErr) > 0
}

func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return lines
}

// normalizeWhitespace normalizes all whitespace in a string for comparison
// - Trims leading and trailing whitespace
// - Replaces multiple consecutive spaces/tabs with single spaces
// - Normalizes different types of whitespace to regular spaces
func normalizeWhitespace(s string) string {
	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)
	
	// Replace all whitespace sequences (spaces, tabs, etc.) with single spaces
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	
	return s
}

func compareOutput(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	
	for i := range actual {
		// Normalize whitespace for comparison
		actualLine := normalizeWhitespace(actual[i])
		expectedLine := normalizeWhitespace(expected[i])
		
		if actualLine != expectedLine {
			return false
		}
	}
	
	return true
}

func generateDiff(expected, actual []string) string {
	var diff strings.Builder
	
	diff.WriteString("Diff:\n")
	diff.WriteString("  Expected output:\n")
	for i, line := range expected {
		diff.WriteString(fmt.Sprintf("  %3d: %s\n", i+1, line))
	}
	
	diff.WriteString("\nActual output:\n")
	for i, line := range actual {
		diff.WriteString(fmt.Sprintf("  %3d: %s\n", i+1, line))
	}
	
	// Show first differing line with normalized comparison
	minLen := len(expected)
	if len(actual) < minLen {
		minLen = len(actual)
	}
	
	for i := 0; i < minLen; i++ {
		expectedNorm := normalizeWhitespace(expected[i])
		actualNorm := normalizeWhitespace(actual[i])
		if expectedNorm != actualNorm {
			diff.WriteString(fmt.Sprintf("\nFirst difference at line %d:\n", i+1))
			diff.WriteString(fmt.Sprintf("  Expected (normalized): %q\n", expectedNorm))
			diff.WriteString(fmt.Sprintf("  Actual (normalized):   %q\n", actualNorm))
			break
		}
	}
	
	return diff.String()
}

// TestSummary provides a summary of test results
type TestSummary struct {
	Total   int
	Passed  int
	Failed  int
	Errors  int
	Skipped int
}

// Summarize creates a summary of test results
func Summarize(results []TestResult) TestSummary {
	summary := TestSummary{
		Total: len(results),
	}
	
	for _, r := range results {
		if r.Error != nil {
			summary.Errors++
		} else if r.Passed {
			summary.Passed++
		} else {
			summary.Failed++
		}
	}
	
	return summary
}

// PrintSummary prints a test summary to stdout
func PrintSummary(summary TestSummary) {
	fmt.Printf("\nTest Summary:\n")
	fmt.Printf("  Total:   %d\n", summary.Total)
	fmt.Printf("  Passed:  %d\n", summary.Passed)
	fmt.Printf("  Failed:  %d\n", summary.Failed)
	fmt.Printf("  Errors:  %d\n", summary.Errors)
	if summary.Skipped > 0 {
		fmt.Printf("  Skipped: %d\n", summary.Skipped)
	}
	
	passRate := float64(summary.Passed) / float64(summary.Total) * 100
	fmt.Printf("  Pass Rate: %.1f%%\n", passRate)
}