package spiker

// NodeBinaryOp binary operator node
type NodeBinaryOp struct {
	Ast
	Left  AstNode
	Op    Symbol
	Right AstNode
}

// Format .
func (bin NodeBinaryOp) Format() string {
	f := " " + string(bin.Op) + " "
	switch bin.Left.(type) {
	case *NodeBinaryOp:
		f = "(" + bin.Left.Format() + ")" + f
	default:
		f = bin.Left.Format() + f
	}
	switch bin.Right.(type) {
	case *NodeBinaryOp:
		f += "(" + bin.Right.Format() + ")"
	default:
		f += bin.Right.Format()
	}

	return f
}

// NodeUnaryOp unary operator node
type NodeUnaryOp struct {
	Ast
	Op    Symbol
	Right AstNode
}

// Format .
func (un NodeUnaryOp) Format() string {
	return string(un.Op) + un.Right.Format()
}

// NodeAssignOp assignment operator node
type NodeAssignOp struct {
	Ast
	Var  NodeVariable
	Op   Symbol
	Expr AstNode
}

// Format .
func (as NodeAssignOp) Format() string {
	return as.Var.Format() + " " + string(as.Op) + " " + as.Expr.Format()
}

// NodeFuncCallOp function call node
type NodeFuncCallOp struct {
	Ast
	Name  NodeVariable
	Param []AstNode
}

// Format .
func (fu NodeFuncCallOp) Format() string {
	ps := ""
	for idx, as := range fu.Param {
		if idx > 0 {
			ps += ", "
		}
		ps += as.Format()
	}

	return fu.Name.Format() + "(" + ps + ")"
}

// NodeVarIndex return the value of the specified index(list, string)
type NodeVarIndex struct {
	Ast
	Var   AstNode
	Index AstNode
}

// Format .
func (vi NodeVarIndex) Format() string {
	return vi.Var.Format() + "[" + vi.Index.Format() + "]"
}
