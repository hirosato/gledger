package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hirosato/gledger/application"
)

// PricesCommand implements the 'prices' command
type PricesCommand struct {
	journal *application.Journal
}

// NewPricesCommand creates a new prices command
func NewPricesCommand(journal *application.Journal) *PricesCommand {
	return &PricesCommand{
		journal: journal,
	}
}

// PriceEntry represents a single price conversion
type PriceEntry struct {
	Date         time.Time
	FromCommodity string
	ToCommodity   string
	Price        float64
}

// Execute runs the prices command
func (c *PricesCommand) Execute(args []string) error {
	// Get commodity filter if provided
	var commodityFilter string
	if len(args) > 0 {
		commodityFilter = args[0]
	}

	// Extract prices from all transactions
	prices := c.extractPrices()

	// Filter by commodity if specified
	if commodityFilter != "" {
		filtered := []PriceEntry{}
		for _, p := range prices {
			if strings.HasPrefix(p.FromCommodity, commodityFilter) || 
			   strings.HasPrefix(p.ToCommodity, commodityFilter) {
				filtered = append(filtered, p)
			}
		}
		prices = filtered
	}

	// Sort by date and commodity
	sort.Slice(prices, func(i, j int) bool {
		if !prices[i].Date.Equal(prices[j].Date) {
			return prices[i].Date.Before(prices[j].Date)
		}
		if prices[i].FromCommodity != prices[j].FromCommodity {
			return prices[i].FromCommodity < prices[j].FromCommodity
		}
		return prices[i].ToCommodity < prices[j].ToCommodity
	})

	// Print prices in ledger format
	for _, p := range prices {
		dateStr := p.Date.Format("2006/01/02")
		
		// Format price with appropriate precision
		var priceStr string
		if p.Price == float64(int(p.Price)) {
			priceStr = fmt.Sprintf("%.2f", p.Price)
		} else {
			// Format with up to 10 decimal places, but remove trailing zeros
			priceStr = fmt.Sprintf("%.10f", p.Price)
			// Remove trailing zeros after decimal point
			priceStr = strings.TrimRight(strings.TrimRight(priceStr, "0"), ".")
		}
		
		// Format output to match ledger's spacing
		// Commodity is left-aligned with proper spacing, price is right-aligned in its field
		fmt.Fprintf(os.Stdout, "%s %-12s%12s %s\n", dateStr, p.FromCommodity, priceStr, p.ToCommodity)
	}

	return nil
}

// extractPrices extracts price information from all transactions
func (c *PricesCommand) extractPrices() []PriceEntry {
	prices := []PriceEntry{}
	seen := make(map[string]bool) // To avoid duplicates

	for _, tx := range c.journal.GetTransactions() {
		for _, posting := range tx.Postings {
			if posting.Price != nil && posting.Amount != nil && posting.Price.Amount != nil {
				// Calculate the unit price
				var unitPrice float64
				if posting.Price.IsTotal {
					// @@ means total price, divide by quantity
					quantity := posting.Amount.ToFloat64()
					totalPrice := posting.Price.Amount.ToFloat64()
					if quantity != 0 {
						unitPrice = totalPrice / quantity
					}
				} else {
					// @ means unit price
					unitPrice = posting.Price.Amount.ToFloat64()
				}

				if unitPrice != 0 {
					fromCommodity := posting.Amount.Commodity.Symbol
					toCommodity := posting.Price.Amount.Commodity.Symbol
					
					// Create a unique key to avoid duplicates
					key := fmt.Sprintf("%s-%s-%s-%.10f", 
						tx.Date.Format("2006-01-02"), 
						fromCommodity, 
						toCommodity, 
						unitPrice)
					
					if !seen[key] {
						seen[key] = true
						prices = append(prices, PriceEntry{
							Date:          tx.Date,
							FromCommodity: fromCommodity,
							ToCommodity:   toCommodity,
							Price:         unitPrice,
						})
					}
				}
			}
		}
	}

	return prices
}