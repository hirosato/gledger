package ledger

import (
	"log"
	"sort"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/gnue/go-disp_width"
)

type Balance struct {
	ParentAccount []Account
	Account       Account
	Amount        decimal.Decimal
}

type Balances struct {
	Balances []*Balance
	Sum      decimal.Decimal
}

func (balances *Balances) PrettyPrint() {
	var accMax, amtMax, accLength, amtLength int
	for i := 0; i < len(balances.Balances); i++ {
		accLength = disp_width.Measure(string(balances.Balances[i].Account))
		amtLength = disp_width.Measure(balances.Balances[i].Amount.String())
		if accLength > accMax {
			accMax = accLength
		}
		if amtLength > amtMax {
			amtMax = amtLength
		}
	}
	for i := 0; i < len(balances.Balances); i++ {
		accLength = disp_width.Measure(string(balances.Balances[i].Account))
		amtLength = disp_width.Measure(balances.Balances[i].Amount.String())
		log.Printf("%s%s: %s%s", balances.Balances[i].Account, strings.Repeat(" ", accMax-accLength+3), strings.Repeat(" ", amtMax-amtLength+3), balances.Balances[i].Amount.String())
	}
}

func toMap(trns []*Transaction) map[string]decimal.Decimal {
	accToAmt := make(map[string]decimal.Decimal)
	var trnItm TransactionItem
	var accs []string
	var acc string
	for i := 0; i < len(trns); i++ {
		for j := 0; j < len(trns[i].Items); j++ {
			trnItm = trns[i].Items[j]
			accs = strings.Split(string(trnItm.Account), ":")
			acc = ""
			for k := 0; k < len(accs); k++ {
				// acc = strings.SplitN()[]
				if k != 0 {
					acc += ":"
				}
				acc += accs[k]
				if value, exists := accToAmt[acc]; !exists {
					accToAmt[acc] = trnItm.Amount
				} else {
					accToAmt[acc] = value.Add(trnItm.Amount)
				}
			}
		}
	}
	return accToAmt
}

func NewBalances(trns []*Transaction) (balances Balances, err error) {
	accToAmt := toMap(trns)

	balances.Balances = make([]*Balance, len(accToAmt))
	count := 0
	for acc, accBalance := range accToAmt {
		balances.Balances[count] = &Balance{
			Account: Account(acc),
			Amount:  accBalance,
		}
		count++
	}

	sort.Slice(balances.Balances, func(i, j int) bool {
		return balances.Balances[i].Account < balances.Balances[j].Account
	})
	return balances, nil
}
