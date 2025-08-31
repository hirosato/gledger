package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/domain"
)

// EquityCommand implements the 'equity' command
type EquityCommand struct {
	journal *application.Journal
}

// NewEquityCommand creates a new equity command
func NewEquityCommand(journal *application.Journal) *EquityCommand {
	return &EquityCommand{
		journal: journal,
	}
}

// Execute runs the equity command
func (c *EquityCommand) Execute(args []string) error {
	// Get account pattern if provided
	var accountPattern string
	if len(args) > 0 {
		accountPattern = args[0]
	}

	// Calculate balances for all accounts
	balances := make(map[string]*domain.Balance)
	
	for _, tx := range c.journal.GetTransactions() {
		for _, posting := range tx.Postings {
			accountName := posting.Account.Name
			
			// Filter by account pattern if specified
			if accountPattern != "" && !strings.Contains(strings.ToLower(accountName), strings.ToLower(accountPattern)) {
				continue
			}
			
			// For equity, we track the actual commodity amounts
			if posting.Amount != nil {
				if balances[accountName] == nil {
					balances[accountName] = domain.NewBalance()
				}
				balances[accountName].Add(posting.Amount)
			}
		}
	}

	// Get the latest transaction date to use for the equity entry
	var latestDate time.Time
	for _, tx := range c.journal.GetTransactions() {
		if tx.Date.After(latestDate) {
			latestDate = tx.Date
		}
	}
	
	// If no transactions, use today's date
	if latestDate.IsZero() {
		latestDate = time.Now()
	}

	// Print the equity transaction
	fmt.Fprintf(os.Stdout, "%s Opening Balances\n", latestDate.Format("2006/01/02"))
	
	// Sort accounts for consistent output
	var accounts []string
	for account := range balances {
		accounts = append(accounts, account)
	}
	sort.Strings(accounts)
	
	// First print all regular account balances
	for _, account := range accounts {
		balance := balances[account]
		for _, amount := range balance.GetAmounts() {
			if !amount.IsZero() {
				// Format the amount
				amountStr := c.formatAmount(amount)
				// Calculate spacing
				spacing := 59 - len(account) - len(amountStr)
				if spacing < 2 {
					spacing = 2
				}
				fmt.Fprintf(os.Stdout, "    %s%*s\n", account, spacing, amountStr)
			}
		}
	}
	
	// Then print the offsetting Equity:Opening Balances entries
	// First collect all equity entries
	type equityEntry struct {
		amount *domain.Amount
		text   string
	}
	var equityEntries []equityEntry
	
	for _, account := range accounts {
		balance := balances[account]
		for _, amount := range balance.GetAmounts() {
			if !amount.IsZero() {
				// Negate the amount for equity account
				negatedAmount := amount.Negate()
				amountStr := c.formatAmount(negatedAmount)
				equityAccount := "Equity:Opening Balances"
				spacing := 59 - len(equityAccount) - len(amountStr)
				if spacing < 2 {
					spacing = 2
				}
				text := fmt.Sprintf("    %s%*s\n", equityAccount, spacing, amountStr)
				equityEntries = append(equityEntries, equityEntry{negatedAmount, text})
			}
		}
	}
	
	// Sort equity entries: negative amounts first
	sort.Slice(equityEntries, func(i, j int) bool {
		return equityEntries[i].amount.ToFloat64() < equityEntries[j].amount.ToFloat64()
	})
	
	// Print sorted equity entries
	for _, entry := range equityEntries {
		fmt.Fprint(os.Stdout, entry.text)
	}
	
	return nil
}

// formatAmount formats an amount for display
func (c *EquityCommand) formatAmount(amount *domain.Amount) string {
	// Format the number part
	var numberStr string
	floatVal := amount.ToFloat64()
	
	// Check if it's a whole number
	if floatVal == float64(int(floatVal)) {
		numberStr = fmt.Sprintf("%d", int(floatVal))
	} else {
		// Use precision from commodity or default to 2
		precision := 2
		if amount.Commodity != nil && amount.Commodity.Precision > 0 {
			precision = amount.Commodity.Precision
		}
		format := fmt.Sprintf("%%.%df", precision)
		numberStr = fmt.Sprintf(format, floatVal)
	}
	
	// Add commodity
	return numberStr + " " + amount.Commodity.Symbol
}