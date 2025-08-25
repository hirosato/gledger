package testutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunner runs test specs and compares with ledger-cli
type TestRunner struct {
	GledgerPath string
	LedgerPath  string
	TempDir     string
}

// NewTestRunner creates a new test runner
func NewTestRunner() (*TestRunner, error) {
	// Find gledger binary - check multiple possible locations
	possiblePaths := []string{
		filepath.Join("build", "gledger"),
		filepath.Join("..", "..", "build", "gledger"),
		filepath.Join(".", "build", "gledger"),
	}
	
	var gledgerPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			gledgerPath = path
			break
		}
	}
	
	if gledgerPath == "" {
		// Try to build it from the root directory
		cmd := exec.Command("make", "build")
		cmd.Dir = filepath.Join("..", "..")
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to build gledger: %w", err)
		}
		// Check again
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				gledgerPath = path
				break
			}
		}
		if gledgerPath == "" {
			return nil, fmt.Errorf("gledger binary not found after build")
		}
	}

	// Find ledger binary
	ledgerPath, err := exec.LookPath("ledger")
	if err != nil {
		return nil, fmt.Errorf("ledger-cli not found in PATH: %w", err)
	}

	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "gledger-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	return &TestRunner{
		GledgerPath: gledgerPath,
		LedgerPath:  ledgerPath,
		TempDir:     tempDir,
	}, nil
}

// Cleanup removes temporary files
func (r *TestRunner) Cleanup() {
	if r.TempDir != "" {
		os.RemoveAll(r.TempDir)
	}
}

// RunSpec runs a single test spec
func (r *TestRunner) RunSpec(t *testing.T, spec *TestSpec) {
	t.Helper()

	// Create temp file with input data
	inputFile := filepath.Join(r.TempDir, "input.ledger")
	if err := os.WriteFile(inputFile, []byte(spec.InputData), 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	// Parse command and build args
	args := r.parseCommand(spec.Command, inputFile)

	// Run gledger
	gledgerOutput, gledgerErr := r.runCommand(r.GledgerPath, args)

	// Check error expectation
	if spec.ShouldFail {
		if gledgerErr == nil {
			t.Errorf("expected command to fail but it succeeded")
		}
		return
	}

	if gledgerErr != nil {
		t.Errorf("gledger failed: %v\nOutput: %s", gledgerErr, gledgerOutput)
		return
	}

	// Compare output
	if !CompareOutput(spec.ExpectedOutput, gledgerOutput) {
		t.Errorf("output mismatch\nExpected:\n%s\n\nActual:\n%s",
			spec.ExpectedOutput, gledgerOutput)
	}
}

// RunComparison runs both gledger and ledger-cli and compares outputs
func (r *TestRunner) RunComparison(t *testing.T, spec *TestSpec) {
	t.Helper()

	// Create temp file with input data
	inputFile := filepath.Join(r.TempDir, "input.ledger")
	if err := os.WriteFile(inputFile, []byte(spec.InputData), 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	// Parse command and build args
	args := r.parseCommand(spec.Command, inputFile)

	// Run both commands
	gledgerOutput, gledgerErr := r.runCommand(r.GledgerPath, args)
	ledgerOutput, ledgerErr := r.runCommand(r.LedgerPath, args)

	// Compare error states
	if (gledgerErr != nil) != (ledgerErr != nil) {
		t.Errorf("error state mismatch\ngledger error: %v\nledger error: %v",
			gledgerErr, ledgerErr)
		return
	}

	// If both failed, that's OK
	if gledgerErr != nil && ledgerErr != nil {
		return
	}

	// Compare outputs
	if !CompareOutput(ledgerOutput, gledgerOutput) {
		t.Errorf("output mismatch\nLedger output:\n%s\n\nGledger output:\n%s",
			ledgerOutput, gledgerOutput)
	}
}

// parseCommand parses a test command into arguments
func (r *TestRunner) parseCommand(cmd string, inputFile string) []string {
	// Split command into parts
	parts := strings.Fields(cmd)
	
	// Add -f flag for input file
	args := []string{"-f", inputFile}
	
	// Add command parts
	args = append(args, parts...)
	
	return args
}

// runCommand executes a command and returns output
func (r *TestRunner) runCommand(cmdPath string, args []string) (string, error) {
	cmd := exec.Command(cmdPath, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	// Combine stdout and stderr for error cases
	output := stdout.String()
	if err != nil && stderr.Len() > 0 {
		output = stderr.String()
	}
	
	return output, err
}