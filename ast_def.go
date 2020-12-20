package spiker

// NodeFuncDef function define
type NodeFuncDef struct {
	Ast
	Name   NodeVariable
	Params []NodeParam
	Body   []AstNode
}

// Format .
func (fn NodeFuncDef) Format() string {
	return ""
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
