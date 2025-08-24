package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hirosato/gledger/domain"
)

// Parser parses ledger journal files
type Parser struct {
	scanner      *bufio.Scanner
	currentLine  string
	lineNumber   int
	transactions []domain.Transaction
	accounts     map[string]bool
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		accounts: make(map[string]bool),
	}
}

// Parse parses a ledger journal from the given reader
func (p *Parser) Parse(reader io.Reader) error {
	p.scanner = bufio.NewScanner(reader)
	p.lineNumber = 0
	p.transactions = []domain.Transaction{}

	for p.advance() {
		// Skip empty lines
		if strings.TrimSpace(p.currentLine) == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(strings.TrimSpace(p.currentLine), ";") {
			continue
		}

		// Check if this is a transaction line
		if p.isTransactionLine() {
			transaction, err := p.parseTransaction()
			if err != nil {
				return fmt.Errorf("line %d: %w", p.lineNumber, err)
			}
			p.transactions = append(p.transactions, *transaction)
		}
	}

	return nil
}

// advance reads the next line
func (p *Parser) advance() bool {
	if p.scanner.Scan() {
		p.currentLine = p.scanner.Text()
		p.lineNumber++
		return true
	}
	return false
}

// isTransactionLine checks if the current line starts a transaction
func (p *Parser) isTransactionLine() bool {
	line := p.currentLine
	
	// Transaction lines start with a date (YYYY-MM-DD or YYYY/MM/DD)
	if len(line) < 10 {
		return false
	}

	// Check for date pattern
	datePart := line[:10]
	if (datePart[4] == '-' || datePart[4] == '/') && 
	   (datePart[7] == '-' || datePart[7] == '/') {
		// Try to parse the date
		_, err := p.parseDate(datePart)
		return err == nil
	}

	return false
}

// parseTransaction parses a complete transaction
func (p *Parser) parseTransaction() (*domain.Transaction, error) {
	// Parse the transaction header line
	date, status, description, err := p.parseTransactionHeader()
	if err != nil {
		return nil, err
	}

	transaction := &domain.Transaction{
		Date:     date,
		Status:   status,
		Payee:    description,
		Postings: []*domain.Posting{},
	}

	// Parse postings (indented lines following the transaction)
	for p.advance() {
		// If line is not indented, we've reached the end of this transaction
		if !strings.HasPrefix(p.currentLine, " ") && !strings.HasPrefix(p.currentLine, "\t") {
			// Back up one line for the next iteration
			p.lineNumber--
			break
		}

		// Skip empty lines and comments within transaction
		trimmed := strings.TrimSpace(p.currentLine)
		if trimmed == "" || strings.HasPrefix(trimmed, ";") {
			continue
		}

		// Parse the posting
		posting, err := p.parsePosting()
		if err != nil {
			return nil, fmt.Errorf("posting error: %w", err)
		}
		
		transaction.Postings = append(transaction.Postings, posting)
	}

	// Validate transaction has at least 2 postings
	if len(transaction.Postings) < 2 {
		return nil, fmt.Errorf("transaction must have at least 2 postings")
	}

	// Apply amount elision if needed
	if err := p.applyAmountElision(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// parseTransactionHeader parses the first line of a transaction
func (p *Parser) parseTransactionHeader() (time.Time, domain.TransactionStatus, string, error) {
	line := p.currentLine
	
	// Extract date (first 10 characters)
	dateStr := line[:10]
	date, err := p.parseDate(dateStr)
	if err != nil {
		return time.Time{}, 0, "", err
	}

	// Parse rest of line
	rest := strings.TrimSpace(line[10:])
	
	// Check for status marker
	status := domain.TransactionStatusPending
	if strings.HasPrefix(rest, "*") {
		status = domain.TransactionStatusCleared
		rest = strings.TrimSpace(rest[1:])
	} else if strings.HasPrefix(rest, "!") {
		status = domain.TransactionStatusPending
		rest = strings.TrimSpace(rest[1:])
	}

	// The rest is the description
	description := rest

	return date, status, description, nil
}

// parseDate parses a date string in YYYY-MM-DD or YYYY/MM/DD format
func (p *Parser) parseDate(dateStr string) (time.Time, error) {
	// Normalize separators
	normalized := strings.ReplaceAll(dateStr, "/", "-")
	
	// Parse as YYYY-MM-DD
	date, err := time.Parse("2006-01-02", normalized)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}
	
	return date, nil
}

// parsePosting parses a posting line
func (p *Parser) parsePosting() (*domain.Posting, error) {
	line := strings.TrimSpace(p.currentLine)
	
	// Find the account name (everything before double space that indicates amount)
	// If there's no double space, the whole line is the account name
	accountName := line
	amountStr := ""
	
	// Look for double space or tab separator
	idx := -1
	for i := 0; i < len(line)-1; i++ {
		if (line[i] == ' ' && line[i+1] == ' ') || line[i] == '\t' {
			idx = i
			break
		}
	}
	
	if idx > 0 {
		accountName = strings.TrimSpace(line[:idx])
		amountStr = strings.TrimSpace(line[idx:])
	} else {
		accountName = strings.TrimSpace(line)
		amountStr = ""
	}
	
	// Register the account
	p.accounts[accountName] = true
	
	// Create account object
	account := domain.NewAccount(accountName)
	account.FullName = accountName
	
	// Parse amount if present
	var amount *domain.Amount
	if amountStr != "" {
		parsedAmount, err := p.parseAmount(amountStr)
		if err == nil {
			amount = parsedAmount
		}
	}

	posting := domain.NewPosting(account)
	posting.Amount = amount

	return posting, nil
}

// parseAmount parses an amount string like "10.00 GBP" or "$25.50"
func (p *Parser) parseAmount(amountStr string) (*domain.Amount, error) {
	amountStr = strings.TrimSpace(amountStr)
	
	// Handle currency symbols at the beginning
	commoditySymbol := ""
	valueStr := amountStr
	
	if strings.HasPrefix(amountStr, "$") {
		commoditySymbol = "$"
		valueStr = strings.TrimSpace(amountStr[1:])
	} else if strings.HasPrefix(amountStr, "£") {
		commoditySymbol = "£"
		valueStr = strings.TrimSpace(amountStr[1:])
	} else if strings.HasPrefix(amountStr, "€") {
		commoditySymbol = "€"
		valueStr = strings.TrimSpace(amountStr[1:])
	} else {
		// Look for commodity at the end
		parts := strings.Fields(amountStr)
		if len(parts) == 2 {
			valueStr = parts[0]
			commoditySymbol = parts[1]
		} else if len(parts) == 1 {
			// Just a number, no commodity
			valueStr = parts[0]
			commoditySymbol = "$" // Default commodity
		}
	}

	// Parse the numeric value
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %s", amountStr)
	}

	// Create commodity
	commodity := domain.NewCommodity(commoditySymbol)
	
	// Create amount
	return domain.NewAmountFromFloat(value, commodity), nil
}

