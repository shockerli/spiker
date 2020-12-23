package spiker

// Symbol lexical symbol
type Symbol string

func (sym Symbol) String() string {
	return string(sym)
}

// Supported symbol
const (
	SymbolIdent  Symbol = "(IDENT)"
	SymbolNumber Symbol = "(NUMBER)"
	SymbolString Symbol = "(STRING)"
	SymbolEOF    Symbol = "(EOF)"
	SymbolTuple  Symbol = "()"
	SymbolArray  Symbol = "[]"
	SymbolMap    Symbol = "{}"
	SymbolPound  Symbol = "#"

	SymbolTrue     Symbol = "true"
	SymbolFalse    Symbol = "false"
	SymbolNone     Symbol = "none"
	SymbolIf       Symbol = "if"
	SymbolElse     Symbol = "else"
	SymbolReturn   Symbol = "return"
	SymbolContinue Symbol = "continue"
	SymbolBreak    Symbol = "break"
	SymbolWhile    Symbol = "while"

	SymbolColon       Symbol = ":"
	SymbolSemicolon   Symbol = ";"
	SymbolLparen      Symbol = "("
	SymbolRparen      Symbol = ")"
	SymbolLbrack      Symbol = "["
	SymbolRbrack      Symbol = "]"
	SymbolLbrace      Symbol = "{"
	SymbolRbrace      Symbol = "}"
	SymbolComma       Symbol = ","
	SymbolFuncDeclare Symbol = "->"

	// mathematical
	SymbolAdd Symbol = "+"
	SymbolSub Symbol = "-"
	SymbolMul Symbol = "*"
	SymbolDiv Symbol = "/"
	SymbolMod Symbol = "%"
	SymbolPow Symbol = "**"
	SymbolIn  Symbol = "in"

	// bit arithmetic
	SymbolAnd Symbol = "&"
	SymbolOr  Symbol = "|"
	SymbolXor Symbol = "^"
	SymbolNot Symbol = "~"
	SymbolSHL Symbol = "<<"
	SymbolSHR Symbol = ">>"

	// logic
	SymbolLogicNot Symbol = "!"
	SymbolLogicAnd Symbol = "&&"
	SymbolLogicOr  Symbol = "||"

	// assignment
	SymbolAssign    Symbol = "="
	SymbolAssignAdd Symbol = "+="
	SymbolAssignSub Symbol = "-="
	SymbolAssignMul Symbol = "*="
	SymbolAssignDiv Symbol = "/="
	SymbolAssignMod Symbol = "%="

	// comparison
	SymbolEQL Symbol = "=="
	SymbolNEQ Symbol = "!="
	SymbolGTR Symbol = ">"
	SymbolGTE Symbol = ">="
	SymbolLSS Symbol = "<"
	SymbolLTE Symbol = "<="
)
