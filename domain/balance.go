package domain

import (
	"fmt"
	"sort"
	"strings"
)

type Balance struct {
	amounts map[string]*Amount
}

func NewBalance() *Balance {
	return &Balance{
		amounts: make(map[string]*Amount),
	}
}

func NewBalanceFromAmount(amount *Amount) *Balance {
	balance := NewBalance()
	balance.Add(amount)
	return balance
}

func (b *Balance) Add(amount *Amount) {
	if amount == nil || amount.IsZero() {
		return
	}
	
	commoditySymbol := amount.Commodity.Symbol
	if existing, ok := b.amounts[commoditySymbol]; ok {
		b.amounts[commoditySymbol] = existing.Add(amount)
	} else {
		b.amounts[commoditySymbol] = amount.Copy()
	}
	
	if b.amounts[commoditySymbol].IsZero() {
		delete(b.amounts, commoditySymbol)
	}
}

func (b *Balance) Subtract(amount *Amount) {
	if amount == nil || amount.IsZero() {
		return
	}
	
	b.Add(amount.Negate())
}

func (b *Balance) AddBalance(other *Balance) {
	for _, amount := range other.amounts {
		b.Add(amount)
	}
}

func (b *Balance) SubtractBalance(other *Balance) {
	for _, amount := range other.amounts {
		b.Subtract(amount)
	}
}

func (b *Balance) GetAmount(commoditySymbol string) *Amount {
	if amount, ok := b.amounts[commoditySymbol]; ok {
		return amount.Copy()
	}
	return nil
}

func (b *Balance) HasAmount(commoditySymbol string) bool {
	_, ok := b.amounts[commoditySymbol]
	return ok
}

func (b *Balance) GetAmounts() []*Amount {
	amounts := make([]*Amount, 0, len(b.amounts))
	
	commodities := make([]string, 0, len(b.amounts))
	for commodity := range b.amounts {
		commodities = append(commodities, commodity)
	}
	sort.Strings(commodities)
	
	for _, commodity := range commodities {
		amounts = append(amounts, b.amounts[commodity].Copy())
	}
	
	return amounts
}

func (b *Balance) GetCommodities() []string {
	commodities := make([]string, 0, len(b.amounts))
	for commodity := range b.amounts {
		commodities = append(commodities, commodity)
	}
	sort.Strings(commodities)
	return commodities
}

func (b *Balance) IsZero() bool {
	return len(b.amounts) == 0
}

func (b *Balance) IsEmpty() bool {
	return b.IsZero()
}

func (b *Balance) HasSingleCommodity() bool {
	return len(b.amounts) == 1
}

func (b *Balance) HasMultipleCommodities() bool {
	return len(b.amounts) > 1
}

func (b *Balance) Clear() {
	b.amounts = make(map[string]*Amount)
}

func (b *Balance) Copy() *Balance {
	copy := NewBalance()
	for commodity, amount := range b.amounts {
		copy.amounts[commodity] = amount.Copy()
	}
	return copy
}

func (b *Balance) Negate() *Balance {
	result := NewBalance()
	for commodity, amount := range b.amounts {
		result.amounts[commodity] = amount.Negate()
	}
	return result
}

func (b *Balance) Abs() *Balance {
	result := NewBalance()
	for commodity, amount := range b.amounts {
		result.amounts[commodity] = amount.Abs()
	}
	return result
}

func (b *Balance) ConvertTo(targetCommodity *Commodity, commodityRepo CommodityRepository) *Balance {
	if b.IsZero() {
		return NewBalance()
	}
	
	if b.HasSingleCommodity() {
		amount := b.GetAmounts()[0]
		if amount.Commodity.Symbol == targetCommodity.Symbol {
			return b.Copy()
		}
		
		price := amount.Commodity.GetLatestPrice(targetCommodity.Symbol)
		if price != nil {
			convertedAmount := amount.ConvertTo(targetCommodity, price.Number)
			return NewBalanceFromAmount(convertedAmount)
		}
	}
	
	result := NewBalance()
	for _, amount := range b.GetAmounts() {
		if amount.Commodity.Symbol == targetCommodity.Symbol {
			result.Add(amount)
		} else {
			price := amount.Commodity.GetLatestPrice(targetCommodity.Symbol)
			if price != nil {
				convertedAmount := amount.ConvertTo(targetCommodity, price.Number)
				result.Add(convertedAmount)
			} else {
				result.Add(amount)
			}
		}
	}
	
	return result
}

func (b *Balance) String() string {
	if b.IsZero() {
		return "0"
	}
	
	amounts := b.GetAmounts()
	if len(amounts) == 1 {
		return amounts[0].Format(true)
	}
	
	parts := make([]string, len(amounts))
	for i, amount := range amounts {
		parts[i] = amount.Format(true)
	}
	
	return strings.Join(parts, ", ")
}

func (b *Balance) Format(separator string, showZero bool) string {
	if b.IsZero() && !showZero {
		return ""
	}
	
	if b.IsZero() && showZero {
		return "0"
	}
	
	amounts := b.GetAmounts()
	parts := make([]string, 0, len(amounts))
	
	for _, amount := range amounts {
		if !amount.IsZero() || showZero {
			parts = append(parts, amount.Format(true))
		}
	}
	
	if len(parts) == 0 && showZero {
		return "0"
	}
	
	return strings.Join(parts, separator)
}

func (b *Balance) Equals(other *Balance) bool {
	if len(b.amounts) != len(other.amounts) {
		return false
	}
	
	for commodity, amount := range b.amounts {
		otherAmount, ok := other.amounts[commodity]
		if !ok || !amount.Equals(otherAmount) {
			return false
		}
	}
	
	return true
}

func (b *Balance) Validate() error {
	for commodity, amount := range b.amounts {
		if amount.Commodity.Symbol != commodity {
			return fmt.Errorf("commodity mismatch: key %s, amount commodity %s", 
				commodity, amount.Commodity.Symbol)
		}
		if amount.IsZero() {
			return fmt.Errorf("zero amount should not be stored for commodity %s", commodity)
		}
	}
	return nil
}