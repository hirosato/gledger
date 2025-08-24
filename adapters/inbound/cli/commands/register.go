package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/domain"
)

// RegisterOptions represents options for the register command
type RegisterOptions struct {
	AccountFilter string // Account pattern filter (e.g., :inve for Assets:Investment)
}

// RegisterCommand implements the 'register' command
type RegisterCommand struct {
	journal *application.Journal
	options RegisterOptions
}

// NewRegisterCommand creates a new register command
func NewRegisterCommand(journal *application.Journal) *RegisterCommand {
	return &RegisterCommand{
		journal: journal,
	}
}

// Execute runs the register command
func (c *RegisterCommand) Execute(args []string) error {
	// Parse command line options
	err := c.parseOptions(args)
	if err != nil {
		return err
	}

	// Get all transactions
	transactions := c.journal.GetTransactions()

	// Track running balances
	runningBalance := domain.NewBalance()

	// Process each transaction
	for _, tx := range transactions {
		c.displayTransaction(&tx, runningBalance)
	}

	return nil
}

// parseOptions parses command line arguments for register options
func (c *RegisterCommand) parseOptions(args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, ":") {
			// Account filter
			c.options.AccountFilter = arg[1:] // Remove the ':' prefix
		}
	}
	return nil
}

// displayTransaction displays a transaction in register format
func (c *RegisterCommand) displayTransaction(tx *domain.Transaction, runningBalance *domain.Balance) {
	// Filter postings based on account filter if specified
	var postingsToShow []*domain.Posting
	if c.options.AccountFilter != "" {
		for _, posting := range tx.Postings {
			if c.matchesAccountFilter(posting.Account.Name, c.options.AccountFilter) {
				postingsToShow = append(postingsToShow, posting)
			}
		}
		if len(postingsToShow) == 0 {
			return // No matching postings, skip this transaction
		}
	} else {
		postingsToShow = tx.Postings
	}

	// Format date (12-Jan-10 format)
	dateStr := c.formatDate(tx.Date)
	
	// Format description (truncated to ~20 chars)
	descStr := c.formatDescription(tx.Payee)

	// Display first posting with date and description
	if len(postingsToShow) > 0 {
		first := postingsToShow[0]
		amountStr := c.formatAmount(first.Amount)
		
		// Update running balance
		runningBalance.Add(first.Amount)
		runningBalanceStr := c.formatBalance(runningBalance)

		fmt.Fprintf(os.Stdout, "%-9s %-20s %-30s %12s %12s\n", 
			dateStr, descStr, first.Account.Name, amountStr, runningBalanceStr)
		
		// Display additional balance lines for multi-commodity
		c.displayAdditionalBalanceLines(runningBalance)

		// Display additional postings (indented, no date/description)
		for _, posting := range postingsToShow[1:] {
			amountStr := c.formatAmount(posting.Amount)
			runningBalance.Add(posting.Amount)
			runningBalanceStr := c.formatBalance(runningBalance)

			fmt.Fprintf(os.Stdout, "%-9s %-20s %-30s %12s %12s\n", 
				"", "", posting.Account.Name, amountStr, runningBalanceStr)
			
			// Display additional balance lines for multi-commodity
			c.displayAdditionalBalanceLines(runningBalance)
		}

		// If we're filtering and showing only some postings, we need to account for 
		// the unseen postings in the running balance
		if c.options.AccountFilter != "" {
			for _, posting := range tx.Postings {
				if !c.matchesAccountFilter(posting.Account.Name, c.options.AccountFilter) {
					runningBalance.Add(posting.Amount)
				}
			}
		}
	}
}

// matchesAccountFilter checks if an account name matches the filter
func (c *RegisterCommand) matchesAccountFilter(accountName, filter string) bool {
	// Simple substring matching for now (matches ":inve" pattern)
	return strings.Contains(strings.ToLower(accountName), strings.ToLower(filter))
}

// formatDate formats a date in register format (12-Jan-10)
func (c *RegisterCommand) formatDate(date time.Time) string {
	months := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	
	year := date.Year() % 100 // Last two digits
	return fmt.Sprintf("%02d-%s-%02d", year, months[date.Month()], date.Day())
}

// formatDescription formats a transaction description (truncated)
func (c *RegisterCommand) formatDescription(payee string) string {
	if len(payee) > 20 {
		return payee[:17] + ".."
	}
	return payee
}

// formatAmount formats an amount for display
func (c *RegisterCommand) formatAmount(amount *domain.Amount) string {
	if amount.Commodity.Symbol == "$" {
		return fmt.Sprintf("%.0f", amount.ToFloat64())
	}
	return amount.Format(true)
}

// formatBalance formats a balance for display (returns first commodity only)
func (c *RegisterCommand) formatBalance(balance *domain.Balance) string {
	if balance.IsZero() {
		return "0"
	}

	amounts := balance.GetAmounts()
	if len(amounts) > 0 {
		amount := amounts[0] // Show only first commodity in main line
		if amount.Commodity.Symbol == "$" {
			return fmt.Sprintf("%.0f", amount.ToFloat64())
		}
		return amount.Format(true)
	}
	
	return "0"
}

// displayAdditionalBalanceLines displays additional balance lines for multi-commodity balances
func (c *RegisterCommand) displayAdditionalBalanceLines(balance *domain.Balance) {
	if balance.IsZero() {
		return
	}

	amounts := balance.GetAmounts()
	// Display remaining commodities (skip the first one which was already shown)
	for i := 1; i < len(amounts); i++ {
		amount := amounts[i]
		amountStr := ""
		if amount.Commodity.Symbol == "$" {
			amountStr = fmt.Sprintf("%.0f", amount.ToFloat64())
		} else {
			amountStr = amount.Format(true)
		}
		
		fmt.Fprintf(os.Stdout, "%-9s %-20s %-30s %12s %12s\n", 
			"", "", "", "", amountStr)
	}
}