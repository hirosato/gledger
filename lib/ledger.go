package ledger

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Account string

type TransactionItem struct {
	Account Account
	Amount  int64
}

type Transaction struct {
	Date  time.Time
	Payee string
	Items []TransactionItem
}

func ParseLedger(reader io.Reader) (trns []*Transaction, err error) {
	var trn *Transaction
	var line string
	var i = 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line = scanner.Text()
		i++
		if len(line) != 0 && trn == nil {
			dateAndPayee := strings.SplitN(line, " ", 2)
			date, err := time.Parse("2006-01-02", dateAndPayee[0])
			if err != nil {
				return nil, fmt.Errorf("%d: Could not parse date: %s", i, line)
			}
			trn = &Transaction{Date: date, Payee: dateAndPayee[1]}
		} else if len(line) != 0 && trn != nil {
			accountAndAmount := strings.Split(strings.TrimLeft(line, " "), " ")
			amount, err := strconv.ParseInt(accountAndAmount[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%d: Could not parse amount: %s", i, accountAndAmount[1])
			}
			trn.Items = append(trn.Items, TransactionItem{Account: Account(accountAndAmount[0]), Amount: amount})
		} else if len(line) == 0 && trn != nil {
			trns = append(trns, trn)
			trn = nil
		}
	}
	if trn != nil {
		return append(trns, trn), nil
	}
	return trns, nil
}

// func parseLedger(reader io.Reader) (trns []*Transaction, err error) {
// 	var trn *Transaction
// 	scanner := bufio.NewScanner(reader)
// 	var line string
// 	var lineCount int
// 	for scanner.Scan() {
// 		line = scanner.Text()
// 		lineCount++
// 		if strings.HasPrefix(line, ";") {
// 			// nop
// 		} else if len(line) == 0 {
// 			if trn != nil {
// 				transErr := balanceTransaction(trans)
// 				if transErr != nil {
// 					return trns, fmt.Errorf("%d: Unable to balance transaction, %s", lineCount, transErr)
// 				}
// 				trns = append(trns, trans)
// 				trn = nil
// 			}
// 		} else if trn == nil {
// 			lineSplit := strings.SplitN(line, " ", 2)
// 			if len(lineSplit) != 2 {
// 				return generalLedger, fmt.Errorf("%d: Unable to parse payee line: %s", lineCount, line)
// 			}
// 			dateString := lineSplit[0]
// 			transDate, dateErr := time.Parse(TransactionDateFormat, dateString)
// 			if dateErr != nil {
// 				return generalLedger, fmt.Errorf("%d: Unable to parse date: %s", lineCount, dateString)
// 			}
// 			payeeString := lineSplit[1]
// 			trans = &Transaction{Payee: payeeString, Date: transDate}
// 		} else {
// 			var accChange Account
// 			lineSplit := strings.Split(line, " ")
// 			nonEmptyWords := []string{}
// 			for _, word := range lineSplit {
// 				if len(word) > 0 {
// 					nonEmptyWords = append(nonEmptyWords, word)
// 				}
// 			}
// 			lastIndex := len(nonEmptyWords) - 1
// 			rationalNum := new(big.Rat)
// 			_, balErr := rationalNum.SetString(nonEmptyWords[lastIndex])
// 			if balErr == false {
// 				// Assuming no balance and whole line is account name
// 				accChange.Name = strings.Join(nonEmptyWords, " ")
// 			} else {
// 				accChange.Name = strings.Join(nonEmptyWords[:lastIndex], " ")
// 				accChange.Balance = rationalNum
// 			}
// 			trans.AccountChanges = append(trans.AccountChanges, accChange)
// 		}
// 	}
// 	sort.Sort(TransactionsByDate{generalLedger})
// 	return generalLedger, scanner.Err()
// }

// func ParseLedger(reader io.Reader) (generalLedger []*Transaction, err error) {
// 	parseLedger("", reader, func(t *Transaction, e error) (stop bool) {
// 		if e != nil {
// 			err = e
// 			stop = true
// 			return
// 		}

// 		generalLedger = append(generalLedger, t)
// 		return
// 	})

// 	return
// }
