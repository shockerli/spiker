package spiker

// NodeFuncDef function define
type NodeFuncDef struct {
	Ast
	Name  NodeVariable
	Param []NodeParam
	Body  []AstNode
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
	Value   interface{}
}

// Format .
func (p NodeParam) Format() string {
	f := p.Name.Format()

	if p.Value != nil {
		f += Interface2String(p.Value)
	} else if p.Default != nil {
		f += Interface2String(p.Default)
	}

	return f
}
