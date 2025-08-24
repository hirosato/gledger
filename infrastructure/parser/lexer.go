package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// TokenType represents the type of a lexical token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenWhitespace
	TokenDate
	TokenStatus       // *, !, ?
	TokenCode         // (code)
	TokenDescription  // Transaction description
	TokenAccount      // Account name
	TokenAmount       // Numeric amount
	TokenCommodity    // Currency/commodity symbol
	TokenComment      // ; comment
	TokenIndent       // Leading whitespace for postings
	TokenEqual        // = for balance assertions
	TokenDoubleEqual  // == for balance checks
	TokenAt           // @ for prices
	TokenDoubleAt     // @@ for lot prices
	TokenSemicolon    // ;
	TokenColon        // :
	TokenMinus        // -
	TokenPlus         // +
	TokenNumber       // Numeric value
	TokenString       // General string
)

// Token represents a lexical token
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Lexer tokenizes ledger input
type Lexer struct {
	reader  *bufio.Reader
	line    int
	column  int
	current rune
	peek    rune
	atEOF   bool
}

// NewLexer creates a new lexer for the given input
func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{
		reader: bufio.NewReader(input),
		line:   1,
		column: 0,
	}
	l.advance() // Read first character
	l.advance() // Set up peek
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	// Skip whitespace except newlines
	for unicode.IsSpace(l.current) && l.current != '\n' {
		l.advance()
	}

	// Check for EOF
	if l.atEOF {
		return Token{Type: TokenEOF, Line: l.line, Column: l.column}
	}

	// Save position for token
	line := l.line
	column := l.column

	// Handle newlines
	if l.current == '\n' {
		l.advance()
		l.line++
		l.column = 0
		return Token{Type: TokenNewline, Value: "\n", Line: line, Column: column}
	}

	// Handle comments
	if l.current == ';' {
		return l.readComment()
	}

	// Handle dates (YYYY-MM-DD or YYYY/MM/DD)
	if unicode.IsDigit(l.current) && l.isDateStart() {
		return l.readDate()
	}

	// Handle amounts (numbers with optional decimal)
	if unicode.IsDigit(l.current) || l.current == '-' {
		return l.readAmount()
	}

	// Handle status markers
	if l.current == '*' || l.current == '!' || l.current == '?' {
		status := string(l.current)
		l.advance()
		return Token{Type: TokenStatus, Value: status, Line: line, Column: column}
	}

	// Handle equals (balance assertions)
	if l.current == '=' {
		l.advance()
		if l.current == '=' {
			l.advance()
			return Token{Type: TokenDoubleEqual, Value: "==", Line: line, Column: column}
		}
		return Token{Type: TokenEqual, Value: "=", Line: line, Column: column}
	}

	// Handle @ (price specifications)
	if l.current == '@' {
		l.advance()
		if l.current == '@' {
			l.advance()
			return Token{Type: TokenDoubleAt, Value: "@@", Line: line, Column: column}
		}
		return Token{Type: TokenAt, Value: "@", Line: line, Column: column}
	}

	// Handle account names and other text
	return l.readText()
}

// advance reads the next character
func (l *Lexer) advance() {
	if l.atEOF {
		return
	}

	l.current = l.peek
	
	r, _, err := l.reader.ReadRune()
	if err != nil {
		l.peek = 0
		if l.current == 0 {
			l.atEOF = true
		}
	} else {
		l.peek = r
	}
	
	l.column++
}

// isDateStart checks if current position might be start of a date
func (l *Lexer) isDateStart() bool {
	// Simple check: 4 digits followed by - or /
	// This is a simplified version; real implementation would be more robust
	return unicode.IsDigit(l.current)
}

// readDate reads a date token
func (l *Lexer) readDate() Token {
	line := l.line
	column := l.column
	var date strings.Builder

	// Read date in format YYYY-MM-DD or YYYY/MM/DD
	for unicode.IsDigit(l.current) || l.current == '-' || l.current == '/' {
		date.WriteRune(l.current)
		l.advance()
	}

	dateStr := date.String()
	// Simple validation: should have format XXXX-XX-XX or XXXX/XX/XX
	if len(dateStr) == 10 && (dateStr[4] == '-' || dateStr[4] == '/') && 
	   (dateStr[7] == '-' || dateStr[7] == '/') {
		return Token{Type: TokenDate, Value: dateStr, Line: line, Column: column}
	}

	// Not a valid date, treat as string
	return Token{Type: TokenString, Value: dateStr, Line: line, Column: column}
}

// readAmount reads a numeric amount
func (l *Lexer) readAmount() Token {
	line := l.line
	column := l.column
	var amount strings.Builder

	// Handle negative sign
	if l.current == '-' {
		amount.WriteRune(l.current)
		l.advance()
	}

	// Read digits before decimal
	for unicode.IsDigit(l.current) || l.current == ',' {
		if l.current != ',' { // Skip thousand separators
			amount.WriteRune(l.current)
		}
		l.advance()
	}

	// Handle decimal point
	if l.current == '.' {
		amount.WriteRune(l.current)
		l.advance()
		
		// Read digits after decimal
		for unicode.IsDigit(l.current) {
			amount.WriteRune(l.current)
			l.advance()
		}
	}

	return Token{Type: TokenAmount, Value: amount.String(), Line: line, Column: column}
}

// readComment reads a comment until end of line
func (l *Lexer) readComment() Token {
	line := l.line
	column := l.column
	var comment strings.Builder

	// Skip the semicolon
	l.advance()

	// Read until newline
	for l.current != '\n' && !l.atEOF {
		comment.WriteRune(l.current)
		l.advance()
	}

	return Token{Type: TokenComment, Value: strings.TrimSpace(comment.String()), Line: line, Column: column}
}

// readText reads account names, descriptions, and other text
func (l *Lexer) readText() Token {
	line := l.line
	column := l.column
	var text strings.Builder

	// Read until we hit a delimiter
	for !l.atEOF && l.current != '\n' && l.current != ';' && 
	    l.current != '=' && l.current != '@' {
		text.WriteRune(l.current)
		l.advance()
		
		// Stop at double space (often separates account from amount)
		if l.current == ' ' && l.peek == ' ' {
			break
		}
	}

	value := strings.TrimSpace(text.String())
	
	// Determine token type based on content
	if strings.Contains(value, ":") {
		return Token{Type: TokenAccount, Value: value, Line: line, Column: column}
	}

	return Token{Type: TokenString, Value: value, Line: line, Column: column}
}

// PeekToken returns the next token without consuming it
func (l *Lexer) PeekToken() Token {
	// This is a simplified implementation
	// A real implementation would maintain a buffer
	panic("PeekToken not implemented")
}

// TokenTypeString returns a string representation of the token type
func TokenTypeString(t TokenType) string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "Newline"
	case TokenWhitespace:
		return "Whitespace"
	case TokenDate:
		return "Date"
	case TokenStatus:
		return "Status"
	case TokenCode:
		return "Code"
	case TokenDescription:
		return "Description"
	case TokenAccount:
		return "Account"
	case TokenAmount:
		return "Amount"
	case TokenCommodity:
		return "Commodity"
	case TokenComment:
		return "Comment"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}