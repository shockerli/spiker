package spiker

import "fmt"

type nudFn func(*Token, *Parser) *Token

type ledFn func(*Token, *Parser, *Token) *Token

// statement function
type stdFn func(*Token, *Parser) *Token

// Token lexical token
type Token struct {
	sym          Symbol
	value        string
	line         int // Line
	col          int // Column
	bindingPower int // Priority
	nud          nudFn
	led          ledFn
	std          stdFn
	key          *Token // for NodeMap
	children     []*Token
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"on line %d:%d, symbol: %s, value: %s",
		t.line, t.col, t.sym, t.value,
	)
}
