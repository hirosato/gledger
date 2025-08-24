package usecases

import (
	"github.com/hirosato/gledger/application"
	"github.com/hirosato/gledger/application/dto"
	"github.com/hirosato/gledger/domain"
)

// ShowRegisterOptions contains options for the register display
type ShowRegisterOptions struct {
	Account   string // Filter by account
	StartDate string // Start date filter
	EndDate   string // End date filter
	Payee     string // Filter by payee
}

// ShowRegister displays transaction register
type ShowRegister struct {
	journal *application.Journal
}

// NewShowRegister creates a new ShowRegister use case
func NewShowRegister(journal *application.Journal) *ShowRegister {
	return &ShowRegister{
		journal: journal,
	}
}

// Execute returns register entries based on the options
func (sr *ShowRegister) Execute(options ShowRegisterOptions) (*dto.RegisterReport, error) {
	transactions := sr.journal.GetTransactions()
	
	report := &dto.RegisterReport{
		Entries: []dto.RegisterEntry{},
	}
	
	runningTotal := domain.NewBalance()
	
	for _, tx := range transactions {
		if sr.shouldInclude(tx, options) {
			for _, posting := range tx.Postings {
				if sr.matchesAccount(posting.Account, options.Account) {
					if posting.Amount != nil {
						runningTotal.Add(posting.Amount)
					}
					
					amountStr := ""
					if posting.Amount != nil {
						amountStr = posting.Amount.String()
					}
					
					report.Entries = append(report.Entries, dto.RegisterEntry{
						Date:         tx.Date.Format("2006/01/02"),
						Payee:        tx.Payee,
						Account:      posting.Account.FullName,
						Amount:       amountStr,
						RunningTotal: runningTotal.String(),
						Note:         posting.Note,
					})
				}
			}
		}
	}
	
	report.RunningTotal = runningTotal.String()
	return report, nil
}

func (sr *ShowRegister) shouldInclude(tx domain.Transaction, options ShowRegisterOptions) bool {
	// Implement filtering logic for date range and payee
	return true
}

func (sr *ShowRegister) matchesAccount(account *domain.Account, pattern string) bool {
	// Implement account pattern matching
	if pattern == "" {
		return true
	}
	return account.FullName == pattern
}