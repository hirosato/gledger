package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirosato/gledger/test/specs"
)

func main() {
	var (
		testDir     = flag.String("dir", "", "Directory containing test files")
		testFile    = flag.String("file", "", "Single test file to run")
		ledgerPath  = flag.String("ledger", "ledger", "Path to ledger executable")
		gledgerPath = flag.String("gledger", "./build/gledger", "Path to gledger executable")
		verbose     = flag.Bool("verbose", false, "Verbose output")
		compare     = flag.Bool("compare", false, "Compare with ledger output")
		listOnly    = flag.Bool("list", false, "List tests without running")
	)

	flag.Parse()

	if *testDir == "" && *testFile == "" {
		// Default to baseline tests
		*testDir = "../../ledger/test/baseline"
	}

	// Create parser
	parser := specs.NewParser(*verbose)

	// Parse test files
	var testFiles []*specs.TestFile
	var err error

	if *testFile != "" {
		// Parse single file
		file, err := parser.ParseFile(*testFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing test file: %v\n", err)
			os.Exit(1)
		}
		testFiles = []*specs.TestFile{file}
	} else if *testDir != "" {
		// Parse directory
		testFiles, err = parser.ParseDirectory(*testDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing test directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Report parsed tests
	totalTests := specs.GetTestCount(testFiles)
	fmt.Printf("Parsed %d test files containing %d tests\n", len(testFiles), totalTests)

	if *listOnly {
		// Just list the tests
		listTests(testFiles)
		return
	}

	// Create runner
	runner, err := specs.NewRunner(*ledgerPath, *gledgerPath, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating test runner: %v\n", err)
		os.Exit(1)
	}
	defer runner.Cleanup()

	// Run tests
	var allResults []specs.TestResult

	for _, file := range testFiles {
		fmt.Printf("\nRunning tests from: %s\n", filepath.Base(file.Path))
		
		for _, test := range file.Tests {
			var result specs.TestResult
			
			if *compare {
				// Compare with ledger
				result = runner.RunComparison(test)
			} else {
				// Run against expected output
				result = runner.RunTest(test)
			}
			
			allResults = append(allResults, result)
			
			// Print result
			printResult(result, *verbose)
		}
	}

	// Print summary
	summary := specs.Summarize(allResults)
	specs.PrintSummary(summary)

	// Exit with error if tests failed
	if summary.Failed > 0 || summary.Errors > 0 {
		os.Exit(1)
	}
}

func listTests(files []*specs.TestFile) {
	byCommand := specs.GetTestsByCommand(files)
	
	fmt.Println("\nTests by command:")
	for cmd, tests := range byCommand {
		fmt.Printf("\n%s (%d tests):\n", cmd, len(tests))
		for _, test := range tests {
			fmt.Printf("  - %s (from %s:%d)\n", 
				test.Name, 
				filepath.Base(test.SourceFile), 
				test.LineNumber)
		}
	}
}

func printResult(result specs.TestResult, verbose bool) {
	status := "✓"
	if !result.Passed {
		status = "✗"
	}
	if result.Error != nil {
		status = "E"
	}

	fmt.Printf("  %s %s\n", status, result.Spec.Name)
	
	if !result.Passed && verbose {
		if result.Error != nil {
			fmt.Printf("    Error: %v\n", result.Error)
		} else if result.Diff != "" {
			fmt.Printf("    Diff:\n%s\n", indent(result.Diff, "      "))
		}
	}
}

func indent(s string, prefix string) string {
	// Simple indentation helper
	return prefix + s // Simplified for now
}