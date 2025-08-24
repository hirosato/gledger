package application

import (
	"strings"
	"testing"
)

func TestJournalBasicFunctionality(t *testing.T) {
	journal := NewJournal()
	
	if journal == nil {
		t.Fatal("Failed to create new journal")
	}
	
	// Test that journal starts empty
	accounts := journal.GetAccounts()
	if len(accounts) != 0 {
		t.Errorf("Expected empty journal to have 0 accounts, got %d", len(accounts))
	}
	
	transactions := journal.GetTransactions()
	if len(transactions) != 0 {
		t.Errorf("Expected empty journal to have 0 transactions, got %d", len(transactions))
	}
}

func TestJournalLoadFromReader(t *testing.T) {
	journal := NewJournal()
	
	input := `2011-01-01 * Opening balance
    Assets:Cash                    10.00 USD
    Equity:Opening balance

2011-01-02 * Purchase
    Expenses:Food                   5.00 USD
    Assets:Cash`
	
	err := journal.LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Failed to load journal: %v", err)
	}
	
	// Check transactions loaded
	transactions := journal.GetTransactions()
	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}
	
	// Check accounts detected
	accounts := journal.GetAccounts()
	expectedAccounts := []string{
		"Assets:Cash",
		"Equity:Opening balance", 
		"Expenses:Food",
	}
	
	if len(accounts) != len(expectedAccounts) {
		t.Errorf("Expected %d accounts, got %d: %v", 
			len(expectedAccounts), len(accounts), accounts)
	}
	
	// Check account names
	for _, expected := range expectedAccounts {
		found := false
		for _, actual := range accounts {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected account '%s' not found in: %v", expected, accounts)
		}
	}
}

func TestJournalAccountsMatching(t *testing.T) {
	journal := NewJournal()
	
	input := `2011-01-01 * Test
    Assets:Cash                    10.00 USD
    Assets:Bank:Checking            5.00 USD
    Assets:Bank:Savings             2.00 USD
    Equity:Opening balance`
	
	err := journal.LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Failed to load journal: %v", err)
	}
	
	// Test pattern matching (should find accounts containing "bank")
	matches := journal.GetAccountsMatching("bank")
	
	// Should match: Assets:Bank:Checking, Assets:Bank:Savings
	// Note: Assets:Bank is not in the transactions, only in the hierarchy
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'bank', got %d: %v", len(matches), matches)
	}
	
	// Check that we get the right accounts (the ones actually used in postings)
	expectedInMatches := []string{
		"Assets:Bank:Checking",
		"Assets:Bank:Savings",
	}
	
	for _, expected := range expectedInMatches {
		found := false
		for _, match := range matches {
			if match == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected account '%s' not found in matches: %v", expected, matches)
		}
	}
}