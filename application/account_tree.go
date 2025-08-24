package application

import (
	"sort"
	"strings"

	"github.com/hirosato/gledger/domain"
)

// AccountTree manages hierarchical account structure
type AccountTree struct {
	root     *domain.Account
	accounts map[string]*domain.Account
}

// NewAccountTree creates a new account tree
func NewAccountTree() *AccountTree {
	return &AccountTree{
		root:     domain.NewAccount(""),
		accounts: make(map[string]*domain.Account),
	}
}

// AddAccount adds an account to the tree
func (t *AccountTree) AddAccount(account *domain.Account) {
	fullName := account.FullName
	if fullName == "" {
		return
	}

	// Register this account
	t.accounts[fullName] = account

	// Create parent accounts if they don't exist
	parts := strings.Split(fullName, ":")
	for i := 1; i <= len(parts); i++ {
		parentName := strings.Join(parts[:i], ":")
		if _, exists := t.accounts[parentName]; !exists {
			parentAccount := domain.NewAccount(parts[i-1])
			parentAccount.FullName = parentName
			t.accounts[parentName] = parentAccount
		}
	}
}

// GetAccount retrieves an account by name
func (t *AccountTree) GetAccount(name string) (*domain.Account, bool) {
	account, exists := t.accounts[name]
	return account, exists
}

// GetAllAccounts returns all accounts sorted by name
func (t *AccountTree) GetAllAccounts() []string {
	names := make([]string, 0, len(t.accounts))
	for name := range t.accounts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetAccountsMatching returns accounts matching a pattern
func (t *AccountTree) GetAccountsMatching(pattern string) []string {
	var matches []string
	pattern = strings.ToLower(pattern)
	
	for name := range t.accounts {
		if strings.Contains(strings.ToLower(name), pattern) {
			matches = append(matches, name)
		}
	}
	
	sort.Strings(matches)
	return matches
}