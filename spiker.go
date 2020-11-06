package spiker

import "strings"

// Execute return the computed result of a string expression
func Execute(code string) (val interface{}, err error) {
	code = padSemicolon(code)

	lexer := NewLexer(code)
	p := Parser{Lexer: lexer}

	stmts, err := p.Statements()
	if err != nil {
		return
	}

	ast, err := Transform(stmts)
	if err != nil {
		return
	}

	return Evaluator(ast)
}

// ExecuteWithScope return the execute result with scope
func ExecuteWithScope(code string, scope *VariableScope) (val interface{}, err error) {
	code = padSemicolon(code)

	lexer := NewLexer(code)
	p := Parser{Lexer: lexer}

	stmts, err := p.Statements()
	if err != nil {
		return
	}

	ast, err := Transform(stmts)
	if err != nil {
		return
	}

	return EvaluateWithScope(ast, scope)
}

// Format return the formatted expression
func Format(code string) (s string, err error) {
	code = padSemicolon(code)

	lexer := NewLexer(code)
	p := Parser{Lexer: lexer}

	stmts, err := p.Statements()
	if err != nil {
		return
	}

	ast, err := Transform(stmts)
	if err != nil {
		return
	}

	return FormatAst(ast)
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
