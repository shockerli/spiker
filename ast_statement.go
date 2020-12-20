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

// Format .
func (ifs NodeIf) Format() string {
	str := ""

	if ifs.Expr != nil {
		str += "if (" + ifs.Expr.Format() + ") {\n"
		for _, node := range ifs.Body {
			str += indentStep
			switch node.(type) {
			case NodeIf, NodeFuncDef, NodeWhile:
				split := strings.Split(node.Format(), "\n")
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
				str += node.Format()
				str += ";"
			}
			str += "\n"
		}
		str += "}"
	}

	if ifs.ElseIf != nil {
		str += " else " + ifs.ElseIf.Format()
	}

	if ifs.Else != nil {
		str += " else {\n"
		for _, node := range ifs.Else {
			str += indentStep + node.Format() + ";\n"
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

// Format .
func (nws NodeWhile) Format() string {
	str := "while (" + nws.Expr.Format() + ") {\n"
	for _, node := range nws.Body {
		str += indentStep
		switch node.(type) {
		case NodeIf, NodeFuncDef, NodeWhile:
			split := strings.Split(node.Format(), "\n")
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
			str += node.Format()
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

// Format .
func (nc NodeContinue) Format() string {
	return SymbolContinue.String()
}

// NodeBreak break node
type NodeBreak struct {
	Ast
}

// Format .
func (nb NodeBreak) Format() string {
	return SymbolBreak.String()
}

// NodeReturn return node
type NodeReturn struct {
	Ast
	Tuples []AstNode
}

// Format .
func (nr NodeReturn) Format() string {
	var str = SymbolReturn.String()

	if len(nr.Tuples) > 0 {
		var ts []string

		for _, v := range nr.Tuples {
			ts = append(ts, v.Format())
		}

		if len(ts) > 1 {
			str += " (" + strings.Join(ts, ", ")
		} else if len(ts) == 1 {
			str += " " + ts[0]
		}
	}

	return str
}
