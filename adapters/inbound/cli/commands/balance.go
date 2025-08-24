package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/domain"
)

// BalanceOptions represents options for the balance command
type BalanceOptions struct {
	Flat     bool // --flat: Show accounts in flat format
	NoTotal  bool // --no-total: Don't show total line
	Empty    bool // -E, --empty: Show accounts with zero balance
	NoRollup bool // -n: Don't roll up account balances to parents
}

// BalanceCommand implements the 'balance' command
type BalanceCommand struct {
	journal *application.Journal
	options BalanceOptions
}

// NewBalanceCommand creates a new balance command
func NewBalanceCommand(journal *application.Journal) *BalanceCommand {
	return &BalanceCommand{
		journal: journal,
	}
}

// Execute runs the balance command
func (c *BalanceCommand) Execute(args []string) error {
	// Parse command line options
	err := c.parseOptions(args)
	if err != nil {
		return err
	}

	// Get account balances
	balances := c.getAccountBalances()

	// Filter based on options
	balances = c.filterBalances(balances)

	// Format and display
	c.displayBalances(balances)

	return nil
}

// parseOptions parses command line arguments for balance options
func (c *BalanceCommand) parseOptions(args []string) error {
	for _, arg := range args {
		switch arg {
		case "--flat":
			c.options.Flat = true
		case "--no-total":
			c.options.NoTotal = true
		case "-E", "--empty":
			c.options.Empty = true
		case "-n", "--no-rollup":
			c.options.NoRollup = true
		}
	}
	return nil
}

// AccountBalance represents balance information for an account
type AccountBalance struct {
	Account string
	Balance *domain.Balance
	Depth   int
}

// getAccountBalances calculates balances for all accounts
func (c *BalanceCommand) getAccountBalances() []AccountBalance {
	var balances []AccountBalance

	// Special case: -n --flat should show nothing
	if c.options.NoRollup && c.options.Flat {
		return balances // return empty slice
	}

	if c.options.NoRollup {
		// Don't roll up - show only parent accounts with balances
		balances = c.getParentAccountBalances()
	} else {
		// Roll up balances to parent accounts
		balances = c.getRolledUpBalances()
	}

	return balances
}

// getParentAccountBalances returns only top-level parent accounts for -n flag
func (c *BalanceCommand) getParentAccountBalances() []AccountBalance {
	var balances []AccountBalance
	parentBalances := make(map[string]*domain.Balance)

	// Calculate balances for all accounts used in transactions
	accounts := c.journal.GetAccounts()
	for _, account := range accounts {
		balance := c.journal.GetBalance(account)
		
		// Get the top-level parent account
		parts := strings.Split(account, ":")
		parentAccount := parts[0]
		
		if _, exists := parentBalances[parentAccount]; !exists {
			parentBalances[parentAccount] = domain.NewBalance()
		}
		parentBalances[parentAccount].AddBalance(balance)
	}

	// Convert to slice and filter
	for account, balance := range parentBalances {
		if !balance.IsZero() || c.options.Empty {
			balances = append(balances, AccountBalance{
				Account: account,
				Balance: balance,
				Depth:   1,
			})
		}
	}

	return balances
}

// getRolledUpBalances calculates balances with parent rollups
func (c *BalanceCommand) getRolledUpBalances() []AccountBalance {
	var balances []AccountBalance
	
	// Calculate balances for all accounts used in transactions
	accounts := c.journal.GetAccounts()
	
	// Get balances for leaf accounts
	leafAccountBalances := make(map[string]*domain.Balance)
	for _, account := range accounts {
		balance := c.journal.GetLeafBalance(account)
		leafAccountBalances[account] = balance
	}

	// In flat mode, only show leaf accounts that have non-zero balances
	if c.options.Flat {
		for account, balance := range leafAccountBalances {
			if !balance.IsZero() || c.options.Empty {
				balances = append(balances, AccountBalance{
					Account: account,
					Balance: balance,
					Depth:   c.getAccountDepth(account),
				})
			}
		}
		return balances
	}

	// For hierarchical mode, determine what should be shown
	accountBalances := make(map[string]*domain.Balance)
	
	// First, identify which accounts have siblings (same parent)
	parentChildren := make(map[string][]string)
	for account := range leafAccountBalances {
		parts := strings.Split(account, ":")
		if len(parts) > 1 {
			parent := parts[0]
			parentChildren[parent] = append(parentChildren[parent], account)
		}
	}
	
	// For each leaf account, decide how to display it
	for account, balance := range leafAccountBalances {
		parts := strings.Split(account, ":")
		
		if len(parts) == 1 {
			// Top-level account, just add it
			accountBalances[account] = balance
		} else {
			parent := parts[0]
			
			// If this account has siblings under the same parent, create parent rollup
			if len(parentChildren[parent]) > 1 {
				// Add parent account
				if _, exists := accountBalances[parent]; !exists {
					accountBalances[parent] = domain.NewBalance()
				}
				accountBalances[parent].AddBalance(balance)
				
				// Add child account
				accountBalances[account] = balance
			} else {
				// No siblings, show full account name (don't create parent)
				accountBalances[account] = balance
			}
		}
	}

	// Convert to slice and filter
	for account, balance := range accountBalances {
		if !balance.IsZero() || c.options.Empty {
			// Calculate depth based on whether this account has a parent in the results
			depth := 1
			parts := strings.Split(account, ":")
			if len(parts) > 1 {
				parent := parts[0]
				if _, exists := accountBalances[parent]; exists {
					// Parent exists in results, this is a child
					depth = 2
				} else {
					// Parent doesn't exist in results, this is top-level
					depth = 1
				}
			}
			
			balances = append(balances, AccountBalance{
				Account: account,
				Balance: balance,
				Depth:   depth,
			})
		}
	}

	return balances
}

