package domain

import (
	"fmt"
	"math/big"
	"strings"
)

type Amount struct {
	Number    *big.Rat
	Commodity *Commodity
}

func NewAmount(number *big.Rat, commodity *Commodity) *Amount {
	return &Amount{
		Number:    new(big.Rat).Set(number),
		Commodity: commodity,
	}
}

func NewAmountFromFloat(value float64, commodity *Commodity) *Amount {
	number := new(big.Rat).SetFloat64(value)
	return NewAmount(number, commodity)
}

func NewAmountFromString(s string, commodity *Commodity) (*Amount, error) {
	number := new(big.Rat)
	_, ok := number.SetString(s)
	if !ok {
		return nil, fmt.Errorf("invalid number format: %s", s)
	}
	return NewAmount(number, commodity), nil
}

func ZeroAmount(commodity *Commodity) *Amount {
	return NewAmount(big.NewRat(0, 1), commodity)
}

func (a *Amount) Add(other *Amount) *Amount {
	if a.Commodity.Symbol != other.Commodity.Symbol {
		panic(fmt.Sprintf("cannot add different commodities: %s and %s", 
			a.Commodity.Symbol, other.Commodity.Symbol))
	}
	
	result := new(big.Rat).Add(a.Number, other.Number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) Subtract(other *Amount) *Amount {
	if a.Commodity.Symbol != other.Commodity.Symbol {
		panic(fmt.Sprintf("cannot subtract different commodities: %s and %s", 
			a.Commodity.Symbol, other.Commodity.Symbol))
	}
	
	result := new(big.Rat).Sub(a.Number, other.Number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) Multiply(number *big.Rat) *Amount {
	result := new(big.Rat).Mul(a.Number, number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) MultiplyByFloat(factor float64) *Amount {
	factorRat := new(big.Rat).SetFloat64(factor)
	return a.Multiply(factorRat)
}

func (a *Amount) Divide(number *big.Rat) *Amount {
	if number.Sign() == 0 {
		panic("division by zero")
	}
	result := new(big.Rat).Quo(a.Number, number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) DivideByFloat(divisor float64) *Amount {
	if divisor == 0 {
		panic("division by zero")
	}
	divisorRat := new(big.Rat).SetFloat64(divisor)
	return a.Divide(divisorRat)
}

func (a *Amount) Negate() *Amount {
	result := new(big.Rat).Neg(a.Number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) Abs() *Amount {
	result := new(big.Rat).Abs(a.Number)
	return NewAmount(result, a.Commodity)
}

func (a *Amount) IsZero() bool {
	return a.Number.Sign() == 0
}

func (a *Amount) IsPositive() bool {
	return a.Number.Sign() > 0
}

func (a *Amount) IsNegative() bool {
	return a.Number.Sign() < 0
}

func (a *Amount) Compare(other *Amount) int {
	if a.Commodity.Symbol != other.Commodity.Symbol {
		panic(fmt.Sprintf("cannot compare different commodities: %s and %s", 
			a.Commodity.Symbol, other.Commodity.Symbol))
	}
	return a.Number.Cmp(other.Number)
}

func (a *Amount) Equals(other *Amount) bool {
	return a.Commodity.Symbol == other.Commodity.Symbol && a.Number.Cmp(other.Number) == 0
}

func (a *Amount) Copy() *Amount {
	return NewAmount(new(big.Rat).Set(a.Number), a.Commodity)
}

func (a *Amount) String() string {
	return a.Format(false)
}

func (a *Amount) Format(showCommodity bool) string {
	numberStr := a.formatNumber()
	
	if !showCommodity || a.Commodity.Symbol == "" {
		return numberStr
	}
	
	commodityStr := a.Commodity.Symbol
	if a.IsNegative() {
		return "-" + commodityStr + strings.TrimPrefix(numberStr, "-")
	}
	
	return commodityStr + numberStr
}

func (a *Amount) formatNumber() string {
	if a.Commodity.Precision == 0 {
		if a.Number.IsInt() {
			return a.Number.Num().String()
		}
		return a.Number.String()
	}
	
	floatValue, _ := a.Number.Float64()
	precision := a.Commodity.Precision
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, floatValue)
}

func (a *Amount) ConvertTo(targetCommodity *Commodity, conversionRate *big.Rat) *Amount {
	convertedNumber := new(big.Rat).Mul(a.Number, conversionRate)
	return NewAmount(convertedNumber, targetCommodity)
}

func (a *Amount) ToFloat64() float64 {
	f, _ := a.Number.Float64()
	return f
}

func (a *Amount) RoundToPrecision() *Amount {
	if a.Commodity.Precision <= 0 {
		return a.Copy()
	}
	
	precision := int64(a.Commodity.Precision)
	multiplier := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil))
	
	temp := new(big.Rat).Mul(a.Number, multiplier)
	rounded := new(big.Int)
	temp.Num().DivMod(temp.Num(), temp.Denom(), rounded)
	
	if temp.Sign() < 0 {
		if new(big.Rat).Sub(temp, new(big.Rat).SetInt(rounded)).Cmp(new(big.Rat).SetFloat64(-0.5)) <= 0 {
			rounded.Sub(rounded, big.NewInt(1))
		}
	} else {
		if new(big.Rat).Sub(temp, new(big.Rat).SetInt(rounded)).Cmp(new(big.Rat).SetFloat64(0.5)) >= 0 {
			rounded.Add(rounded, big.NewInt(1))
		}
	}
	
	result := new(big.Rat).SetInt(rounded)
	result.Quo(result, multiplier)
	
	return NewAmount(result, a.Commodity)
}