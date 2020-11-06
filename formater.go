package spiker

import (
	"fmt"
)

const indentStep = "    "

// FormatAst format AST to source code
func FormatAst(nodeList []AstNode) (f string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	for _, node := range nodeList {
		if node == nil {
			continue
		}
		f += node.String()
		switch node.(type) {
		case *NodeIf, *NodeFuncDef, *NodeWhile:
		default:
			f += ";"
		}
		f += "\n"
	}

	return
}
