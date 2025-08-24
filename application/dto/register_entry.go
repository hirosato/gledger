package dto

// RegisterReport represents the result of a register query
type RegisterReport struct {
	Entries      []RegisterEntry
	RunningTotal string
}

// RegisterEntry represents a single entry in the register
type RegisterEntry struct {
	Date         string
	Payee        string
	Account      string
	Amount       string
	RunningTotal string
	Note         string
}