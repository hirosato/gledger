package presenters

import (
	"fmt"
	"strings"

	"github.com/hirosato/gledger/application/dto"
)

// BalancePresenter formats balance reports for CLI output
type BalancePresenter struct {
	flat bool
}

// NewBalancePresenter creates a new balance presenter
func NewBalancePresenter(flat bool) *BalancePresenter {
	return &BalancePresenter{
		flat: flat,
	}
}

// Present formats a balance report for display
func (bp *BalancePresenter) Present(report *dto.BalanceReport) string {
	var output strings.Builder
	
	for _, account := range report.Accounts {
		if bp.flat {
			// Flat format: no indentation
			output.WriteString(fmt.Sprintf("%20s  %s\n", account.Balance, account.Name))
		} else {
			// Hierarchical format: with indentation
			indent := strings.Repeat("  ", account.Level)
			output.WriteString(fmt.Sprintf("%20s  %s%s\n", account.Balance, indent, account.Name))
		}
	}
	
	// Add total line if present
	if report.Total != nil {
		output.WriteString(strings.Repeat("-", 20))
		output.WriteString("\n")
		output.WriteString(fmt.Sprintf("%20s\n", report.Total.Balance))
	}
	
	return output.String()
}