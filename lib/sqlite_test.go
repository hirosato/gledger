package ledger_test

import (
	"testing"

	ledger "github.com/hirosato/ledger/lib"
)

func TestSqlite(t *testing.T) {
	ledger.NewDB()
}
