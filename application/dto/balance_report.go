package dto

// BalanceReport represents the result of a balance calculation
type BalanceReport struct {
	Accounts []AccountBalance
	Total    *AccountBalance // Optional total line
}

// AccountBalance represents a single account's balance in the report
type AccountBalance struct {
	Name     string
	Balance  string
	Level    int  // Indentation level for hierarchical display
	IsTotal  bool // Whether this is a total line
	IsEmpty  bool // Whether this account has zero balance
}