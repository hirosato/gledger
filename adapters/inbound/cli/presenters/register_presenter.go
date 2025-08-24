package presenters

import (
	"fmt"
	"strings"

	"github.com/hirosato/gledger/application/dto"
)

// RegisterPresenter formats register reports for CLI output
type RegisterPresenter struct{}

// NewRegisterPresenter creates a new register presenter
func NewRegisterPresenter() *RegisterPresenter {
	return &RegisterPresenter{}
}

// Present formats a register report for display
func (rp *RegisterPresenter) Present(report *dto.RegisterReport) string {
	var output strings.Builder
	
	for _, entry := range report.Entries {
		// Format: date payee account amount running_total
		output.WriteString(fmt.Sprintf("%s %-30s %-30s %15s %15s\n",
			entry.Date,
			truncate(entry.Payee, 30),
			truncate(entry.Account, 30),
			entry.Amount,
			entry.RunningTotal,
		))
		
		if entry.Note != "" {
			output.WriteString(fmt.Sprintf("    ; %s\n", entry.Note))
		}
	}
	
	return output.String()
}

// truncate truncates a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}