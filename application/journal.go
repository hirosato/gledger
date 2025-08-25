package application

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hirosato/gledger/domain"
	"github.com/hirosato/gledger/domain/ports"
)

// Journal represents a complete ledger journal with all transactions and accounts
type Journal struct {
	transactions      []domain.Transaction
	accounts          map[string]*domain.Account
	accountTree       *AccountTree
	directives        []domain.Directive
	commodityRegistry map[string]*domain.Commodity
	defaultCommodity  *domain.Commodity
	parser            ports.Parser
}

// NewJournal creates a new empty journal with injected dependencies
func NewJournal(parser ports.Parser) *Journal {
	return &Journal{
		transactions:      []domain.Transaction{},
		accounts:          make(map[string]*domain.Account),
		accountTree:       NewAccountTree(),
		directives:        []domain.Directive{},
		commodityRegistry: make(map[string]*domain.Commodity),
		parser:            parser,
	}
}

// LoadFromReader loads journal data from an io.Reader
func (j *Journal) LoadFromReader(reader io.Reader) error {
	if j.parser == nil {
		return fmt.Errorf("parser not initialized")
	}
	
	transactions, directives, err := j.parser.Parse(reader)
	if err != nil {
		return fmt.Errorf("failed to parse journal: %w", err)
	}

	// Store parsed data
	j.transactions = transactions
	j.directives = directives

	// Build account tree and commodity registry from transactions
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Account != nil {
				j.registerAccount(posting.Account.FullName)
			}
			// Register commodity from amount
			if posting.Amount != nil && posting.Amount.Commodity != nil {
				j.RegisterCommodity(posting.Amount.Commodity)
			}
		}
	}

	return nil
}

// registerAccount registers an account and all its parent accounts
func (j *Journal) registerAccount(fullName string) {
	// Register this account and all parent accounts
	parts := strings.Split(fullName, ":")
	for i := 1; i <= len(parts); i++ {
		accountName := strings.Join(parts[:i], ":")
		if _, exists := j.accounts[accountName]; !exists {
			account := &domain.Account{
				Name:     parts[i-1],
				FullName: accountName,
			}
			j.accounts[accountName] = account
			j.accountTree.AddAccount(account)
		}
	}
}

// GetTransactions returns all transactions in the journal
func (j *Journal) GetTransactions() []domain.Transaction {
	return j.transactions
}

// GetAccounts returns all account names that are used in transactions, sorted alphabetically
func (j *Journal) GetAccounts() []string {
	// Only return accounts that are actually used in postings
	usedAccounts := make(map[string]bool)
	
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Account != nil {
				usedAccounts[posting.Account.FullName] = true
			}
		}
	}
	
	accounts := make([]string, 0, len(usedAccounts))
	for name := range usedAccounts {
		accounts = append(accounts, name)
	}
	sort.Strings(accounts)
	return accounts
}

// GetAccountsMatching returns accounts matching the given pattern
func (j *Journal) GetAccountsMatching(pattern string) []string {
	var matches []string
	
	pattern = strings.ToLower(pattern)
	
	// Only match accounts that are actually used in transactions
	usedAccounts := make(map[string]bool)
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Account != nil {
				usedAccounts[posting.Account.FullName] = true
			}
		}
	}
	
	for name := range usedAccounts {
		nameLower := strings.ToLower(name)
		
		// Simple substring matching for now
		if strings.Contains(nameLower, pattern) {
			matches = append(matches, name)
		}
	}
	
	sort.Strings(matches)
	return matches
}

// GetAccount returns a specific account by name
func (j *Journal) GetAccount(name string) (*domain.Account, bool) {
	account, exists := j.accounts[name]
	return account, exists
}

// GetBalance calculates the balance for an account (includes sub-accounts)
func (j *Journal) GetBalance(accountName string) *domain.Balance {
	balance := domain.NewBalance()
	
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Account != nil {
				accountFullName := posting.Account.FullName
				if accountFullName == accountName || strings.HasPrefix(accountFullName, accountName+":") {
					if posting.Amount != nil {
						balance.Add(posting.Amount)
					}
				}
			}
		}
	}
	
	return balance
}