// applyAmountElision fills in missing amounts in postings
func (p *Parser) applyAmountElision(transaction *domain.Transaction) error {
	// Count postings without amounts
	var missingIndex = -1
	missingCount := 0
	
	for i, posting := range transaction.Postings {
		if posting.Amount == nil {
			missingIndex = i
			missingCount++
		}
	}

	// Can only elide one amount
	if missingCount > 1 {
		return fmt.Errorf("only one posting can have an elided amount")
	}

	// If one amount is missing, calculate it to balance the transaction
	if missingCount == 1 {
		// Sum all known amounts by commodity
		sums := make(map[string]float64)
		commodities := make(map[string]*domain.Commodity)
		
		for i, posting := range transaction.Postings {
			if i != missingIndex && posting.Amount != nil {
				commoditySymbol := posting.Amount.Commodity.Symbol
				if commoditySymbol == "" {
					commoditySymbol = "$" // Default commodity
				}
				sums[commoditySymbol] += posting.Amount.ToFloat64()
				commodities[commoditySymbol] = posting.Amount.Commodity
			}
		}

		// The elided amount should be the negative of the sum
		// For simplicity, assume single commodity for now
		if len(sums) == 1 {
			for commoditySymbol, sum := range sums {
				commodity := commodities[commoditySymbol]
				if commodity == nil {
					commodity = domain.NewCommodity(commoditySymbol)
				}
				transaction.Postings[missingIndex].Amount = domain.NewAmountFromFloat(-sum, commodity)
			}
		} else if len(sums) > 1 {
			return fmt.Errorf("cannot elide amount with multiple commodities")
		}
	}

	return nil
}

// GetTransactions returns all parsed transactions
func (p *Parser) GetTransactions() []domain.Transaction {
	return p.transactions
}

// GetAccounts returns all account names found during parsing
func (p *Parser) GetAccounts() []string {
	accounts := make([]string, 0, len(p.accounts))
	for account := range p.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}