package spiker

// AstNode AST node interface
type AstNode interface {
	String() string
	Raw() *Token
}

// Ast syntax parsing infrastructure
type Ast struct {
	raw *Token
}

// Raw return the token
func (ast Ast) Raw() *Token {
	return ast.raw
}
