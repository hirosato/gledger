package ledger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TransactionItem struct {
	Account Account
	Amount  float64
	Comment string
}

type Transaction struct {
	Date  time.Time
	Payee string
	Items []TransactionItem
}

type Line struct {
	s string
	n int
}

type Lines []Line

func (lines Lines) header() Line {
	return lines[0]
}

func (lines Lines) body() []Line {
	var result []Line
	for i := 1; i < len(lines); i++ {
		result = append(result, lines[i])
	}
	return result
}

func parseHeader(lines *Lines) (date time.Time, payee string, err error) {
	header := lines.header()
	dateAndPayee := strings.SplitN(header.s, " ", 2)
	date, err = time.Parse("2006-01-02", dateAndPayee[0])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("%d: Could not parse transaction header: %s", header.n, header.s)
	}
	return date, dateAndPayee[1], nil
}

func parseBodyLine(body Line) (account string, amount float64, comment string, err error) {
	accountAndAmount, comment, _ := strings.Cut(body.s, ";")
	account, sAmount, _ := strings.Cut(strings.TrimLeft(accountAndAmount, " "), "  ")
	amount, err = strconv.ParseFloat(strings.Trim(sAmount, " "), 64)
	if err != nil {
		return "", 0, "", fmt.Errorf("%d: Could not parse amount: %s", body.n, sAmount)
	}
	return account, amount, comment, nil
}

func parseTransaction(lines *Lines) (trn *Transaction, err error) {
	date, payee, err := parseHeader(lines)
	if err != nil {
		return nil, err
	}
	trn = &Transaction{Date: date, Payee: payee}

	body := lines.body()
	var account string
	var amount float64
	var comment string
	for i := 0; i < len(body); i++ {
		account, amount, comment, err = parseBodyLine(body[i])
		if err != nil {
			return nil, err
		}
		trn.Items = append(trn.Items, TransactionItem{Account: Account(account), Amount: amount, Comment: comment})
	}
	return trn, nil
}
