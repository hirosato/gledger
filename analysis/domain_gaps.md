# Domain Model Gap Analysis

Based on analysis of complex baseline tests, here are the gaps in our current domain model:

## 1. Balance Assignment vs. Assertion Distinction

**Current State**: `BalanceAssertion` struct exists but doesn't distinguish between:
- `= amount` (assignment - sets the balance)
- `== amount` (assertion - checks the balance)

**Required Enhancement**:
```go
type BalanceAssertion struct {
    Amount      *Amount
    Date        *time.Time
    Inclusive   bool
    IsAssignment bool  // NEW: true for =, false for ==
}
```

## 2. Price Specification Types

**Current State**: Basic `Price` field in Posting
**Missing**: Distinction between:
- `@ price` (per-unit price)
- `@@ total_price` (total price)

**Required Enhancement**:
```go
type PriceSpec struct {
    Amount     *Amount
    IsTotal    bool    // true for @@, false for @
}
// Update Posting.Price to *PriceSpec
```

## 3. Lot Tracking Enhancement

**Current State**: Basic `CostBasis` struct
**Missing**: Lot notation parsing like `{{$50}}`

**Required Enhancement**: Already supported, just needs parser support.

## 4. Directive Support

**Missing Completely**: Need new domain entities:

```go
type Directive interface {
    Type() DirectiveType
}

type DirectiveType int
const (
    DirectiveTypeAccount DirectiveType = iota
    DirectiveTypeCommodity
    DirectiveTypePrice
    DirectiveTypeAlias
    DirectiveTypeInclude
)

type AccountDirective struct {
    Name string
    Note string
}

type CommodityDirective struct {
    Symbol    string
    Format    string
    Precision int
}

type PriceDirective struct {
    Date      time.Time
    Commodity string
    Price     *Amount
}
```

## 5. Journal Enhancement

**Missing**: 
- Directive storage and management
- Default commodity tracking
- Commodity registry with conversion rates

```go
type Journal struct {
    // ... existing fields ...
    Directives       []Directive
    DefaultCommodity *Commodity
    CommodityRegistry map[string]*Commodity
}
```

## 6. Parser Enhancement Requirements

**Missing Parsing Support**:
1. Balance assignment/assertion syntax (`=` vs `==`)
2. Price specification syntax (`@` vs `@@`) 
3. Lot tracking syntax (`{{amount}}`)
4. Directive parsing (`account`, `commodity`, `P`, etc.)
5. Expression amounts (`= ($4,000.00 + $100.00)`)
6. Account pattern matching in commands

## 7. Command Enhancement Requirements

**Missing Command Features**:
1. Balance report formatting options (`-n`, `-E`, `--flat`)
2. Account filtering by patterns (`:inve` matches `Assets:Investment`)
3. Multi-commodity display formatting
4. Hierarchical account tree display

## 8. Priority Order for Implementation

### Phase 3.2 (Next) - Core Command Support
1. ✅ Balance assignments/assertions parsing
2. ✅ Basic price specifications (`@`, `@@`)
3. ✅ Balance command implementation
4. ✅ Register command basic implementation

### Phase 3.3 - Advanced Features  
1. ✅ Directive parsing (`account`, `commodity`, `P`)
2. ✅ Lot tracking with cost basis
3. ✅ Expression amount parsing
4. ✅ Advanced command options

### Phase 3.4 - Polish & Compatibility
1. ✅ Full command option support
2. ✅ Multi-commodity formatting
3. ✅ Account pattern matching
4. ✅ Error handling & validation

This analysis ensures we build a robust foundation that can handle the complexity we see in the baseline tests without requiring major architectural changes later.