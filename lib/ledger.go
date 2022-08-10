package ledger

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type linekind int

const (
	num      = iota
	empty    = iota
	comment  = iota
	indented = iota
	unknown  = iota
)

func whatIsThisLine(line string) linekind {
	if len(line) == 0 {
		return empty
	}

	char := string([]rune(line)[0])
	trim := strings.Trim(line, " ")
	if char == ";" {
		return comment
	}

	if trim == "" {
		return empty
	}

	if char == " " {
		return indented
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

		case indented: // there should be account and amount in this line
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
