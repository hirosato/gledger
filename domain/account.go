package domain

import (
	"strings"
)

type AccountType int

const (
	AccountTypeAsset AccountType = iota
	AccountTypeLiability
	AccountTypeEquity
	AccountTypeIncome
	AccountTypeExpense
)

type Account struct {
	Name        string
	FullName    string
	Type        AccountType
	Parent      *Account
	Children    []*Account
	Alias       string
	Note        string
	IsVirtual   bool
	IsBracketed bool
}

func NewAccount(name string) *Account {
	return &Account{
		Name:     name,
		FullName: name,
		Children: make([]*Account, 0),
	}
}

func (a *Account) AddChild(child *Account) {
	child.Parent = a
	child.updateFullName()
	a.Children = append(a.Children, child)
}

func (a *Account) updateFullName() {
	if a.Parent != nil {
		a.FullName = a.Parent.FullName + ":" + a.Name
	} else {
		a.FullName = a.Name
	}
	for _, child := range a.Children {
		child.updateFullName()
	}
}

func (a *Account) FindChild(name string) *Account {
	for _, child := range a.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (a *Account) FindOrCreateChild(name string) *Account {
	child := a.FindChild(name)
	if child == nil {
		child = NewAccount(name)
		a.AddChild(child)
	}
	return child
}

func (a *Account) IsDescendantOf(ancestor *Account) bool {
	current := a.Parent
	for current != nil {
		if current == ancestor {
			return true
		}
		current = current.Parent
	}
	return false
}

func DetermineAccountType(name string) AccountType {
	lowerName := strings.ToLower(name)
	parts := strings.Split(lowerName, ":")
	
	if len(parts) > 0 {
		firstPart := parts[0]
		switch {
		case strings.HasPrefix(firstPart, "asset"):
			return AccountTypeAsset
		case strings.HasPrefix(firstPart, "liabilit"):
			return AccountTypeLiability
		case strings.HasPrefix(firstPart, "equit"):
			return AccountTypeEquity
		case strings.HasPrefix(firstPart, "income") || strings.HasPrefix(firstPart, "revenue"):
			return AccountTypeIncome
		case strings.HasPrefix(firstPart, "expense"):
			return AccountTypeExpense
		}
	}
	
	return AccountTypeAsset
}

