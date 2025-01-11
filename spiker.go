package spiker

import (
	"crypto/sha1"
	"strings"
	"sync"
)

// EnableAstCache enable ast cache, default is false
var EnableAstCache = false

// cachedAst cached ast node map
var cachedAst = sync.Map{}

// Execute return the computed result of a string expression
func Execute(code string) (val interface{}, err error) {
	ast, err := ParseAst(code)
	if err != nil {
		return
	}

	return Evaluator(ast)
}

// ExecuteWithScope return the execute result with scope
func ExecuteWithScope(code string, scope *VariableScope) (val interface{}, err error) {
	ast, err := ParseAst(code)
	if err != nil {
		return
	}

	return EvaluateWithScope(ast, scope)
}

// Format return the formatted expression
func Format(code string) (s string, err error) {
	ast, err := ParseAst(code)
	if err != nil {
		return
	}

	return FormatAst(ast)
}

// ParseAst lexer, statements, transform, and return the ast nodes
// if EnableAstCache is true, it will cache the ast nodes
func ParseAst(code string) (ast []AstNode, err error) {
	// get cached ast nodes
	hashKey := sha1.Sum([]byte(code))
	if EnableAstCache {
		if ast, ok := cachedAst.Load(hashKey); ok {
			return ast.([]AstNode), nil
		}
	}

	// padding semicolon
	code = padSemicolon(code)

	// lexer, parser, transform
	lexer := NewLexer(code)
	p := Parser{Lexer: lexer}

	// parse statements
	stmts, err := p.Statements()
	if err != nil {
		return
	}

	// transform to ast nodes
	ast, err = Transform(stmts)
	if err != nil {
		return
	}

	// cache ast nodes
	if EnableAstCache {
		cachedAst.Store(hashKey, ast)
	}

	return
}

// padding semicolon
func padSemicolon(code string) string {
	code = strings.TrimSpace(code)
	last := code[len(code)-1:]
	if last != SymbolSemicolon.String() && last != SymbolRbrace.String() {
		code += SymbolSemicolon.String()
	}
	return code
}
