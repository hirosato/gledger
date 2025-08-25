package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/hirosato/gledger/application"
)

// StatsCommand implements the 'stats' command
type StatsCommand struct {
	journal *application.Journal
}

// NewStatsCommand creates a new stats command
func NewStatsCommand(journal *application.Journal) *StatsCommand {
	return &StatsCommand{
		journal: journal,
	}
}

// Execute runs the stats command
func (c *StatsCommand) Execute(args []string) error {
	transactions := c.journal.GetTransactions()
	
	if len(transactions) == 0 {
		fmt.Fprintln(os.Stdout, "No transactions found")
		return nil
	}

	// Calculate basic statistics
	var earliestDate, latestDate time.Time
	var totalPostings int
	var unclearedPostings int
	payees := make(map[string]bool)
	accounts := make(map[string]bool)

	for i, tx := range transactions {
		if i == 0 {
			earliestDate = tx.Date
			latestDate = tx.Date
		} else {
			if tx.Date.Before(earliestDate) {
				earliestDate = tx.Date
			}
			if tx.Date.After(latestDate) {
				latestDate = tx.Date
			}
		}

		if tx.Payee != "" {
			payees[tx.Payee] = true
		}

		for _, posting := range tx.Postings {
			totalPostings++
			if posting.Account != nil {
				accounts[posting.Account.FullName] = true
			}
			// Count uncleared (pending) postings
			if tx.Status == 0 { // TransactionStatusPending
				unclearedPostings++
			}
		}
	}

	// Calculate time period
	days := int(latestDate.Sub(earliestDate).Hours()/24) + 1
	postsPerDay := float64(totalPostings) / float64(days)

	// Output statistics
	fmt.Fprintf(os.Stdout, "Time period: %s to %s (%d days)\n",
		earliestDate.Format("06-Jan-02"),
		latestDate.Format("06-Jan-02"),
		days)
	fmt.Fprintln(os.Stdout)

	fmt.Fprintln(os.Stdout, "  Files these postings came from:")
	fmt.Fprintln(os.Stdout, "    [journal file]") // Placeholder
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "  Unique payees:               %d\n", len(payees))
	fmt.Fprintf(os.Stdout, "  Unique accounts:             %d\n", len(accounts))
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "  Number of postings:          %d (%.2f per day)\n", totalPostings, postsPerDay)
	fmt.Fprintf(os.Stdout, "  Uncleared postings:          %d\n", unclearedPostings)
	fmt.Fprintln(os.Stdout)

	// Calculate days since last post
	now := time.Now()
	daysSinceLastPost := int(now.Sub(latestDate).Hours() / 24)
	fmt.Fprintf(os.Stdout, "  Days since last post:        %d\n", daysSinceLastPost)

	// Count posts in time periods (simplified implementation)
	last7Days := 0
	last30Days := 0
	thisMonth := 0

	cutoff7 := now.AddDate(0, 0, -7)
	cutoff30 := now.AddDate(0, 0, -30)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for _, tx := range transactions {
		if tx.Date.After(cutoff7) || tx.Date.Equal(cutoff7) {
			last7Days += len(tx.Postings)
		}
		if tx.Date.After(cutoff30) || tx.Date.Equal(cutoff30) {
			last30Days += len(tx.Postings)
		}
		if tx.Date.After(firstOfMonth) || tx.Date.Equal(firstOfMonth) {
			thisMonth += len(tx.Postings)
		}
	}

	fmt.Fprintf(os.Stdout, "  Posts in last 7 days:        %d\n", last7Days)
	fmt.Fprintf(os.Stdout, "  Posts in last 30 days:       %d\n", last30Days)
	fmt.Fprintf(os.Stdout, "  Posts seen this month:       %d\n", thisMonth)

	return nil
}