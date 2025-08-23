package specs

import (
	"path/filepath"
	"testing"
	
	"github.com/hirosato/gledger/test/testutil"
)

// TestBaselineSpecs runs all baseline test specifications
func TestBaselineSpecs(t *testing.T) {
	t.Skip("Skipping until domain models are implemented")
	
	testFiles, err := filepath.Glob("../fixtures/baseline/*.test")
	if err != nil {
		t.Fatalf("Failed to find baseline test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("No baseline test files found. Run 'make import-tests' first.")
	}

	runner, err := testutil.NewTestRunner()
	if err != nil {
		t.Fatalf("Failed to create test runner: %v", err)
	}
	defer runner.Cleanup()

	for _, testFile := range testFiles {
		testName := filepath.Base(testFile)
		t.Run(testName, func(t *testing.T) {
			specs, err := testutil.ParseFile(testFile)
			if err != nil {
				t.Fatalf("Failed to parse test file: %v", err)
			}

			for _, spec := range specs {
				t.Run(spec.Name, func(t *testing.T) {
					runner.RunSpec(t, spec)
				})
			}
		})
	}
}

// TestRegressionSpecs runs all regression test specifications
func TestRegressionSpecs(t *testing.T) {
	t.Skip("Skipping until domain models are implemented")
	
	testFiles, err := filepath.Glob("../fixtures/regress/*.test")
	if err != nil {
		t.Fatalf("Failed to find regression test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("No regression test files found. Run 'make import-tests' first.")
	}

	runner, err := testutil.NewTestRunner()
	if err != nil {
		t.Fatalf("Failed to create test runner: %v", err)
	}
	defer runner.Cleanup()

	for _, testFile := range testFiles {
		testName := filepath.Base(testFile)
		t.Run(testName, func(t *testing.T) {
			specs, err := testutil.ParseFile(testFile)
			if err != nil {
				t.Fatalf("Failed to parse test file: %v", err)
			}

			for _, spec := range specs {
				t.Run(spec.Name, func(t *testing.T) {
					runner.RunSpec(t, spec)
				})
			}
		})
	}
}

// TestComparisonWithLedger runs comparison tests between gledger and ledger-cli
func TestComparisonWithLedger(t *testing.T) {
	t.Skip("Skipping until domain models are implemented")
	
	testFiles, err := filepath.Glob("../fixtures/baseline/cmd-*.test")
	if err != nil {
		t.Fatalf("Failed to find command test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Skip("No command test files found. Run 'make import-tests' first.")
	}

	runner, err := testutil.NewTestRunner()
	if err != nil {
		t.Skipf("Cannot run comparison tests: %v", err)
	}
	defer runner.Cleanup()

	for _, testFile := range testFiles {
		testName := filepath.Base(testFile)
		t.Run(testName, func(t *testing.T) {
			specs, err := testutil.ParseFile(testFile)
			if err != nil {
				t.Fatalf("Failed to parse test file: %v", err)
			}

			for _, spec := range specs {
				t.Run(spec.Name, func(t *testing.T) {
					runner.RunComparison(t, spec)
				})
			}
		})
	}
}