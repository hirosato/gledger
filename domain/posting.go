package domain

import (
	"time"
)

type PostingType int

const (
	PostingTypeNormal PostingType = iota
	PostingTypeVirtual
	PostingTypeBracket
)

type BalanceAssertion struct {
	Amount    *Amount
	Date      *time.Time
	Inclusive bool
}

type CostBasis struct {
	Amount        *Amount
	Date          *time.Time
	Label         string
	PerUnitAmount *Amount
}

type Posting struct {
	Account          *Account
	Amount           *Amount
	Cost             *CostBasis
	Price            *Amount
	BalanceAssertion *BalanceAssertion
	Note             string
	Metadata         map[string]string
	Transaction      *Transaction
	Type             PostingType
	IsGenerated      bool
}

func NewPosting(account *Account) *Posting {
	return &Posting{
		Account:  account,
		Type:     PostingTypeNormal,
		Metadata: make(map[string]string),
	}
}

func (p *Posting) SetAmount(amount *Amount) {
	p.Amount = amount
}

func (p *Posting) SetCost(cost *CostBasis) {
	p.Cost = cost
}

func (p *Posting) SetPrice(price *Amount) {
	p.Price = price
}

func (p *Posting) SetBalanceAssertion(assertion *BalanceAssertion) {
	p.BalanceAssertion = assertion
}

func (p *Posting) GetMetadata(key string) (string, bool) {
	value, exists := p.Metadata[key]
	return value, exists
}

func (p *Posting) SetMetadata(key, value string) {
	p.Metadata[key] = value
}

func (p *Posting) IsVirtual() bool {
	return p.Type == PostingTypeVirtual || p.Type == PostingTypeBracket
}

func (p *Posting) IsBracketed() bool {
	return p.Type == PostingTypeBracket
}

func (p *Posting) HasCost() bool {
	return p.Cost != nil
}

func (p *Posting) HasPrice() bool {
	return p.Price != nil
}

func (p *Posting) HasBalanceAssertion() bool {
	return p.BalanceAssertion != nil
}

func (p *Posting) GetCostAmount() *Amount {
	if p.Cost == nil {
		return nil
	}
	if p.Cost.PerUnitAmount != nil && p.Amount != nil {
		return p.Cost.PerUnitAmount.Multiply(p.Amount.Number)
	}
	return p.Cost.Amount
}

func (p *Posting) GetDisplayAmount() *Amount {
	if p.Amount == nil {
		return nil
	}
	return p.Amount
}

func (p *Posting) GetMarketValue() *Amount {
	if p.HasPrice() && p.Amount != nil {
		return p.Price.Multiply(p.Amount.Number)
	}
	return p.Amount
}

func (p *Posting) Copy() *Posting {
	copy := &Posting{
		Account:     p.Account,
		Note:        p.Note,
		Type:        p.Type,
		IsGenerated: p.IsGenerated,
		Metadata:    make(map[string]string),
	}
	
	if p.Amount != nil {
		copy.Amount = p.Amount.Copy()
	}
	
	if p.Cost != nil {
		costCopy := &CostBasis{
			Label: p.Cost.Label,
		}
		if p.Cost.Amount != nil {
			costCopy.Amount = p.Cost.Amount.Copy()
		}
		if p.Cost.Date != nil {
			date := *p.Cost.Date
			costCopy.Date = &date
		}
		if p.Cost.PerUnitAmount != nil {
			costCopy.PerUnitAmount = p.Cost.PerUnitAmount.Copy()
		}
		copy.Cost = costCopy
	}
	
	if p.Price != nil {
		copy.Price = p.Price.Copy()
	}
	
	if p.BalanceAssertion != nil {
		assertionCopy := &BalanceAssertion{
			Inclusive: p.BalanceAssertion.Inclusive,
		}
		if p.BalanceAssertion.Amount != nil {
			assertionCopy.Amount = p.BalanceAssertion.Amount.Copy()
		}
		if p.BalanceAssertion.Date != nil {
			date := *p.BalanceAssertion.Date
			assertionCopy.Date = &date
		}
		copy.BalanceAssertion = assertionCopy
	}
	
	for k, v := range p.Metadata {
		copy.Metadata[k] = v
	}
	
	return copy
}