package spiker

import (
	"fmt"
	"strings"
)

// NodeFuncDef function define
type NodeFuncDef struct {
	Ast
	Name       NodeVariable
	Params     []NodeParam
	Body       []AstNode
	SingleStmt bool
}

// Format .
func (fn NodeFuncDef) Format() string {
	var ps []string
	for _, v := range fn.Params {
		ps = append(ps, v.Format())
	}
	var p = strings.Join(ps, ", ")
	if len(ps) > 1 {
		p = "(" + p + ")"
	}

	var b string
	var l = len(fn.Body)
	if l == 0 {
		b = "{}"
	} else if l == 1 && fn.SingleStmt {
		b = fn.Body[0].Format()
	} else {
		var str = formatBody(fn.Body)
		b = "{\n" + str + "}"
	}

	return fmt.Sprintf("%s = %s -> %s;", fn.Name.Format(), p, b)
}

// NodeParam function param
type NodeParam struct {
	Ast
	Default interface{}
	Name    NodeVariable
}

// Format .
func (p NodeParam) Format() string {
	f := p.Name.Format()

	if p.Default != nil {
		f += " = " + Interface2String(p.Default)
	}

	return f
}
