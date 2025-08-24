package domain

import (
	"testing"
)

func TestNewAccount(t *testing.T) {
	account := NewAccount("Assets:Cash")
	
	if account.Name != "Assets:Cash" {
		t.Errorf("Expected account name 'Assets:Cash', got '%s'", account.Name)
	}
	
	if account.FullName != "Assets:Cash" {
		t.Errorf("Expected full name 'Assets:Cash', got '%s'", account.FullName)
	}
}

func TestAccountHierarchy(t *testing.T) {
	parent := NewAccount("Assets")
	child := NewAccount("Cash")
	
	parent.AddChild(child)
	
	if child.Parent != parent {
		t.Error("Child should reference parent")
	}
	
	if child.FullName != "Assets:Cash" {
		t.Errorf("Expected child full name 'Assets:Cash', got '%s'", child.FullName)
	}
	
	if len(parent.Children) != 1 {
		t.Errorf("Expected parent to have 1 child, got %d", len(parent.Children))
	}
}

func TestDetermineAccountType(t *testing.T) {
	tests := []struct {
		name     string
		expected AccountType
	}{
		{"Assets:Cash", AccountTypeAsset},
		{"Liabilities:Loan", AccountTypeLiability},
		{"Equity:Opening", AccountTypeEquity},
		{"Income:Salary", AccountTypeIncome},
		{"Expenses:Food", AccountTypeExpense},
	}
	
	for _, test := range tests {
		result := DetermineAccountType(test.name)
		if result != test.expected {
			t.Errorf("For account '%s', expected type %v, got %v", 
				test.name, test.expected, result)
		}
	}
}