package domain

type AccountRepository interface {
	FindAccount(fullName string) *Account
	CreateAccount(fullName string) *Account
	FindOrCreateAccount(fullName string) *Account
	GetRootAccount() *Account
	GetAllAccounts() []*Account
}

type CommodityRepository interface {
	FindCommodity(symbol string) *Commodity
	CreateCommodity(symbol string) *Commodity
	FindOrCreateCommodity(symbol string) *Commodity
	RegisterCommodity(commodity *Commodity)
	GetAllCommodities() []*Commodity
	SetDefaultCommodity(commodity *Commodity)
	GetDefaultCommodity() *Commodity
}

type TransactionRepository interface {
	FindTransaction(id string) *Transaction
	SaveTransaction(transaction *Transaction) error
	GetAllTransactions() []*Transaction
	GetTransactionsByDateRange(start, end string) []*Transaction
}

type JournalRepository interface {
	LoadFromFile(filename string) error
	SaveToFile(filename string) error
	GetTransactions() []*Transaction
	GetAccounts() []*Account
	GetCommodities() []*Commodity
}