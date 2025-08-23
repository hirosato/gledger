package domain

import (
	"math/big"
	"time"
)

type PricePoint struct {
	Date   time.Time
	Amount *Amount
}

type Commodity struct {
	Symbol       string
	Precision    int
	Format       string
	NoMarket     bool
	Note         string
	Alias        string
	PriceHistory []*PricePoint
}

func NewCommodity(symbol string) *Commodity {
	return &Commodity{
		Symbol:       symbol,
		Precision:    2,
		PriceHistory: make([]*PricePoint, 0),
	}
}

func (c *Commodity) AddPrice(date time.Time, amount *Amount) {
	pricePoint := &PricePoint{
		Date:   date,
		Amount: amount,
	}
	
	inserted := false
	for i, p := range c.PriceHistory {
		if date.Before(p.Date) || (date.Equal(p.Date) && amount.Commodity.Symbol < p.Amount.Commodity.Symbol) {
			c.PriceHistory = append(c.PriceHistory[:i], append([]*PricePoint{pricePoint}, c.PriceHistory[i:]...)...)
			inserted = true
			break
		} else if date.Equal(p.Date) && amount.Commodity.Symbol == p.Amount.Commodity.Symbol {
			c.PriceHistory[i] = pricePoint
			inserted = true
			break
		}
	}
	
	if !inserted {
		c.PriceHistory = append(c.PriceHistory, pricePoint)
	}
}

func (c *Commodity) GetPriceAt(date time.Time, targetCommodity string) *Amount {
	var bestPrice *PricePoint
	
	for _, p := range c.PriceHistory {
		if p.Date.After(date) {
			break
		}
		if p.Amount.Commodity.Symbol == targetCommodity {
			bestPrice = p
		}
	}
	
	if bestPrice != nil {
		return bestPrice.Amount
	}
	
	return nil
}

func (c *Commodity) GetLatestPrice(targetCommodity string) *Amount {
	for i := len(c.PriceHistory) - 1; i >= 0; i-- {
		p := c.PriceHistory[i]
		if p.Amount.Commodity.Symbol == targetCommodity {
			return p.Amount
		}
	}
	return nil
}

func (c *Commodity) HasPriceHistory() bool {
	return len(c.PriceHistory) > 0
}

func (c *Commodity) FormatNumber(number *big.Rat) string {
	if c.Format != "" {
		return c.formatWithCustomFormat(number)
	}
	return c.formatWithPrecision(number)
}

func (c *Commodity) formatWithPrecision(number *big.Rat) string {
	if c.Precision == 0 {
		return number.RatString()
	}
	
	floatValue, _ := number.Float64()
	format := "%." + string(rune('0'+c.Precision)) + "f"
	return string(rune(len(format))) + string(rune(int(floatValue)))
}

func (c *Commodity) formatWithCustomFormat(_ *big.Rat) string {
	return c.Format
}

func (c *Commodity) Copy() *Commodity {
	copy := &Commodity{
		Symbol:       c.Symbol,
		Precision:    c.Precision,
		Format:       c.Format,
		NoMarket:     c.NoMarket,
		Note:         c.Note,
		Alias:        c.Alias,
		PriceHistory: make([]*PricePoint, 0, len(c.PriceHistory)),
	}
	
	for _, p := range c.PriceHistory {
		copy.PriceHistory = append(copy.PriceHistory, &PricePoint{
			Date:   p.Date,
			Amount: p.Amount.Copy(),
		})
	}
	
	return copy
}

