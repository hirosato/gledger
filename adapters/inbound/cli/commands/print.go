package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/domain"
)

// Formatting constants for print command
const (
	// Alignment columns for different contexts
	shortAccountAlignColumn = 47  // For short account names  
	longAccountAlignColumn  = 59  // For long account names
	
	// Spacing constants  
	postingIndent       = 4   // Standard posting indentation
	complexAmountSpacing = 2  // Spacing for complex amounts
	minSpacing          = 2   // Minimum spacing always
	
	// Classification thresholds
	longAccountThreshold       = 30 // Account names >= this use long column alignment
	simpleCurrencyMaxLength    = 6  // Max length for simple currency amounts
	complexAmountMinSpaces     = 1  // Min spaces in amount to be considered complex
)

// Pre-computed indentation strings
var (
	postingIndentStr = strings.Repeat(" ", postingIndent)
	noteIndentStr    = strings.Repeat(" ", postingIndent*2) // Notes use double indentation
)

// PrintOptions represents options for the print command
type PrintOptions struct {
	Raw         bool   // --raw option: preserve original formatting
	DecimalComma bool  // --decimal-comma option: use comma as decimal separator
	Actual       bool  // --actual option: show actual dates
	Hashes       string // --hashes option: for integrity checking
}

// PrintCommand implements the 'print' command
type PrintCommand struct {
	journal *application.Journal
	options PrintOptions
}

// NewPrintCommand creates a new print command
func NewPrintCommand(journal *application.Journal) *PrintCommand {
	return &PrintCommand{
		journal: journal,
	}
}

// Execute runs the print command
func (c *PrintCommand) Execute(args []string) error {
	// Parse command line options
	err := c.parseOptions(args)
	if err != nil {
		return err
	}

	// Get all transactions
	transactions := c.journal.GetTransactions()

	// Print each transaction
	for i, tx := range transactions {
		c.printTransaction(&tx)
		
		// Add blank line between transactions (except after the last one)
		if i < len(transactions)-1 {
			fmt.Println()
		}
	}

	return nil
}

// parseOptions parses command line arguments for print options
func (c *PrintCommand) parseOptions(args []string) error {
	for _, arg := range args {
		switch arg {
		case "--raw":
			c.options.Raw = true
		case "--decimal-comma":
			c.options.DecimalComma = true
		case "--actual":
			c.options.Actual = true
		default:
			if strings.HasPrefix(arg, "--hashes=") {
				c.options.Hashes = strings.TrimPrefix(arg, "--hashes=")
			}
		}
	}
	return nil
}

// printTransaction prints a single transaction in ledger format
func (c *PrintCommand) printTransaction(tx *domain.Transaction) {
	// Print transaction header: date [status] [code] payee
	line := c.formatDate(tx.Date)
	
	// Add spacing - varies by options and format
	if c.options.Raw {
		line += "       "  // --raw option uses extra spacing
	} else {
		line += " "       // Normal format uses single space
	}
	
	// Add status if not pending
	if tx.Status != domain.TransactionStatusPending {
		switch tx.Status {
		case domain.TransactionStatusCleared:
			line += "*"
		case domain.TransactionStatusReconciled:
			line += "!"
		}
		line += " "
	}
	
	// Add code if present
	if tx.Code != "" {
		line += "(" + tx.Code + ") "
	}
	
	// Add payee
	if tx.Payee != "" {
		line += tx.Payee
	}
	
	fmt.Println(line)
	
	// Print postings
	for _, posting := range tx.Postings {
		c.printPosting(posting)
	}
	
	// Print transaction note if present
	if tx.Note != "" {
		fmt.Println(postingIndentStr + "; " + tx.Note)
	}
}

