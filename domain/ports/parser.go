package ports

import (
	"io"

	"github.com/hirosato/gledger/domain"
)

// Parser defines the interface for parsing ledger journal files
type Parser interface {
	// Parse reads from the provided reader and returns transactions and directives
	Parse(reader io.Reader) ([]domain.Transaction, []domain.Directive, error)
}