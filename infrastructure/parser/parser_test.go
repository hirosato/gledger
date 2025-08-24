package parser

import (
	"strings"
	"testing"
)

func TestParseDate(t *testing.T) {
	p := NewParser()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"2011-01-01", "2011-01-01"},
		{"2011/01/01", "2011-01-01"},
		{"2012-12-31", "2012-12-31"},
	}
	
	for _, test := range tests {
		date, err := p.parseDate(test.input)
		if err != nil {
			t.Errorf("Failed to parse date '%s': %v", test.input, err)
			continue
		}
		
		result := date.Format("2006-01-02")
		if result != test.expected {
			t.Errorf("For date '%s', expected '%s', got '%s'", 
				test.input, test.expected, result)
		}
	}
}

func TestParseAmount(t *testing.T) {
	p := NewParser()
	
	tests := []struct {
		input             string
		expectedValue     float64
		expectedCommodity string
	}{
		{"10.00 GBP", 10.00, "GBP"},
		{"$25.50", 25.50, "$"},
		{"100", 100.00, "$"},
		{"-5.25 EUR", -5.25, "EUR"},
	}
	
	for _, test := range tests {
		amount, err := p.parseAmount(test.input)
		if err != nil {
			t.Errorf("Failed to parse amount '%s': %v", test.input, err)
			continue
		}
		
		if amount.ToFloat64() != test.expectedValue {
			t.Errorf("For amount '%s', expected value %f, got %f", 
				test.input, test.expectedValue, amount.ToFloat64())
		}
		
		if amount.Commodity.Symbol != test.expectedCommodity {
			t.Errorf("For amount '%s', expected commodity '%s', got '%s'", 
				test.input, test.expectedCommodity, amount.Commodity.Symbol)
		}
	}
}

func TestParseSimpleTransaction(t *testing.T) {
	p := NewParser()
	
	input := `2011-01-01 * Opening balance
    Assets:Cash                    10.00 USD
    Equity:Opening balance`
	
	err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Failed to parse transaction: %v", err)
	}
	
	transactions := p.GetTransactions()
	if len(transactions) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(transactions))
	}
	
	tx := transactions[0]
	if tx.Payee != "Opening balance" {
		t.Errorf("Expected payee 'Opening balance', got '%s'", tx.Payee)
	}
	
	if len(tx.Postings) != 2 {
		t.Fatalf("Expected 2 postings, got %d", len(tx.Postings))
	}
	
	// Check first posting
	posting1 := tx.Postings[0]
	if posting1.Account.FullName != "Assets:Cash" {
		t.Errorf("Expected first posting account 'Assets:Cash', got '%s'", 
			posting1.Account.FullName)
	}
	
	if posting1.Amount == nil {
		t.Error("Expected first posting to have amount")
	} else if posting1.Amount.ToFloat64() != 10.00 {
		t.Errorf("Expected first posting amount 10.00, got %f", 
			posting1.Amount.ToFloat64())
	}
	
	// Check second posting (should have elided amount)
	posting2 := tx.Postings[1]
	if posting2.Account.FullName != "Equity:Opening balance" {
		t.Errorf("Expected second posting account 'Equity:Opening balance', got '%s'", 
			posting2.Account.FullName)
	}
	
	if posting2.Amount == nil {
		t.Error("Expected second posting to have elided amount")
	} else if posting2.Amount.ToFloat64() != -10.00 {
		t.Errorf("Expected second posting amount -10.00, got %f", 
			posting2.Amount.ToFloat64())
	}
}