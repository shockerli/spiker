package spiker

import (
	"fmt"
)

// Parser Lexer parser
type Parser struct {
	Lexer *Lexer
}

func (psr *Parser) expression(rbp int) *Token {
	var left *Token
	t := psr.Lexer.next()

	if t.nud != nil {
		left = t.nud(t, psr)
	} else {
		panic(fmt.Sprint("syntax error: NOT PREFIX ", t.Format()))
	}
	for rbp < psr.Lexer.peek().bindingPower {
		t := psr.Lexer.next()
		if t.led != nil {
			left = t.led(t, psr, left)
		} else {
			panic(fmt.Sprint("syntax error: NOT INFIX ", t.Format()))
		}
	}

	return left
}

// Statements parse statements
func (psr *Parser) Statements() (stmts []*Token, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	next := psr.Lexer.peek()
	for next.sym != SymbolEOF && next.sym != SymbolRbrace {
		stmts = append(stmts, psr.statement())
		next = psr.Lexer.peek()
	}

	return
}

func (psr *Parser) block() *Token {
	t := psr.Lexer.next()
	if t.sym != SymbolLbrace {
		panic(fmt.Sprint("syntax error: was looking for block start ", t.Format()))
	}
	return t.std(t, psr)
}

func (psr *Parser) statement() *Token {
	t := psr.Lexer.peek()
	if t.std != nil {
		t = psr.Lexer.next()
		return t.std(t, psr)
	}
	res := psr.expression(0)
	psr.advance(SymbolSemicolon)
	return res
}

func (psr *Parser) advance(sym Symbol) *Token {
	line := psr.Lexer.line
	col := psr.Lexer.col
	t := psr.Lexer.next()
	if t.sym != sym {
		panic(fmt.Sprintf(`syntax error: expected "%s", but got "%s", on line %d:%d`, sym, t.sym, line, col))
	}
	return t
}
