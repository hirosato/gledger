package usecases

import (
	"strings"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/application/dto"
)

// ListAccountsOptions contains options for listing accounts
type ListAccountsOptions struct {
	Pattern string // Optional pattern to filter accounts
	Used    bool   // Only show accounts with transactions
}

// ListAccounts returns a list of all accounts
type ListAccounts struct {
	journal *application.Journal
}

// NewListAccounts creates a new ListAccounts use case
func NewListAccounts(journal *application.Journal) *ListAccounts {
	return &ListAccounts{
		journal: journal,
	}
}

// Execute returns a list of accounts matching the criteria
func (la *ListAccounts) Execute(options ListAccountsOptions) (*dto.AccountList, error) {
	accounts := la.journal.GetAccounts()
	
	result := &dto.AccountList{
		Accounts: []dto.AccountInfo{},
	}
	
	for _, accountName := range accounts {
		if la.shouldInclude(accountName, options) {
			// Extract the last part of the account name
			parts := strings.Split(accountName, ":")
			name := parts[len(parts)-1]
			level := len(parts) - 1
			
			var parent string
			if len(parts) > 1 {
				parent = strings.Join(parts[:len(parts)-1], ":")
			}
			
			result.Accounts = append(result.Accounts, dto.AccountInfo{
				FullName:        accountName,
				Name:            name,
				Level:           level,
				Parent:          parent,
				HasTransactions: true, // All returned accounts have transactions
			})
		}
	}
	
	return result, nil
}

func (la *ListAccounts) shouldInclude(accountName string, options ListAccountsOptions) bool {
	// Implement filtering logic
	if options.Pattern != "" {
		return strings.Contains(strings.ToLower(accountName), strings.ToLower(options.Pattern))
	}
	return true
}