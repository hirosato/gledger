package ledger_test

import (
	"bytes"
	"testing"

	ledger "github.com/hirosato/ledger/lib"
)

func TestBalance(t *testing.T) {
	testLedger := bytes.NewBufferString(`
; Comment
1970/01/01 Payee
  Expense:Food  100
  Assets:Cash  -100
1970/01/01 Payee
  Expense:その他  100
  Assets:Cash  -100
1970/01/01 Payee
  Assets:Fixed Asset  100
  Assets:Cash  -100
`)
	ledgers, _ := ledger.ParseLedger(testLedger)
	balances, _ := ledger.NewBalances(ledgers)
	balances.PrettyPrint()
}
