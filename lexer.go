package spiker

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Lexer lexical analyzer
type Lexer struct {
	tokReg *tokenRegistry
	source string
	index  int
	line   int
	col    int
	tok    *Token
	cached bool
}

func (lex *Lexer) nextString() *Token {
	var text bytes.Buffer
	r, size := utf8.DecodeRuneInString(lex.source[lex.index:])
	for size > 0 {
		if r == '"' {
			lex.col++
			lex.index += size
			break
		}
		if r == '\n' {
			panic(fmt.Sprint("UNTERMINATED STRING AT ", lex.line, ":", lex.col))
		}
		if r == '\\' {
			lex.col++
			lex.index += size
			r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
			if size > 0 && !unicode.IsSpace(r) {
				if r == 'r' {
					lex.consumeRune(&text, '\r', size)
				} else if r == 'n' {
					lex.consumeRune(&text, '\n', size)
				} else if r == 't' {
					lex.consumeRune(&text, '\t', size)
				} else {
					lex.consumeRune(&text, r, size)
				}
				r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
				continue
			}
		}
		lex.consumeRune(&text, r, size)
		r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
	}
	return lex.tokReg.token(SymbolString, text.String(), lex.line, lex.col)
}

func (lex *Lexer) next() *Token {
	// invalidate peekable cache
	lex.cached = false

	tmpIndex := -1
	for lex.index != tmpIndex {
		tmpIndex = lex.index
		lex.consumeWhitespace()
		lex.consumeComments()
	}

	// end of file
	if len(lex.source[lex.index:]) == 0 {
		return lex.tokReg.token(SymbolEOF, "EOF", lex.line, lex.col)
	}

	var text bytes.Buffer
	r, size := utf8.DecodeRuneInString(lex.source[lex.index:])
	for size > 0 {
		if r == '"' { // parse string
			lex.col++
			lex.index += size
			return lex.nextString()
		} else if isFirstIdentChar(r) { // parse identifiers/keywords
			col := lex.col
			lex.consumeRune(&text, r, size)
			for {
				r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
				if size > 0 && isIdentChar(r) {
					lex.consumeRune(&text, r, size)
				} else {
					break
				}
			}
			symbol := text.String()
			if lex.tokReg.defined(Symbol(symbol)) {
				return lex.tokReg.token(Symbol(symbol), symbol, lex.line, col)
			}

			return lex.tokReg.token(SymbolIdent, symbol, lex.line, col)
		} else if unicode.IsDigit(r) { // parse numbers
			col := lex.col
			lex.consumeRune(&text, r, size)
			for {
				r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
				if size > 0 && unicode.IsDigit(r) {
					lex.consumeRune(&text, r, size)
				} else {
					break
				}
			}
			if size > 0 && r == '.' {
				lex.consumeRune(&text, r, size)
				for {
					r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
					if size > 0 && unicode.IsDigit(r) {
						lex.consumeRune(&text, r, size)
					} else {
						break
					}
				}
			}
			return lex.tokReg.token(SymbolNumber, text.String(), lex.line, col)
		} else if isOperatorChar(r) { // parse operators
			col := lex.col
			lex.consumeRune(&text, r, size)

			// try to parse operators made of two characters
			var twoChar bytes.Buffer
			twoChar.WriteRune(r)
			r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
			if size > 0 && isOperatorChar(r) {
				twoChar.WriteRune(r)
				if lex.tokReg.defined(Symbol(twoChar.String())) {
					lex.consumeRune(&text, r, size)
					textStr := text.String()
					return lex.tokReg.token(Symbol(textStr), textStr, lex.line, col)
				}
			}

			// single character operator
			textStr := text.String()
			if lex.tokReg.defined(Symbol(textStr)) {
				return lex.tokReg.token(Symbol(textStr), textStr, lex.line, col)
			}
		} else {
			break
		}
	}
	panic(fmt.Sprint("INVALID CHARACTER ", lex.line, lex.col))
}

func (lex *Lexer) consumeWhitespace() {
	r, size := utf8.DecodeRuneInString(lex.source[lex.index:])
	for size > 0 && unicode.IsSpace(r) {
		if r == '\n' {
			lex.line++
			lex.col = 1
		} else {
			lex.col++
		}
		lex.index += size
		r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
	}
}

func (lex *Lexer) consumeComments() {
	r, size := utf8.DecodeRuneInString(lex.source[lex.index:])
	if Symbol(r) == SymbolComment {
		for size > 0 && r != '\n' {
			lex.col++
			lex.index += size
			r, size = utf8.DecodeRuneInString(lex.source[lex.index:])
		}
	}
}

