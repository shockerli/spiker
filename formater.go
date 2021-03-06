package spiker

import (
	"fmt"
	"strings"
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
		f += node.Format()
		switch node.(type) {
		case *NodeIf, *NodeFuncDef, *NodeWhile:
		default:
			f += ";"
		}
		f += "\n"
	}
	f = strings.TrimRight(f, "\n")

	return
}
