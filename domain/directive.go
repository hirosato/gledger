package domain

import (
	"time"
)

// DirectiveType represents the type of ledger directive
type DirectiveType int

const (
	DirectiveTypeAccount DirectiveType = iota
	DirectiveTypeCommodity
	DirectiveTypePrice
	DirectiveTypeAlias
	DirectiveTypeInclude
	DirectiveTypeApply
	DirectiveBucket
	DirectiveTypeAssert
	DirectiveTypeCheck
)

// Directive represents a ledger directive
type Directive interface {
	Type() DirectiveType
	String() string
}

// AccountDirective represents an account declaration
type AccountDirective struct {
	Name string
	Note string
}

func (d *AccountDirective) Type() DirectiveType {
	return DirectiveTypeAccount
}

func (d *AccountDirective) String() string {
	if d.Note != "" {
		return "account " + d.Name + " ; " + d.Note
	}
	return "account " + d.Name
}

// CommodityDirective represents a commodity declaration
type CommodityDirective struct {
	Symbol    string
	Format    string
	Precision int
	Note      string
}

func (d *CommodityDirective) Type() DirectiveType {
	return DirectiveTypeCommodity
}

func (d *CommodityDirective) String() string {
	if d.Note != "" {
		return "commodity " + d.Symbol + " ; " + d.Note
	}
	return "commodity " + d.Symbol
}

// PriceDirective represents a price declaration (P directive)
type PriceDirective struct {
	Date      time.Time
	Commodity string
	Price     *Amount
}

func (d *PriceDirective) Type() DirectiveType {
	return DirectiveTypePrice
}

func (d *PriceDirective) String() string {
	return "P " + d.Date.Format("2006/01/02") + " " + d.Commodity + " " + d.Price.String()
}

// AliasDirective represents an alias declaration
type AliasDirective struct {
	Name  string
	Value string
}

func (d *AliasDirective) Type() DirectiveType {
	return DirectiveTypeAlias
}

func (d *AliasDirective) String() string {
	return "alias " + d.Name + "=" + d.Value
}

// IncludeDirective represents an include statement
type IncludeDirective struct {
	Path string
}

func (d *IncludeDirective) Type() DirectiveType {
	return DirectiveTypeInclude
}

func (d *IncludeDirective) String() string {
	return "include " + d.Path
}

// ApplyDirective represents an apply account directive
type ApplyDirective struct {
	Account string
}

func (d *ApplyDirective) Type() DirectiveType {
	return DirectiveTypeApply
}

func (d *ApplyDirective) String() string {
	return "apply account " + d.Account
}

// BucketDirective represents a bucket directive
type BucketDirective struct {
	Account string
}

func (d *BucketDirective) Type() DirectiveType {
	return DirectiveBucket
}

func (d *BucketDirective) String() string {
	return "bucket " + d.Account
}

// AssertDirective represents an assert directive
type AssertDirective struct {
	Expression string
}

func (d *AssertDirective) Type() DirectiveType {
	return DirectiveTypeAssert
}

func (d *AssertDirective) String() string {
	return "assert " + d.Expression
}

// CheckDirective represents a check directive
type CheckDirective struct {
	Expression string
}

func (d *CheckDirective) Type() DirectiveType {
	return DirectiveTypeCheck
}

func (d *CheckDirective) String() string {
	return "check " + d.Expression
}