package spiker

import "strings"

// NodeIf if statement node
type NodeIf struct {
	Ast
	Expr   AstNode
	Body   []AstNode
	ElseIf *NodeIf
	Else   []AstNode
}

func (ifs NodeIf) String() string {
	str := ""

	if ifs.Expr != nil {
		str += "if (" + ifs.Expr.String() + ") {\n"
		for _, node := range ifs.Body {
			str += indentStep
			switch node.(type) {
			case NodeIf, NodeFuncDef, NodeWhile:
				split := strings.Split(node.String(), "\n")
				cnt := len(split)
				for idx, s := range split {
					if idx > 0 {
						str += indentStep
					}
					if len(s) > 0 {
						str += s
					}
					if idx+1 < cnt {
						str += "\n"
					}
				}
			default:
				str += node.String()
				str += ";"
			}
			str += "\n"
		}
		str += "}"
	}

	if ifs.ElseIf != nil {
		str += " else " + ifs.ElseIf.String()
	}

	if ifs.Else != nil {
		str += " else {\n"
		for _, node := range ifs.Else {
			str += indentStep + node.String() + ";\n"
		}
		str += "}"
	}

	return str
}

// NodeWhile while statement node
type NodeWhile struct {
	Ast
	Expr AstNode
	Body []AstNode
}

func (nws NodeWhile) String() string {
	str := "while (" + nws.Expr.String() + ") {\n"
	for _, node := range nws.Body {
		str += indentStep
		switch node.(type) {
		case NodeIf, NodeFuncDef, NodeWhile:
			split := strings.Split(node.String(), "\n")
			cnt := len(split)
			for idx, s := range split {
				if idx > 0 {
					str += indentStep
				}
				if len(s) > 0 {
					str += s
				}
				if idx+1 < cnt {
					str += "\n"
				}
			}
		default:
			str += node.String()
			str += ";"
		}
		str += "\n"
	}
	str += "}"

	return str
}

// NodeContinue continue node
type NodeContinue struct {
	Ast
}

func (nc NodeContinue) String() string {
	return SymbolContinue.String()
}

// NodeBreak break node
type NodeBreak struct {
	Ast
}

func (nb NodeBreak) String() string {
	return SymbolBreak.String()
}

// NodeReturn return node
type NodeReturn struct {
	Ast
}

func (nr NodeReturn) String() string {
	return SymbolReturn.String()
}
