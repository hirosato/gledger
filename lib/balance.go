package ledger

import (
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
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
	var amtMax, amtLength int
	for i := 0; i < len(balances.Balances); i++ {
		amtLength = len(balances.Balances[i].Amount.String())
		if amtLength > amtMax {
			amtMax = amtLength
		}
	}
	for i := 0; i < len(balances.Balances); i++ {
		log.Printf("%"+strconv.Itoa(amtMax+3)+"s | %s", balances.Balances[i].Amount, balances.Balances[i].Account)
	}
	log.Println("--------------------")
	log.Printf("%"+strconv.Itoa(amtMax+3)+"s", balances.Sum.String())

}

func toMap(trns []*Transaction) (accToAmt map[string]decimal.Decimal, sum decimal.Decimal) {
	accToAmt = make(map[string]decimal.Decimal)
	var trnItm TransactionItem
	var accs []string
	var acc string
	for i := 0; i < len(trns); i++ {
		for j := 0; j < len(trns[i].Items); j++ {
			trnItm = trns[i].Items[j]
			accs = strings.Split(string(trnItm.Account), ":")
			acc = ""
			for k := 0; k < len(accs); k++ {
				if k != 0 {
					acc += ":"
				}
				acc += accs[k]
				if value, exists := accToAmt[acc]; !exists {
					accToAmt[acc] = trnItm.Amount
				} else {
					accToAmt[acc] = value.Add(trnItm.Amount)
				}
				sum = sum.Add(trnItm.Amount)
			}
		}
	}
	return accToAmt, sum
}

func NewBalances(trns []*Transaction) (balances Balances, err error) {
	accToAmt, sum := toMap(trns)

	balances.Sum = sum
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