func (lex *Lexer) consumeRune(text *bytes.Buffer, r rune, size int) {
	text.WriteRune(r)
	lex.col++
	lex.index += size
}

func (lex *Lexer) peek() *Token {
	if lex.cached {
		return lex.tok
	}
	// save current state
	index := lex.index
	line := lex.line
	col := lex.col

	// get Token and cache it
	nextToken := lex.next()
	lex.tok = nextToken
	lex.cached = true

	// restore state
	lex.index = index
	lex.line = line
	lex.col = col

	return nextToken
}

// NewLexer return Lexer instance
func NewLexer(source string) *Lexer {
	return &Lexer{tokReg: getTokenRegistry(), source: source, index: 0, line: 1, col: 1}
}

// Register token
func getTokenRegistry() *tokenRegistry {
	t := &tokenRegistry{symTable: make(map[Symbol]*Token)}
	t.symbol(SymbolIdent)  // (IDENT)
	t.symbol(SymbolNumber) // (NUMBER)
	t.symbol(SymbolString) // (STRING)

	t.symbol(SymbolTrue)  // true
	t.symbol(SymbolFalse) // false
	t.symbol(SymbolNone)  // none

	t.consumable(SymbolColon)     // :
	t.consumable(SymbolSemicolon) // ;
	t.consumable(SymbolRparen)    // )
	t.consumable(SymbolRbrack)    // ]
	t.consumable(SymbolComma)     // ,
	t.consumable(SymbolElse)      // else

	t.consumable(SymbolEOF)    // (EOF)
	t.consumable(SymbolLbrace) // {
	t.consumable(SymbolRbrace) // }

	t.infix(SymbolIn, 70)  // in
	t.infix(SymbolPow, 68) // **

	t.infix(SymbolAdd, 60) // +
	t.infix(SymbolSub, 60) // -

	t.infix(SymbolMul, 65) // *
	t.infix(SymbolDiv, 65) // /
	t.infix(SymbolMod, 65) // %

	t.infix(SymbolSHL, 50) // <<
	t.infix(SymbolSHR, 50) // >>

	t.infix(SymbolLSS, 40) // <
	t.infix(SymbolGTR, 40) // >
	t.infix(SymbolLTE, 40) // <=
	t.infix(SymbolGTE, 40) // >=
	t.infix(SymbolEQL, 35) // ==
	t.infix(SymbolNEQ, 35) // !=

	t.infix(SymbolAnd, 33) // &
	t.infix(SymbolXor, 32) // ^
	t.infix(SymbolOr, 31)  // |

	t.infix(SymbolLogicAnd, 25) // &&
	t.infix(SymbolLogicOr, 25)  // ||

	// (
	t.infixLed(SymbolLparen, 90, func(token *Token, p *Parser, left *Token) *Token {
		if left.sym != SymbolIdent && left.sym != SymbolLbrack && left.sym != SymbolLparen && left.sym != SymbolFuncDeclare {
			panic(fmt.Sprint("BAD FUNC CALL LEFT OPERAND: ", left))
		}
		token.children = append(token.children, left)
		t := p.Lexer.peek()
		if t.sym != SymbolRparen {
			for {
				exp := p.expression(0)
				token.children = append(token.children, exp)
				if p.Lexer.peek().sym != SymbolComma {
					break
				}
				p.advance(SymbolComma)
			}
			p.advance(SymbolRparen)
		} else {
			p.advance(SymbolRparen)
		}
		return token
	})

	// [
	t.infixLed(SymbolLbrack, 80, func(token *Token, p *Parser, left *Token) *Token {
		if left.sym != SymbolIdent && left.sym != SymbolLbrack && left.sym != SymbolLparen {
			panic(fmt.Sprint("BAD ARRAY LEFT OPERAND: ", left))
		}
		token.children = append(token.children, left)
		t := p.Lexer.peek()
		if t.sym != SymbolRbrack {
			for {
				exp := p.expression(0)
				token.children = append(token.children, exp)
				if p.Lexer.peek().sym != SymbolComma {
					break
				}
				p.advance(SymbolComma)
			}
			p.advance(SymbolRbrack)
		} else {
			p.advance(SymbolRbrack)
		}
		return token
	})

	t.infixRight(SymbolAssign, 10)    // =
	t.infixRight(SymbolAssignAdd, 10) // +=
	t.infixRight(SymbolAssignSub, 10) // -=
	t.infixRight(SymbolAssignMul, 10) // *=
	t.infixRight(SymbolAssignDiv, 10) // /=
	t.infixRight(SymbolAssignMod, 10) // %=

	// ->
	t.infixRightLed(SymbolFuncDeclare, 10, func(token *Token, p *Parser, left *Token) *Token {
		if left.sym != SymbolTuple && left.sym != SymbolIdent {
			panic(fmt.Sprint("INVALID FUNC DECLARATION TUPLE: ", left))
		}
		if left.sym == SymbolTuple && len(left.children) != 0 {
			named := true
			for _, child := range left.children {
				if child.sym != SymbolIdent {
					named = false
					break
				}
			}
			if !named {
				panic(fmt.Sprint("INVALID FUNC DECLARATION TUPLE: ", left))
			}
		}
		token.children = append(token.children, left)
		if p.Lexer.peek().sym == SymbolLbrace {
			token.children = append(token.children, p.block())
		} else {
			token.children = append(token.children, p.expression(0))
		}
		return token
	})

	t.prefix(SymbolSub)      // -
	t.prefix(SymbolLogicNot) // !
	t.prefix(SymbolNot)      // ~

	// (
	t.prefixNud(SymbolLparen, func(t *Token, p *Parser) *Token {
		comma := false
		if p.Lexer.peek().sym != SymbolRparen {
			for {
				if p.Lexer.peek().sym == SymbolRparen {
					break
				}
				t.children = append(t.children, p.expression(0))
				if p.Lexer.peek().sym != SymbolComma {
					break
				}
				comma = true
				p.advance(SymbolComma)
			}
		}
		p.advance(SymbolRparen)
		if len(t.children) == 0 || comma {
			t.sym = SymbolTuple
			t.value = "TUPLE"
			return t
		}
		return t.children[0]
	})

	// ARRAY[]
	t.prefixNud(SymbolLbrack, func(t *Token, p *Parser) *Token {
		isArray := true
		children := make([]*Token, 0)
		if p.Lexer.peek().sym != SymbolRbrack {
			for {
				// ]
				if p.Lexer.peek().sym == SymbolRbrack {
					break
				}

				item := p.expression(0)

				// :
				if p.Lexer.peek().sym == SymbolColon {
					p.advance(SymbolColon)
					isArray = false

					item.key = p.expression(0)
				}
				children = append(children, item)

				// ;
				if p.Lexer.peek().sym != SymbolComma {
					break
				}
				p.advance(SymbolComma)
			}
		}
		p.advance(SymbolRbrack)
		t.children = children
		if isArray {
			t.sym = SymbolArray
			t.value = "ARRAY"
		} else {
			t.sym = SymbolMap
			t.value = "MAP"
		}
		return t
	})

	// if
	t.stmt(SymbolIf, func(t *Token, p *Parser) *Token {
		t.children = append(t.children, p.expression(0))
		t.children = append(t.children, p.block())
		next := p.Lexer.peek()
		if next.value == "else" {
			p.Lexer.next()
			next = p.Lexer.peek()
			if next.value == "if" {
				t.children = append(t.children, p.statement())
			} else {
				t.children = append(t.children, p.block())
			}
		}
		return t
	})

	// while
	t.stmt(SymbolWhile, func(t *Token, p *Parser) *Token {
		t.children = append(t.children, p.expression(0))
		t.children = append(t.children, p.block())
		return t
	})

	// {
	t.stmt(SymbolLbrace, func(t *Token, p *Parser) *Token {
		stmts, err := p.Statements()
		if err != nil {
			panic(err.Error())
		}
		t.children = append(t.children, stmts...)
		p.advance(SymbolRbrace)
		return t
	})

	// break
	t.stmt(SymbolBreak, func(t *Token, p *Parser) *Token {
		p.advance(SymbolSemicolon)
		return t
	})

	// continue
	t.stmt(SymbolContinue, func(t *Token, p *Parser) *Token {
		p.advance(SymbolSemicolon)
		return t
	})

	// return
	t.stmt(SymbolReturn, func(t *Token, p *Parser) *Token {
		if p.Lexer.peek().sym != SymbolSemicolon {
			t.children = append(t.children, p.expression(0))
		}
		p.advance(SymbolSemicolon)
		return t
	})

	return t
}

// Is first ident char
func isFirstIdentChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r == '_')
}

// Is ident char
func isIdentChar(r rune) bool {
	return isFirstIdentChar(r) || unicode.IsDigit(r)
}

// Is operator char
func isOperatorChar(r rune) bool {
	operators := "~&!#%^*()-+=/?.,:;\"|/{}[]><"
	for _, c := range operators {
		if c == r {
			return true
		}
	}
	return false
}
