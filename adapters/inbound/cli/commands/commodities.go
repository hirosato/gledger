package commands

import (
	"fmt"
	"os"

	"github.com/hirosato/gledger/application"
)

// CommoditiesCommand implements the 'commodities' command
type CommoditiesCommand struct {
	journal *application.Journal
}

// NewCommoditiesCommand creates a new commodities command
func NewCommoditiesCommand(journal *application.Journal) *CommoditiesCommand {
	return &CommoditiesCommand{
		journal: journal,
	}
}

// Execute runs the commodities command
func (c *CommoditiesCommand) Execute(args []string) error {
	var commodities []string
	
	// If account pattern is provided, filter by account
	if len(args) > 0 {
		accountPattern := args[0]
		commodities = c.journal.GetCommoditiesForAccount(accountPattern)
	} else {
		// Get all commodities
		commodities = c.journal.GetCommodities()
	}

	// Print each commodity on a new line
	for _, commodity := range commodities {
		fmt.Fprintln(os.Stdout, commodity)
	}

	return nil
}