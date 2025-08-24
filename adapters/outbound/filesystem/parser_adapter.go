package filesystem

import (
	"io"

	"github.com/hirosato/gledger/domain"
	"github.com/hirosato/gledger/domain/ports"
	"github.com/hirosato/gledger/infrastructure/parser"
)

// ParserAdapter adapts the infrastructure parser to the domain port
type ParserAdapter struct {
	parser *parser.Parser
}

// NewParserAdapter creates a new parser adapter
func NewParserAdapter() ports.Parser {
	return &ParserAdapter{
		parser: parser.NewParser(),
	}
}

// Parse implements the Parser interface
func (p *ParserAdapter) Parse(reader io.Reader) ([]domain.Transaction, []domain.Directive, error) {
	// Use the existing parser implementation
	if err := p.parser.Parse(reader); err != nil {
		return nil, nil, err
	}
	
	// Return the parsed data
	transactions := p.parser.GetTransactions()
	// TODO: The current parser doesn't support directives yet
	// This will need to be implemented in the parser
	directives := []domain.Directive{}
	
	return transactions, directives, nil
}