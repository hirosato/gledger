package commands

import (
	"fmt"
	"os"

	"github.com/hirosato/gledger/application"
)

// AccountsCommand implements the 'accounts' command
type AccountsCommand struct {
	journal *application.Journal
}

// NewAccountsCommand creates a new accounts command
func NewAccountsCommand(journal *application.Journal) *AccountsCommand {
	return &AccountsCommand{
		journal: journal,
	}
}

// Execute runs the accounts command
func (c *AccountsCommand) Execute(args []string) error {
	var accounts []string
	
	// If pattern is provided, filter accounts
	if len(args) > 0 {
		pattern := args[0]
		accounts = c.journal.GetAccountsMatching(pattern)
	} else {
		// Get all accounts
		accounts = c.journal.GetAccounts()
	}

	// Print each account on a new line
	for _, account := range accounts {
		fmt.Fprintln(os.Stdout, account)
	}

	return nil
}