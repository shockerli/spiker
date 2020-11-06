package spiker

// NodeBinaryOp binary operator node
type NodeBinaryOp struct {
	Ast
	Left  AstNode
	Op    Symbol
	Right AstNode
}

func (bin NodeBinaryOp) String() string {
	f := " " + string(bin.Op) + " "
	switch bin.Left.(type) {
	case *NodeBinaryOp:
		f = "(" + bin.Left.String() + ")" + f
	default:
		f = bin.Left.String() + f
	}
	switch bin.Right.(type) {
	case *NodeBinaryOp:
		f += "(" + bin.Right.String() + ")"
	default:
		f += bin.Right.String()
	}

	return f
}

// NodeUnaryOp unary operator node
type NodeUnaryOp struct {
	Ast
	Op    Symbol
	Right AstNode
}

func (un NodeUnaryOp) String() string {
	return string(un.Op) + un.Right.String()
}

// NodeAssignOp assignment operator node
type NodeAssignOp struct {
	Ast
	Var  NodeVariable
	Op   Symbol
	Expr AstNode
}

func (as NodeAssignOp) String() string {
	return as.Var.String() + " " + string(as.Op) + " " + as.Expr.String()
}

// NodeFuncCallOp function call node
type NodeFuncCallOp struct {
	Ast
	Name  NodeVariable
	Param []AstNode
}

func (fu NodeFuncCallOp) String() string {
	ps := ""
	for idx, as := range fu.Param {
		if idx > 0 {
			ps += ", "
		}
		ps += as.String()
	}

	return fu.Name.String() + "(" + ps + ")"
}

// NodeVarIndex return the value of the specified index(list, string)
type NodeVarIndex struct {
	Ast
	Var   AstNode
	Index AstNode
}

func (vi NodeVarIndex) String() string {
	return vi.Var.String() + "[" + vi.Index.String() + "]"
}
