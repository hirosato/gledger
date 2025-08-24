package specs

import (
	"strings"
	"testing"
)

func TestParseSimpleTest(t *testing.T) {
	input := `2011-01-01 * Opening balance
    Assets:Bank                    10.00 GBP
    Equity:Opening balance

test accounts
Assets:Bank
Equity:Opening balance
end test

test balance
           10.00 GBP  Assets:Bank
          -10.00 GBP  Equity:Opening balance
--------------------
                   0
end test`

	parser := NewParser(false)
	testFile, err := parser.parse(strings.NewReader(input), "test.ledger")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Check input data was captured (3 transaction lines + 1 blank line)
	if len(testFile.InputData) != 4 {
		t.Errorf("Expected 4 lines of input data, got %d", len(testFile.InputData))
	}

	// Check we found 2 tests
	if len(testFile.Tests) != 2 {
		t.Fatalf("Expected 2 tests, got %d", len(testFile.Tests))
	}

	// Check first test
	test1 := testFile.Tests[0]
	if test1.Name != "accounts" {
		t.Errorf("Expected first test name 'accounts', got '%s'", test1.Name)
	}
	if len(test1.ExpectedOut) != 2 {
		t.Errorf("Expected 2 lines of output for first test, got %d", len(test1.ExpectedOut))
	}

	// Check second test
	test2 := testFile.Tests[1]
	if test2.Name != "balance" {
		t.Errorf("Expected second test name 'balance', got '%s'", test2.Name)
	}
	if len(test2.ExpectedOut) != 4 {
		t.Errorf("Expected 4 lines of output for second test, got %d", len(test2.ExpectedOut))
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		testName string
		expected string
	}{
		{"accounts", "accounts"},
		{"balance --flat", "balance --flat"},
		{"register checking", "register checking"},
		{"print --market", "print --market"},
	}

	parser := NewParser(false)
	for _, tc := range tests {
		result := parser.parseCommand(tc.testName)
		if result != tc.expected {
			t.Errorf("For test name '%s', expected command '%s', got '%s'",
				tc.testName, tc.expected, result)
		}
	}
}

func TestGetTestsByCommand(t *testing.T) {
	// Create test data
	files := []*TestFile{
		{
			Path: "test1.ledger",
			Tests: []TestSpec{
				{Name: "accounts", Command: "accounts"},
				{Name: "balance", Command: "balance"},
				{Name: "balance --flat", Command: "balance --flat"},
			},
		},
		{
			Path: "test2.ledger",
			Tests: []TestSpec{
				{Name: "register", Command: "register"},
				{Name: "balance", Command: "balance"},
			},
		},
	}

	byCommand := GetTestsByCommand(files)

	// Check we have the right commands
	if len(byCommand["accounts"]) != 1 {
		t.Errorf("Expected 1 'accounts' test, got %d", len(byCommand["accounts"]))
	}
	if len(byCommand["balance"]) != 3 {
		t.Errorf("Expected 3 'balance' tests, got %d", len(byCommand["balance"]))
	}
	if len(byCommand["register"]) != 1 {
		t.Errorf("Expected 1 'register' test, got %d", len(byCommand["register"]))
	}
}