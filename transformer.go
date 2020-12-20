package spiker

import (
	"fmt"
	"strconv"
)

// Transform token list to AST tree
func Transform(tokList []*Token) (nodeList []AstNode, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	for _, token := range tokList {
		nodeList = append(nodeList, transNode(token))
	}

	return
}

// Transform token to AST node
func transNode(token *Token) AstNode {
	switch token.sym {

	// Variable
	case SymbolIdent:
		return &NodeVariable{
			Ast:   Ast{raw: token},
			Value: token.value,
		}

	// Assignment
	case SymbolAssign, SymbolAssignAdd, SymbolAssignSub, SymbolAssignMul, SymbolAssignDiv, SymbolAssignMod:
		if token.children[0].sym != SymbolIdent {
			panic(fmt.Sprintf("EXPECTED IDENT, BUT GOT: %v", token.children[0].sym))
		}

		// Function declare
		if len(token.children) >= 2 && token.children[1].sym == SymbolFuncDeclare {
			return transFuncDef(token)
		}

		// Assign
		return &NodeAssignOp{
			Ast: Ast{raw: token},
			Var: NodeVariable{
				Ast:   Ast{raw: token},
				Value: token.children[0].value,
			},
			Op:   token.sym,
			Expr: transNode(token.children[1]),
		}

	// Minus(-)
	case SymbolSub:
		if len(token.children) > 1 {
			return &NodeBinaryOp{
				Ast:   Ast{raw: token},
				Left:  transNode(token.children[0]),
				Op:    token.sym,
				Right: transNode(token.children[1]),
			}
		} else if len(token.children) == 1 {
			return &NodeUnaryOp{
				Ast:   Ast{raw: token},
				Op:    token.sym,
				Right: transNode(token.children[0]),
			}
		}

	// Unary
	case SymbolLogicNot, SymbolNot:
		return &NodeUnaryOp{
			Ast:   Ast{raw: token},
			Op:    token.sym,
			Right: transNode(token.children[0]),
		}

	// Binary
	case SymbolAdd, SymbolMul, SymbolDiv, SymbolMod, SymbolPow, // +, -, *, /, %, **
		SymbolSHL, SymbolSHR, // >>, <<
		SymbolAnd, SymbolOr, SymbolXor, SymbolLogicAnd, SymbolLogicOr, // &, |, ^, &&, ||
		SymbolEQL, SymbolNEQ, SymbolGTR, SymbolGTE, SymbolLSS, SymbolLTE, // ==, !=, >, >=, <, <=
		SymbolIn: // in
		return &NodeBinaryOp{
			Ast:   Ast{raw: token},
			Left:  transNode(token.children[0]),
			Op:    token.sym,
			Right: transNode(token.children[1]),
		}

	// Number
	case SymbolNumber:
		num, _ := strconv.ParseFloat(token.value, 64)
		return &NodeNumber{
			Ast:   Ast{raw: token},
			Value: num,
		}

	// String
	case SymbolString:
		return &NodeString{
			Ast:   Ast{raw: token},
			Value: token.value,
		}

	// True
	case SymbolTrue:
		return &NodeBool{
			Ast:   Ast{raw: token},
			Value: true,
		}

	// False
	case SymbolFalse:
		return &NodeBool{
			Ast:   Ast{raw: token},
			Value: false,
		}

	// [A,B,...]
	case SymbolArray:
		arr := &NodeList{
			Ast:  Ast{raw: token},
			List: make([]AstNode, 0),
		}
		for _, subNode := range token.children {
			arr.List = append(arr.List, transNode(subNode))
		}

		return arr

	// [A:AA,BB,C:CC,...]
	case SymbolMap:
		var index float64
		dict := &NodeMap{
			Ast: Ast{raw: token},
			Map: make(map[AstNode]AstNode),
		}

		for _, subNode := range token.children {
			var key AstNode
			item := transNode(subNode)
			if subNode.key != nil {
				key = item
				item = transNode(subNode.key)
			} else {
				key = &NodeNumber{
					Ast:   Ast{raw: nil},
					Value: index,
				}

				index++
			}
			dict.Map[key] = item
		}

		return dict

	// var[i]
	case SymbolLbrack:
		idx := &NodeVarIndex{
			Ast:   Ast{raw: token},
			Var:   transNode(token.children[0]),
			Index: transNode(token.children[1]),
		}

		return idx

	// If
	case SymbolIf:
		return transIfStmt(token)

	// (
	case SymbolLparen:
		// function call
		if len(token.children) > 0 && token.children[0].sym == SymbolIdent {
			fc := &NodeFuncCallOp{
				Ast: Ast{raw: token},
				Name: NodeVariable{
					Ast:   Ast{raw: token.children[0]},
					Value: token.children[0].value,
				},
				Param: make([]AstNode, 0),
			}
			for _, pt := range token.children[1:] {
				fc.Param = append(fc.Param, transNode(pt))
			}
			return fc
		}

	// while
	case SymbolWhile:
		if len(token.children) < 2 {
			panic("Missing judgment expression for while")
		}

		nws := &NodeWhile{
			Ast:  Ast{raw: token},
			Expr: transNode(token.children[0]),
			Body: make([]AstNode, 0),
		}

		if token.children[1].sym == SymbolLbrace {
			for _, stmt := range token.children[1].children {
				nws.Body = append(nws.Body, transNode(stmt))
			}
		}

		return nws

	// continue
	case SymbolContinue:
		return &NodeContinue{
			Ast{raw: token},
		}

	// break
	case SymbolBreak:
		return &NodeBreak{
			Ast{raw: token},
		}

	}

	return nil
}