// hasParentInResults checks if an account has a parent account in the results
func (c *BalanceCommand) hasParentInResults(account string, accountBalances map[string]*domain.Balance) bool {
	parts := strings.Split(account, ":")
	for i := 1; i < len(parts); i++ {
		parentAccount := strings.Join(parts[:i], ":")
		if _, exists := accountBalances[parentAccount]; exists {
			return true
		}
	}
	return false
}

// getAccountDepth returns the depth of an account in the hierarchy
func (c *BalanceCommand) getAccountDepth(account string) int {
	return len(strings.Split(account, ":"))
}

// filterBalances applies filtering options
func (c *BalanceCommand) filterBalances(balances []AccountBalance) []AccountBalance {
	var filtered []AccountBalance

	for _, bal := range balances {
		// Apply empty filter
		if !c.options.Empty && bal.Balance.IsZero() {
			continue
		}
		filtered = append(filtered, bal)
	}

	// Sort accounts
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Account < filtered[j].Account
	})

	return filtered
}

// displayBalances formats and displays the balance report
func (c *BalanceCommand) displayBalances(balances []AccountBalance) {
	if len(balances) == 0 {
		return
	}

	// Calculate total from leaf accounts only (to avoid double counting)
	total := domain.NewBalance()
	leafAccounts := make(map[string]bool)
	
	// First, identify leaf accounts (accounts that appear as full names in our results)
	// In hierarchical view, a leaf account either has no parent shown, or it's explicitly marked as a leaf
	for _, bal := range balances {
		parts := strings.Split(bal.Account, ":")
		if len(parts) > 1 {
			// This is potentially a child account
			parentAccount := parts[0]
			isParentShown := false
			for _, other := range balances {
				if other.Account == parentAccount {
					isParentShown = true
					break
				}
			}
			if isParentShown {
				// This is a child account, count it as leaf
				leafAccounts[bal.Account] = true
			} else {
				// Parent not shown, this is effectively a leaf
				leafAccounts[bal.Account] = true
			}
		} else {
			// Top-level account, check if it has children shown
			hasChildren := false
			for _, other := range balances {
				if strings.HasPrefix(other.Account, bal.Account+":") {
					hasChildren = true
					break
				}
			}
			if !hasChildren {
				// No children, it's a leaf
				leafAccounts[bal.Account] = true
			}
		}
	}
	
	// Calculate total from original account balances, not rollup
	accounts := c.journal.GetAccounts()
	for _, account := range accounts {
		balance := c.journal.GetLeafBalance(account)
		total.AddBalance(balance)
	}

	// Display each balance
	for _, bal := range balances {
		c.displayAccountBalance(bal)
	}

	// Display total line
	if !c.options.NoTotal {
		c.displayTotalLine(total)
	}
}

// displayAccountBalance displays a single account balance line
func (c *BalanceCommand) displayAccountBalance(bal AccountBalance) {
	balanceStr := c.formatBalance(bal.Balance)
	accountStr := bal.Account

	if c.options.Flat {
		// Flat format: no indentation
		fmt.Fprintf(os.Stdout, "%20s  %s\n", balanceStr, accountStr)
	} else {
		// Hierarchical format with indentation
		indent := strings.Repeat(" ", (bal.Depth-1)*2)
		
		// For depth 1 accounts, show the full name
		// For depth > 1 accounts, show just the leaf name
		var accountName string
		if bal.Depth == 1 {
			accountName = accountStr
		} else {
			accountName = c.getAccountLeafName(accountStr)
		}
		
		fmt.Fprintf(os.Stdout, "%20s  %s%s\n", balanceStr, indent, accountName)
	}
}

// getAccountLeafName returns the leaf name of an account
func (c *BalanceCommand) getAccountLeafName(account string) string {
	parts := strings.Split(account, ":")
	return parts[len(parts)-1]
}

// displayTotalLine displays the total line
func (c *BalanceCommand) displayTotalLine(total *domain.Balance) {
	fmt.Fprintln(os.Stdout, "--------------------")
	balanceStr := c.formatBalance(total)
	fmt.Fprintf(os.Stdout, "%20s\n", balanceStr)
}

// formatBalance formats a balance for display
func (c *BalanceCommand) formatBalance(balance *domain.Balance) string {
	if balance.IsZero() {
		return "0"
	}

	amounts := balance.GetAmounts()
	if len(amounts) == 1 {
		// Single commodity
		amount := amounts[0]
		if amount.Commodity.Symbol == "$" {
			// Default commodity - show without symbol
			return fmt.Sprintf("%.0f", amount.ToFloat64())
		}
		return amount.Format(true)
	}

	// Multiple commodities
	var parts []string
	for _, amount := range amounts {
		if amount.Commodity.Symbol == "$" {
			parts = append(parts, fmt.Sprintf("%.0f", amount.ToFloat64()))
		} else {
			parts = append(parts, amount.Format(true))
		}
	}
	return strings.Join(parts, ", ")
}