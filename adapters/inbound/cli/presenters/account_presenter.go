package presenters

import (
	"fmt"
	"strings"

	"github.com/hirosato/gledger/application/dto"
)

// AccountPresenter formats account lists for CLI output
type AccountPresenter struct{}

// NewAccountPresenter creates a new account presenter
func NewAccountPresenter() *AccountPresenter {
	return &AccountPresenter{}
}

// Present formats an account list for display
func (ap *AccountPresenter) Present(list *dto.AccountList) string {
	var output strings.Builder
	
	for _, account := range list.Accounts {
		output.WriteString(fmt.Sprintf("%s\n", account.FullName))
	}
	
	return output.String()
}