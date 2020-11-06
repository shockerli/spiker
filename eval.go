package spiker

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Evaluator run the expression and evaluate the value
func Evaluator(nodeList []AstNode) (res interface{}, err error) {
	globalScope := NewScopeTable("GLOBAL", 1, nil)

	return EvaluateWithScope(nodeList, globalScope)
}

// EvaluateWithScope same as Evaluator, evaluate with scope
func EvaluateWithScope(nodeList []AstNode, scope *VariableScope) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	for _, node := range nodeList {
		res = evalExpr(node, scope)

		// special function export(), return the specified value
		if isExport(node) {
			return
		}
	}

	// return the last statement value
	return
}

// returns the value of the expression
func evalExpr(node AstNode, scope *VariableScope) interface{} {
	switch node := node.(type) {
	case *NodeAssignOp:
		return evalAssign(node, scope)

	case *NodeUnaryOp:
		return evalUnary(node, scope)

	case *NodeBinaryOp:
		return evalBinary(node, scope)

	case *NodeVariable:
		return evalVariable(node, scope)

	case *NodeNumber:
		return node.Value

	case *NodeString:
		return node.Value

	case *NodeBool:
		return node.Value

	case *NodeList:
		return evalList(node, scope)

	case *NodeMap:
		return evalMap(node, scope)

	case *NodeVarIndex:
		return evalVarIndex(node, scope)

	case *NodeIf:
		return evalIfStmt(node, scope)

	case *NodeFuncCallOp:
		return evalFuncCall(node, scope)

	case *NodeWhile:
		return evalWhileStmt(node, scope)

	case *NodeContinue, *NodeBreak:
		scope.directPush(node)

	}

	return nil
}

// return the variable value from scope
func evalVariable(expr *NodeVariable, scope *VariableScope) interface{} {
	if val, ok := scope.Get(expr.Value); ok {
		return val
	}
	return nil
}

// init the map value
func evalMap(expr *NodeMap, scope *VariableScope) interface{} {
	dict := make(ValueMap)
	for idx, val := range expr.Map {
		dict[Interface2String(evalExpr(idx, scope))] = evalExpr(val, scope)
	}
	return dict
}

// init the list value
func evalList(expr *NodeList, scope *VariableScope) interface{} {
	list := make(ValueList, 0)
	for _, sub := range expr.List {
		list = append(list, evalExpr(sub, scope))
	}
	return list
}

// assign and return a value
func evalAssign(expr *NodeAssignOp, scope *VariableScope) interface{} {
	name := expr.Var.Value
	exprVal := evalExpr(expr.Expr, scope)
	initVal, ok := scope.Get(name) // original value
	// initial value
	if !ok {
		initVal = 0
		// string concat
		if expr.Op == SymbolAssignAdd && !IsNumber(Interface2String(exprVal)) {
			initVal = ""
		}
	}

	switch expr.Op {
	case SymbolAssign:
		scope.Set(name, exprVal)

	case SymbolAssignAdd:
		scope.Set(name, calcMath(SymbolAdd, initVal, exprVal))

	case SymbolAssignSub:
		scope.Set(name, calcMath(SymbolSub, initVal, exprVal))

	case SymbolAssignMul:
		scope.Set(name, calcMath(SymbolMul, initVal, exprVal))

	case SymbolAssignDiv:
		scope.Set(name, calcMath(SymbolDiv, initVal, exprVal))

	case SymbolAssignMod:
		scope.Set(name, calcMath(SymbolMod, initVal, exprVal))
	}

	if val, ok := scope.Get(name); ok {
		return val
	}

	return nil
}

// evalUnary unary operation
func evalUnary(expr *NodeUnaryOp, scope *VariableScope) interface{} {
	right := evalExpr(expr.Right, scope)

	switch expr.Op {
	case SymbolLogicNot:
		return !IsTrue(right)

	case SymbolNot:
		rightNumber, _ := ParseNumber(Interface2String(right))
		return ^int(rightNumber)

	case SymbolSub:
		return -Interface2Float64(right)
	}

	return nil
}

