package ports

import (
	"io"

	"github.com/hirosato/gledger/domain"
)

// Storage defines the interface for persisting and retrieving journal data
type Storage interface {
	// Load reads journal data from the storage
	Load() ([]domain.Transaction, []domain.Directive, error)
	
	// Save writes journal data to the storage
	Save(transactions []domain.Transaction, directives []domain.Directive) error
	
	// LoadFromReader reads journal data from a reader
	LoadFromReader(reader io.Reader) ([]domain.Transaction, []domain.Directive, error)
	
	// SaveToWriter writes journal data to a writer
	SaveToWriter(writer io.Writer, transactions []domain.Transaction, directives []domain.Directive) error
}