// transform token to FuncDef statement
func transFuncDef(token *Token) *NodeFuncDef {
	fnd := &NodeFuncDef{
		Ast: Ast{raw: token},
		Name: NodeVariable{
			Ast:   Ast{raw: token.children[0]},
			Value: token.children[0].value,
		},
		Param: []NodeParam{},
		Body:  nil,
	}

	tokFnd := token.children[1]
	if len(tokFnd.children) < 2 {
		panic("FUNC DECLARE EXPECTED PARAMS AND BODY")
	}

	// params
	switch tokFnd.children[0].sym {
	case SymbolIdent: // single parameter
		fnd.Param = append(fnd.Param, NodeParam{
			Ast: Ast{raw: tokFnd.children[0]},
			Name: NodeVariable{
				Ast:   Ast{raw: tokFnd.children[0]},
				Value: tokFnd.children[0].value,
			},
		})

	case SymbolTuple: // multi parameters, use tuple
		if len(tokFnd.children[0].children) > 0 {
			for _, v := range tokFnd.children[0].children {
				fnd.Param = append(fnd.Param, NodeParam{
					Ast: Ast{raw: v},
					Name: NodeVariable{
						Ast:   Ast{raw: v},
						Value: v.value,
					},
				})
			}
		}
	}

	// body
	if tokFnd.children[1].sym == SymbolLbrace {
		for _, v := range tokFnd.children[1].children {
			fnd.Body = append(fnd.Body, transNode(v))
		}
	} else {
		fnd.Body = append(fnd.Body, transNode(tokFnd.children[1]))
	}

	return fnd
}

// transform token to IF statement
func transIfStmt(token *Token) *NodeIf {
	// if
	ifStmt := &NodeIf{
		Ast:  Ast{raw: token},
		Expr: transNode(token.children[0]),
		Body: make([]AstNode, 0),
	}

	// if ... { ... <body> ... }
	if token.children[1] != nil && token.children[1].sym == SymbolLbrace {
		for _, ifBodyNode := range token.children[1].children {
			ifStmt.Body = append(ifStmt.Body, transNode(ifBodyNode))
		}
	}

	if len(token.children) < 3 {
		return ifStmt
	}

	// if ... <else if> ...
	if token.children[2] != nil && token.children[2].sym == SymbolIf {
		ifStmt.ElseIf = transIfStmt(token.children[2])
	}

	// if ... <else> ...
	if token.children[2] != nil && token.children[2].sym == SymbolLbrace {
		ifStmt.Else = make([]AstNode, 0)
		for _, elseBodyNode := range token.children[2].children {
			ifStmt.Else = append(ifStmt.Else, transNode(elseBodyNode))
		}
	}

	return ifStmt
}

// is a func call statement
func isFuncCall(token *Token) bool {
	return token.sym == SymbolLparen && len(token.children) > 0 && token.children[0].sym == SymbolIdent
}