// evalBinary binary operator
func evalBinary(expr *NodeBinaryOp, scope *VariableScope) interface{} {
	left := evalExpr(expr.Left, scope)
	right := evalExpr(expr.Right, scope)

	switch expr.Op {
	case SymbolAdd, SymbolSub, SymbolMul, SymbolDiv, SymbolMod, SymbolPow,
		SymbolAnd, SymbolOr, SymbolXor, SymbolSHR, SymbolSHL:
		return calcMath(expr.Op, left, right)

	case SymbolLogicAnd:
		return IsTrue(left) && IsTrue(right)

	case SymbolLogicOr:
		return IsTrue(left) || IsTrue(right)

	case SymbolEQL, SymbolNEQ, SymbolGTR, SymbolGTE, SymbolLSS, SymbolLTE:
		return calcComparison(expr.Op, left, right)

	case SymbolIn:
		return calcIn(left, right)
	}

	return nil
}

// report whether element is within a value
func calcIn(elem interface{}, set interface{}) bool {
	leftString := Interface2String(elem)
	switch set := set.(type) {
	case ValueList:
		for _, v := range set {
			if leftString == Interface2String(v) {
				return true
			}
		}

	case ValueMap:
		for _, v := range set {
			if leftString == Interface2String(v) {
				return true
			}
		}

	case string:
		return strings.Contains(set, leftString)

	case int, float64:
		return strings.Contains(Interface2String(set), leftString)

	}

	return false
}

// mathematical calculation
func calcMath(symbol Symbol, left interface{}, right interface{}) interface{} {
	bigLeft := new(big.Float).SetFloat64(Interface2Float64(left))
	bigRight := new(big.Float).SetFloat64(Interface2Float64(right))

	leftString := Interface2String(left)
	rightString := Interface2String(right)

	leftNumber, leftErr := ParseNumber(leftString)
	rightNumber, rightErr := ParseNumber(rightString)
	isNumberExpr := leftErr == nil && rightErr == nil && IsNumber(leftString) && IsNumber(rightString)

	var bigNumber *big.Float
	switch symbol {
	case SymbolAdd:
		if !isNumberExpr {
			// concat string
			return Interface2String(left) + Interface2String(right)
		}

		// number addition
		bigNumber = new(big.Float).Add(bigLeft, bigRight)

	case SymbolSub:
		bigNumber = new(big.Float).Sub(bigLeft, bigRight)

	case SymbolMul:
		bigNumber = new(big.Float).Mul(bigLeft, bigRight)

	case SymbolDiv:
		if bigRight == new(big.Float).SetFloat64(0) {
			panic("RUNTIME ERROR: division by zero")
		}
		bigNumber = new(big.Float).Quo(bigLeft, bigRight)

	case SymbolMod:
		return int(leftNumber) % int(rightNumber)

	case SymbolPow:
		return math.Pow(leftNumber, rightNumber)

	case SymbolAnd:
		return int(leftNumber) & int(rightNumber)

	case SymbolOr:
		return int(leftNumber) | int(rightNumber)

	case SymbolXor:
		return int(leftNumber) ^ int(rightNumber)

	case SymbolSHR:
		return int(leftNumber) >> int(rightNumber)

	case SymbolSHL:
		return int(leftNumber) << int(rightNumber)
	}

	if bigNumber != nil {
		res, _ := strconv.ParseFloat(bigNumber.String(), 64)
		return res
	} else if isNumberExpr {
		return 0
	}

	return ""
}

// compare two value
func calcComparison(symbol Symbol, left interface{}, right interface{}) bool {
	leftString := Interface2String(left)
	rightString := Interface2String(right)
	leftNumber, leftErr := ParseNumber(leftString)
	rightNumber, rightErr := ParseNumber(rightString)
	isNumberExpr := leftErr == nil && rightErr == nil && IsNumber(leftString) && IsNumber(rightString)

	switch symbol {
	case SymbolEQL:
		if isNumberExpr {
			return leftNumber == rightNumber
		}
		return leftString == rightString
	case SymbolNEQ:
		if isNumberExpr {
			return leftNumber != rightNumber
		}
		return leftString != rightString
	case SymbolGTR:
		if isNumberExpr {
			return leftNumber > rightNumber
		}
		return leftString > rightString
	case SymbolGTE:
		if isNumberExpr {
			return leftNumber >= rightNumber
		}
		return leftString >= rightString
	case SymbolLSS:
		if isNumberExpr {
			return leftNumber < rightNumber
		}
		return leftString < rightString
	case SymbolLTE:
		if isNumberExpr {
			return leftNumber <= rightNumber
		}
		return leftString <= rightString
	}

	return false
}

