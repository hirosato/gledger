package ledger

import (
	"bufio"
	"io"
	"strconv"
)

type Account string

type linekind int

const (
	num     = iota
	empty   = iota
	comment = iota
	space   = iota
	unknown = iota
)

func whatIsThisLine(line string) linekind {
	if len(line) == 0 {
		return empty
	}

	char := string([]rune(line)[0])
	if char == ";" {
		return comment
	}

	if char == " " {
		return space
	}

	_, err := strconv.Atoi(char)
	if err == nil {
		return num
	}
	return unknown
}

func ParseLedger(reader io.Reader) (trns []*Transaction, err error) {
	var trn *Transaction
	var line Line
	var lines Lines
	i := 0

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		i++
		line = Line{
			s: scanner.Text(),
			n: i,
		}

		switch whatIsThisLine(line.s) {
		case num: //start of a transaction
			if len(lines) > 0 {
				trn, err = parseTransaction(&lines)
				if err != nil {
					return []*Transaction{}, err
				}
				trns = append(trns, trn)
			}
			lines = []Line{line}

		case empty, comment:

		case space: // there should be account and amount in this line
			lines = append(lines, line)
		}

	}

	if len(lines) > 0 {
		trn, err = parseTransaction(&lines)
		if err != nil {
			return []*Transaction{}, err
		}
		trns = append(trns, trn)
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
