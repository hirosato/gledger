package dto

// AccountList represents a list of accounts
type AccountList struct {
	Accounts []AccountInfo
}

// AccountInfo represents information about a single account
type AccountInfo struct {
	FullName string
	Name     string
	Level    int
	Parent   string
	HasTransactions bool
}