// function call
func evalFuncCall(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	localScope := NewScopeTable("func_"+fnc.Name.Value, scope.scopeLevel+1, scope)

	// NOTICE: current only support built-in function
	res := execBuiltinFunc(fnc, localScope)

	return res
}

// return the index value
func evalVarIndex(vi *NodeVarIndex, scope *VariableScope) interface{} {
	varVal := evalExpr(vi.Var, scope)
	switch varVal := varVal.(type) {
	case string:
		idx := int(evalExpr(vi.Index, scope).(float64))
		r := []rune(varVal)
		if len(r) > idx {
			return string(r[idx])
		}
		panic(fmt.Sprintf("RUNTIME ERROR: undefined offset %d", idx))

	case float64:
		idx := int(evalExpr(vi.Index, scope).(float64))
		r := strconv.FormatFloat(varVal, 'f', -1, 64)
		if len(r) > idx {
			return r[idx]
		}
		panic(fmt.Sprintf("RUNTIME ERROR: undefined offset %d", idx))

	case int:
		idx := int(evalExpr(vi.Index, scope).(float64))
		r := strconv.Itoa(varVal)
		if len(r) > idx {
			return r[idx]
		}
		panic(fmt.Sprintf("RUNTIME ERROR: undefined offset %d", idx))

	case ValueList:
		idx := int(evalExpr(vi.Index, scope).(float64))
		r := varVal
		if len(r) > idx {
			return r[idx]
		}
		panic(fmt.Sprintf("RUNTIME ERROR: undefined offset %d", idx))

	case ValueMap:
		idx := Interface2String(evalExpr(vi.Index, scope))
		r := varVal
		if val, ok := r[idx]; ok {
			return val
		}
		panic(fmt.Sprintf("RUNTIME ERROR: undefined offset %s", idx))
	}

	return nil
}

// if-else statement
func evalIfStmt(expr *NodeIf, scope *VariableScope) interface{} {
	if expr.Expr != nil {
		ifExprVal := evalExpr(expr.Expr, scope)
		if IsTrue(ifExprVal) {
			for _, node := range expr.Body {
				if breakLevel, an := needBreak(scope); an != nil {
					if breakLevel < 0 {
						return evalExpr(node, scope)
					}
					break
				}
				r := evalExpr(node, scope)
				if isExport(node) {
					return r
				}
			}
		} else if expr.ElseIf != nil {
			return evalIfStmt(expr.ElseIf, scope)
		} else {
			for _, node := range expr.Else {
				evalExpr(node, scope)
			}
		}
	}

	return nil
}

// while statement
func evalWhileStmt(expr *NodeWhile, scope *VariableScope) interface{} {
	breakLevel := 0
	for IsTrue(evalExpr(expr.Expr, scope)) && breakLevel == 0 {
		for _, sub := range expr.Body {
			var an AstNode
			if breakLevel, an = needBreak(scope); an != nil {
				if breakLevel < 0 {
					return evalExpr(sub, scope)
				}
				break
			}
			if isExport(sub) {
				return evalExpr(sub, scope)
			}
		}
	}

	return nil
}

func needBreak(scope *VariableScope) (breakLevel int, an AstNode) {
	breakLevel = 0
	for idx, direct := range scope.directive {
		switch direct.(type) {
		case *NodeContinue:
			an = direct
			scope.directDel(idx)

		case *NodeBreak:
			an = direct
			scope.directDel(idx)
			breakLevel = 1

		default:
			if isExport(direct) {
				an = direct
				breakLevel = -1
			}
		}
	}

	return breakLevel, an
}
