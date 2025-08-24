package usecases

import (
	"strings"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/application/dto"
	"github.com/hirosato/gledger/domain"
)

// GetBalanceOptions contains options for the balance calculation
type GetBalanceOptions struct {
	Flat     bool     // Show accounts in flat format
	NoTotal  bool     // Don't show total line
	Empty    bool     // Show accounts with zero balance
	NoRollup bool     // Don't roll up account balances to parents
	Accounts []string // Filter by account patterns
}

// GetBalance calculates and returns account balances
type GetBalance struct {
	journal *application.Journal
}

// NewGetBalance creates a new GetBalance use case
func NewGetBalance(journal *application.Journal) *GetBalance {
	return &GetBalance{
		journal: journal,
	}
}

// Execute performs the balance calculation and returns a BalanceReport
func (gb *GetBalance) Execute(options GetBalanceOptions) (*dto.BalanceReport, error) {
	report := &dto.BalanceReport{
		Accounts: []dto.AccountBalance{},
	}

	// Calculate balances for all accounts
	balances := gb.calculateBalances(options)

	// Convert to DTOs
	for _, acc := range balances {
		if !options.Empty && acc.Balance.IsZero() {
			continue
		}

		report.Accounts = append(report.Accounts, dto.AccountBalance{
			Name:    acc.Name,
			Balance: acc.Balance.String(),
			Level:   acc.Level,
			IsEmpty: acc.Balance.IsZero(),
		})
	}

	// Add total if needed
	if !options.NoTotal && len(report.Accounts) > 0 {
		total := gb.calculateTotal(balances)
		report.Total = &dto.AccountBalance{
			Name:    "Total",
			Balance: total.String(),
			IsTotal: true,
		}
	}

	return report, nil
}

// calculateBalances calculates balances for all accounts
func (gb *GetBalance) calculateBalances(options GetBalanceOptions) []accountBalance {
	// Get all account names
	accounts := gb.journal.GetAccounts()
	
	var balances []accountBalance
	for _, accountName := range accounts {
		if gb.shouldIncludeAccountName(accountName, options.Accounts) {
			// Calculate balance for this account
			bal := gb.journal.GetBalance(accountName)
			
			// Calculate account level (depth)
			level := len(strings.Split(accountName, ":")) - 1
			
			balances = append(balances, accountBalance{
				Name:    accountName,
				Balance: *bal,
				Level:   level,
			})
		}
	}
	
	return balances
}

// shouldIncludeAccountName checks if an account should be included based on filters
func (gb *GetBalance) shouldIncludeAccountName(accountName string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	// Implement pattern matching logic
	for _, pattern := range patterns {
		if strings.Contains(accountName, pattern) {
			return true
		}
	}
	return false
}

// calculateTotal calculates the total of all balances
func (gb *GetBalance) calculateTotal(balances []accountBalance) domain.Balance {
	total := domain.NewBalance()
	for _, bal := range balances {
		// Add balance to total
		for _, amount := range bal.Balance.GetAmounts() {
			total.Add(amount)
		}
	}
	return *total
}

type accountBalance struct {
	Name    string
	Balance domain.Balance
	Level   int
}