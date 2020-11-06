package spiker

// NodeFuncDef function define
type NodeFuncDef struct {
	Ast
	Name  NodeVariable
	Param []NodeParam
	Body  []AstNode
}

// TODO
func (fn NodeFuncDef) String() string {
	return ""
}

// NodeParam function param
type NodeParam struct {
	Ast
	DefaultValue interface{}
	Name         NodeVariable
	Value        interface{}
}

func (p NodeParam) String() string {
	f := p.Name.String()

	if p.Value != nil {
		f += Interface2String(p.Value)
	} else if p.DefaultValue != nil {
		f += Interface2String(p.DefaultValue)
	}

	return f
}
