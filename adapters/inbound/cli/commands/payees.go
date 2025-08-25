package commands

import (
	"fmt"
	"os"

	"github.com/hirosato/gledger/application"
)

// PayeesCommand implements the 'payees' command
type PayeesCommand struct {
	journal *application.Journal
}

// NewPayeesCommand creates a new payees command
func NewPayeesCommand(journal *application.Journal) *PayeesCommand {
	return &PayeesCommand{
		journal: journal,
	}
}

// Execute runs the payees command
func (c *PayeesCommand) Execute(args []string) error {
	var payees []string
	
	// If pattern is provided, filter payees
	if len(args) > 0 {
		pattern := args[0]
		payees = c.journal.GetPayeesMatching(pattern)
	} else {
		// Get all payees
		payees = c.journal.GetPayees()
	}

	// Print each payee on a new line
	for _, payee := range payees {
		fmt.Fprintln(os.Stdout, payee)
	}

	return nil
}