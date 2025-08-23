package domain

import (
	"time"
)

type TransactionStatus int

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusCleared
	TransactionStatusReconciled
)

type Transaction struct {
	ID       string
	Date     time.Time
	AuxDate  *time.Time
	Status   TransactionStatus
	Code     string
	Payee    string
	Note     string
	Postings []*Posting
	Metadata map[string]string
}

func NewTransaction(date time.Time) *Transaction {
	return &Transaction{
		Date:     date,
		Status:   TransactionStatusPending,
		Postings: make([]*Posting, 0),
		Metadata: make(map[string]string),
	}
}

func (t *Transaction) AddPosting(posting *Posting) {
	posting.Transaction = t
	t.Postings = append(t.Postings, posting)
}

func (t *Transaction) IsBalanced() bool {
	if len(t.Postings) < 2 {
		return false
	}
	
	balances := make(map[string]*Amount)
	
	for _, posting := range t.Postings {
		if posting.Amount != nil {
			commodity := posting.Amount.Commodity.Symbol
			if existing, ok := balances[commodity]; ok {
				balances[commodity] = existing.Add(posting.Amount)
			} else {
				balances[commodity] = posting.Amount.Copy()
			}
		}
	}
	
	for _, balance := range balances {
		if !balance.IsZero() {
			return false
		}
	}
	
	return true
}

func (t *Transaction) GetMetadata(key string) (string, bool) {
	value, exists := t.Metadata[key]
	return value, exists
}

func (t *Transaction) SetMetadata(key, value string) {
	t.Metadata[key] = value
}

func (t *Transaction) Copy() *Transaction {
	copy := &Transaction{
		ID:       t.ID,
		Date:     t.Date,
		Status:   t.Status,
		Code:     t.Code,
		Payee:    t.Payee,
		Note:     t.Note,
		Postings: make([]*Posting, 0, len(t.Postings)),
		Metadata: make(map[string]string),
	}
	
	if t.AuxDate != nil {
		auxDate := *t.AuxDate
		copy.AuxDate = &auxDate
	}
	
	for _, posting := range t.Postings {
		copy.AddPosting(posting.Copy())
	}
	
	for k, v := range t.Metadata {
		copy.Metadata[k] = v
	}
	
	return copy
}

func (t *Transaction) String() string {
	return t.Date.Format("2006/01/02") + " " + t.Payee
}