package ledger_test

import (
	"bytes"
	"testing"

	ledger "github.com/hirosato/ledger/lib"
)

func TestLedger(t *testing.T) {
	testLedger := bytes.NewBufferString(`
; Comment
1970-01-01 Payee
  Expense:Food  100
  Assets:Cash  -100
1970-01-02 Payee
  Expense:Food  200;item comment
  資産:流動資産:現金  -200
2021-12-01 Checking balance
  Assets:Bank:普通預金  1000.00
  Equity:Account with Spaces  -1000.00
`)
	ledger, _ := ledger.ParseLedger(testLedger)
	if len(ledger) != 3 {
		t.Error("There should be 3 ledgers")
	}
}