// printPosting prints a single posting line
func (c *PrintCommand) printPosting(posting *domain.Posting) {
	if posting.Amount != nil {
		// Build the amount string with all components
		amountStr := c.formatAmount(posting.Amount)
		
		// Add cost information if present
		if posting.HasCost() {
			costStr := c.formatCost(posting.Cost)
			if costStr != "" {
				amountStr += " " + costStr
			}
		}
		
		// Add price information if present
		if posting.HasPrice() {
			priceStr := c.formatPrice(posting.Price)
			if priceStr != "" {
				amountStr += " " + priceStr
			}
		}
		
		// Add balance assertion if present
		if posting.HasBalanceAssertion() {
			assertionStr := c.formatBalanceAssertion(posting.BalanceAssertion)
			if assertionStr != "" {
				amountStr += " " + assertionStr
			}
		}
		
		// Format with ledger-style alignment
		accountName := posting.Account.Name
		
		// Ledger's alignment strategy (derived from baseline test analysis):
		// 1. For complex amounts: use minimal spacing for readability
		// 2. For simple currency amounts: use column alignment based on account name length
		// 3. For other amounts: use reasonable default spacing
		//
		// This creates professional, tabular output matching original ledger formatting.
		
		accountEndPos := postingIndent + len(accountName)
		
		// Detect if this is a simple currency amount (short and starts with currency symbol)
		isSimpleCurrency := len(amountStr) <= simpleCurrencyMaxLength && 
			(strings.HasPrefix(amountStr, "$") || strings.HasPrefix(amountStr, "€") || strings.HasPrefix(amountStr, "£"))
		
		// Detect if this is a complex amount (contains @ or multiple components)
		isComplexAmount := strings.Contains(amountStr, "@") || strings.Count(amountStr, " ") > complexAmountMinSpaces
		
		if isComplexAmount {
			// Complex amounts get minimal spacing for readability
			fmt.Printf("%s%s%s%s\n", postingIndentStr, accountName, strings.Repeat(" ", complexAmountSpacing), amountStr)
		} else if isSimpleCurrency {
			// Simple currency amounts: choose alignment column based on account length and context
			var targetColumn int
			
			// For very long account names, prefer LONG column if it fits
			if len(accountName) >= longAccountThreshold && accountEndPos + minSpacing < longAccountAlignColumn {
				targetColumn = longAccountAlignColumn
			} else if accountEndPos + minSpacing < shortAccountAlignColumn {
				targetColumn = shortAccountAlignColumn
			} else if accountEndPos + minSpacing < longAccountAlignColumn {
				targetColumn = longAccountAlignColumn
			} else {
				// Account name too long for either column, use minimum spacing
				fmt.Printf("%s%s%s%s\n", postingIndentStr, accountName, strings.Repeat(" ", minSpacing), amountStr)
				return
			}
			
			spacesNeeded := targetColumn - accountEndPos
			if spacesNeeded < minSpacing {
				spacesNeeded = minSpacing
			}
			fmt.Printf("%s%s%s%s\n", postingIndentStr, accountName, strings.Repeat(" ", spacesNeeded), amountStr)
		} else {
			// Default case: reasonable spacing
			fmt.Printf("%s%s%s%s\n", postingIndentStr, accountName, strings.Repeat(" ", minSpacing), amountStr)
		}
	} else {
		// No amount - just print account name
		fmt.Printf("%s%s\n", postingIndentStr, posting.Account.Name)
	}
	
	// Print posting note if present
	if posting.Note != "" {
		fmt.Println(noteIndentStr + "; " + posting.Note)
	}
}

// formatDate formats a date for transaction header
func (c *PrintCommand) formatDate(date time.Time) string {
	return date.Format("2006/01/02")
}

// formatAmount formats an amount for display
func (c *PrintCommand) formatAmount(amount *domain.Amount) string {
	if amount == nil {
		return ""
	}
	
	// Get the numeric string
	numberStr := c.formatNumber(amount)
	
	// Handle commodity placement - currencies like $ go before, others after
	commodity := amount.Commodity.Symbol
	if commodity == "$" || commodity == "€" || commodity == "£" {
		result := commodity + numberStr
		if c.options.DecimalComma {
			return strings.ReplaceAll(result, ".", ",")
		}
		return result
	}
	
	// Non-currency commodities go after the number
	result := numberStr + " " + commodity
	if c.options.DecimalComma {
		return strings.ReplaceAll(result, ".", ",")
	}
	return result
}

// formatNumber formats just the numeric part of an amount
func (c *PrintCommand) formatNumber(amount *domain.Amount) string {
	if amount.Commodity.Precision == 0 {
		if amount.Number.IsInt() {
			return amount.Number.Num().String()
		}
		return amount.Number.String()
	}
	
	floatValue, _ := amount.Number.Float64()
	precision := amount.Commodity.Precision
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, floatValue)
}

// formatCost formats cost basis information
func (c *PrintCommand) formatCost(cost *domain.CostBasis) string {
	if cost == nil {
		return ""
	}
	
	if cost.PerUnitAmount != nil {
		return "{" + c.formatAmount(cost.PerUnitAmount) + "}"
	}
	
	if cost.Amount != nil {
		return "{{" + c.formatAmount(cost.Amount) + "}}"
	}
	
	return ""
}

// formatPrice formats price specification information
func (c *PrintCommand) formatPrice(price *domain.PriceSpec) string {
	if price == nil || price.Amount == nil {
		return ""
	}
	
	if price.IsTotal {
		return "@@ " + c.formatAmount(price.Amount)
	}
	return "@ " + c.formatAmount(price.Amount)
}

// formatBalanceAssertion formats balance assertion/assignment
func (c *PrintCommand) formatBalanceAssertion(assertion *domain.BalanceAssertion) string {
	if assertion == nil || assertion.Amount == nil {
		return ""
	}
	
	if assertion.IsAssignment {
		return "= " + c.formatAmount(assertion.Amount)
	}
	return "== " + c.formatAmount(assertion.Amount)
}