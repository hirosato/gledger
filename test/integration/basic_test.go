package integration

import (
	"testing"
	
	"github.com/hirosato/gledger/test/testutil"
)

func TestBasicTransactions(t *testing.T) {
	// Create a simple test spec
	spec := &testutil.TestSpec{
		Name: "basic balance",
		InputData: `2024-01-01 Opening Balance
    Assets:Bank                 $1000.00
    Equity:Opening

2024-01-02 Grocery Store
    Expenses:Food                $50.00
    Assets:Bank

2024-01-03 Paycheck
    Assets:Bank                $2000.00
    Income:Salary
`,
		Command:        "balance",
		ExpectedOutput: `              $2950.00  Assets:Bank
             $-1000.00  Equity:Opening
                $50.00  Expenses:Food
             $-2000.00  Income:Salary
--------------------
                     0`,
	}

	// TODO: Implement actual test when domain models are ready
	t.Skip("Skipping until domain models are implemented")
	
	runner, err := testutil.NewTestRunner()
	if err != nil {
		t.Fatalf("Failed to create test runner: %v", err)
	}
	defer runner.Cleanup()

	runner.RunSpec(t, spec)
}