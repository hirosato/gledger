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
	// Parse command line options
	var accountPattern string
	var showLotPrices bool
	var showLots bool
	var dateFormat string = "2006/01/02"
	
	// Process arguments
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--lot-prices" {
			showLotPrices = true
			i++
		} else if arg == "--lots" {
			showLots = true
			i++
		} else if arg == "--date-format" && i+1 < len(args) {
			i++
			// Convert ledger date format to Go format
			dateFormat = convertDateFormat(args[i])
			i++
		} else if !strings.HasPrefix(arg, "-") {
			// This is the account pattern
			accountPattern = arg
			i++
		} else {
			i++
		}
	}

	// Structure to track lots when needed
	type lot struct {
		amount    *domain.Amount
		price     *domain.Amount
		date      time.Time
	}
	
	// Calculate balances for all accounts
	balances := make(map[string]*domain.Balance)
	accountLots := make(map[string][]lot) // Track lots per account if needed
	
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
				
				// If we're showing lot prices or lots, track each lot separately for price specs
				if (showLotPrices || showLots) && posting.HasPrice() {
					// Store lot information
					accountLots[accountName] = append(accountLots[accountName], lot{
						amount: posting.Amount,
						price:  posting.Price.Amount,
						date:   tx.Date,
					})
				} else {
					// Normal balance tracking for all accounts
					balances[accountName].Add(posting.Amount)
				}
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
		// Check if we have lots for this account
		if lots, hasLots := accountLots[account]; hasLots && (showLotPrices || showLots) {
			// Print each lot separately
			for _, lot := range lots {
				amountStr := c.formatAmount(lot.amount)
				
				// Add lot price if requested (lots also include prices)
				if showLotPrices || showLots {
					priceStr := c.formatAmount(lot.price)
					amountStr = fmt.Sprintf("%s {%s}", amountStr, priceStr)
				}
				
				// Add lot date if requested
				if showLots {
					dateStr := lot.date.Format(dateFormat)
					amountStr = fmt.Sprintf("%s [%s]", amountStr, dateStr)
				}
				
				// Use right-aligned formatting with fixed width
				fmt.Fprintf(os.Stdout, "    %-27s%32s\n", account, amountStr)
			}
		} else {
			// Regular balance output
			balance := balances[account]
			for _, amount := range balance.GetAmounts() {
				if !amount.IsZero() {
					// Format the amount
					amountStr := c.formatAmount(amount)
					// Use right-aligned formatting with fixed width
					fmt.Fprintf(os.Stdout, "    %-27s%32s\n", account, amountStr)
				}
			}
		}
	}
	
	// Check if we're filtering to a specific account pattern
	// If we have exactly one commodity in the filtered results, we can elide the equity amount
	totalCommodities := 0
	for _, balance := range balances {
		totalCommodities += len(balance.GetAmounts())
	}
	
	shouldElideEquity := accountPattern != "" && totalCommodities == 1
	
	// Then print the offsetting Equity:Opening Balances entries
	if shouldElideEquity {
		// When filtering to a single account with one commodity, elide the amount
		fmt.Fprintf(os.Stdout, "    Equity:Opening Balances\n")
	} else {
		// First collect all equity entries
		type equityEntry struct {
			amount *domain.Amount
			text   string
		}
		var equityEntries []equityEntry
		
		// Track GBP totals when showing lots
		totalGBP := 0.0
		
		for _, account := range accounts {
			// Handle lots for equity entries
			if lots, hasLots := accountLots[account]; hasLots && (showLotPrices || showLots) {
				for _, lot := range lots {
					// Negate the amount for equity account
					negatedAmount := lot.amount.Negate()
					amountStr := c.formatAmount(negatedAmount)
					
					// Add lot price if requested (lots also include prices)
					if showLotPrices || showLots {
						priceStr := c.formatAmount(lot.price)
						amountStr = fmt.Sprintf("%s {%s}", amountStr, priceStr)
					}
					
					// Add lot date if requested
					if showLots {
						dateStr := lot.date.Format(dateFormat)
						amountStr = fmt.Sprintf("%s [%s]", amountStr, dateStr)
					}
					
					// Track GBP total for final equity entry
					marketValue := lot.price.ToFloat64() * lot.amount.ToFloat64()
					totalGBP += marketValue
					
					equityAccount := "Equity:Opening Balances"
					text := fmt.Sprintf("    %-27s%32s\n", equityAccount, amountStr)
					equityEntries = append(equityEntries, equityEntry{negatedAmount, text})
				}
			} else {
				// Regular balance entries
				balance := balances[account]
				for _, amount := range balance.GetAmounts() {
					if !amount.IsZero() {
						// Negate the amount for equity account
						negatedAmount := amount.Negate()
						amountStr := c.formatAmount(negatedAmount)
						equityAccount := "Equity:Opening Balances"
						text := fmt.Sprintf("    %-27s%32s\n", equityAccount, amountStr)
						equityEntries = append(equityEntries, equityEntry{negatedAmount, text})
					}
				}
			}
		}
		
		// Note: We don't add the GBP equity entry here because it comes from the regular balances
		// The Assets:Bank account balance provides the GBP offset
		
		// Sort equity entries: negative amounts first
		sort.Slice(equityEntries, func(i, j int) bool {
			return equityEntries[i].amount.ToFloat64() < equityEntries[j].amount.ToFloat64()
		})
		
		// Print sorted equity entries
		for _, entry := range equityEntries {
			fmt.Fprint(os.Stdout, entry.text)
		}
	}
	
	return nil
}

// formatAmount formats an amount for display
func (c *EquityCommand) formatAmount(amount *domain.Amount) string {
	// Format the number part
	var numberStr string
	floatVal := amount.ToFloat64()
	
	// Use precision from commodity or default to 2
	precision := 2
	if amount.Commodity != nil && amount.Commodity.Precision >= 0 {
		precision = amount.Commodity.Precision
	}
	
	// Always use the precision for consistency with ledger-cli
	format := fmt.Sprintf("%%.%df", precision)
	numberStr = fmt.Sprintf(format, floatVal)
	
	// Add commodity
	return numberStr + " " + amount.Commodity.Symbol
}

// convertDateFormat converts ledger date format to Go date format
func convertDateFormat(ledgerFormat string) string {
	// Simple conversion - extend as needed
	goFormat := strings.ReplaceAll(ledgerFormat, "%Y", "2006")
	goFormat = strings.ReplaceAll(goFormat, "%m", "01")
	goFormat = strings.ReplaceAll(goFormat, "%d", "02")
	return goFormat
}