// GetLeafBalance calculates the balance for an account (excludes sub-accounts)
func (j *Journal) GetLeafBalance(accountName string) *domain.Balance {
	balance := domain.NewBalance()
	
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Account != nil {
				accountFullName := posting.Account.FullName
				if accountFullName == accountName {
					if posting.Amount != nil {
						balance.Add(posting.Amount)
					}
				}
			}
		}
	}
	
	return balance
}

// GetTotalBalance calculates the total balance across all accounts
func (j *Journal) GetTotalBalance() *domain.Balance {
	balance := domain.NewBalance()
	
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			if posting.Amount != nil {
				balance.Add(posting.Amount)
			}
		}
	}
	
	return balance
}

// AddDirective adds a directive to the journal
func (j *Journal) AddDirective(directive domain.Directive) {
	j.directives = append(j.directives, directive)
	
	// Process directive effects
	switch d := directive.(type) {
	case *domain.CommodityDirective:
		commodity := domain.NewCommodity(d.Symbol)
		commodity.Precision = d.Precision
		commodity.Format = d.Format
		j.commodityRegistry[d.Symbol] = commodity
	
	case *domain.PriceDirective:
		if commodity, exists := j.commodityRegistry[d.Commodity]; exists {
			commodity.AddPrice(d.Date, d.Price)
		}
	}
}

// GetDirectives returns all directives
func (j *Journal) GetDirectives() []domain.Directive {
	return j.directives
}

// GetCommodity gets a commodity from the registry
func (j *Journal) GetCommodity(symbol string) (*domain.Commodity, bool) {
	commodity, exists := j.commodityRegistry[symbol]
	return commodity, exists
}

// RegisterCommodity registers a commodity in the registry
func (j *Journal) RegisterCommodity(commodity *domain.Commodity) {
	j.commodityRegistry[commodity.Symbol] = commodity
}

// GetCommodities returns all registered commodities
func (j *Journal) GetCommodities() []string {
	commodities := make([]string, 0, len(j.commodityRegistry))
	for symbol := range j.commodityRegistry {
		commodities = append(commodities, symbol)
	}
	sort.Strings(commodities)
	return commodities
}

// GetCommoditiesForAccount returns commodities used in transactions affecting the given account
func (j *Journal) GetCommoditiesForAccount(accountPattern string) []string {
	commoditySet := make(map[string]bool)
	
	for _, tx := range j.transactions {
		for _, posting := range tx.Postings {
			// Check if this posting affects an account matching the pattern
			if posting.Account != nil && j.matchesAccountPattern(posting.Account.FullName, accountPattern) {
				// Add the commodity from this posting
				if posting.Amount != nil && posting.Amount.Commodity != nil {
					commoditySet[posting.Amount.Commodity.Symbol] = true
				}
			}
		}
	}
	
	var commodities []string
	for symbol := range commoditySet {
		commodities = append(commodities, symbol)
	}
	sort.Strings(commodities)
	return commodities
}

// matchesAccountPattern checks if an account name matches the pattern (simple substring match for now)
func (j *Journal) matchesAccountPattern(accountName, pattern string) bool {
	return strings.Contains(strings.ToLower(accountName), strings.ToLower(pattern))
}

// SetDefaultCommodity sets the default commodity
func (j *Journal) SetDefaultCommodity(commodity *domain.Commodity) {
	j.defaultCommodity = commodity
}

// GetDefaultCommodity returns the default commodity
func (j *Journal) GetDefaultCommodity() *domain.Commodity {
	return j.defaultCommodity
}

// GetPayees returns all unique payees from transactions
func (j *Journal) GetPayees() []string {
	payeeSet := make(map[string]bool)
	
	for _, tx := range j.transactions {
		if tx.Payee != "" {
			payeeSet[tx.Payee] = true
		}
	}
	
	var payees []string
	for payee := range payeeSet {
		payees = append(payees, payee)
	}
	sort.Strings(payees)
	return payees
}

// GetPayeesMatching returns payees matching the given pattern
func (j *Journal) GetPayeesMatching(pattern string) []string {
	var matches []string
	pattern = strings.ToLower(pattern)
	
	payeeSet := make(map[string]bool)
	for _, tx := range j.transactions {
		if tx.Payee != "" {
			payeeLower := strings.ToLower(tx.Payee)
			if strings.Contains(payeeLower, pattern) {
				payeeSet[tx.Payee] = true
			}
		}
	}
	
	for payee := range payeeSet {
		matches = append(matches, payee)
	}
	sort.Strings(matches)
	return matches
}