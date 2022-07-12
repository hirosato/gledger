package main

import (
	"bytes"
	"testing"
)

func TestLedger(t *testing.T) {
	testLedger := bytes.NewBufferString(`
1970-01-01 Payee
  Expense:Food 100
  Assets:Cash -100

1970-01-02 Payee
  Expense:Food 200
  Assets:Cash -200
`)
	ledger, _ := ParseLedger(testLedger)
	if len(ledger) != 2 {
		t.Error("There should be 2 ledgers")
	}
}
