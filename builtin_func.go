package spiker

import (
	"fmt"
	"unicode/utf8"
)

func execBuiltinFunc(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	switch fnc.Name.Value {
	case "export":
		if len(fnc.Param) != 1 {
			panic(fmt.Sprintf("export() expects 1 parameters, %d given", len(fnc.Param)))
		}
		return export(fnc.Param[0], scope)

	case "len":
		if len(fnc.Param) != 1 {
			panic(fmt.Sprintf("len() expects 1 parameters, %d given", len(fnc.Param)))
		}
		return length(fnc.Param[0], scope)

	case "exist":
		if len(fnc.Param) != 1 {
			panic(fmt.Sprintf("exist() expects 1 parameters, %d given", len(fnc.Param)))
		}
		switch fnc.Param[0].(type) {
		case *NodeVarIndex, *NodeVariable:
		default:
			panic("exist() expects parameter variable or index")
		}
		return exist(fnc.Param[0], scope)

	case "del":
		if len(fnc.Param) < 1 {
			panic(fmt.Sprintf("del() expects least 1 parameters, %d given", len(fnc.Param)))
		}

		for _, node := range fnc.Param {
			del(node, scope)
		}

		return nil

	default:
		panic(fmt.Sprintf("call to undefined function %s()", fnc.Name.Value))
	}
}

// The export built-in function return the variable value
func export(v AstNode, scope *VariableScope) interface{} {
	return evalExpr(v, scope)
}

func isExport(node AstNode) bool {
	if node == nil || node.Raw() == nil || node.Raw().children == nil {
		return false
	}
	return isFuncCall(node.Raw()) && node.Raw().children[0].value == "export"
}

// The length built-in function return the length of v
func length(v AstNode, scope *VariableScope) int {
	val := evalExpr(v, scope)
	switch val := val.(type) {
	case int, float64:
		return -1
	case string:
		return utf8.RuneCountInString(val)
	case ValueList:
		return len(val)
	case ValueMap:
		return len(val)
	}

	return -1
}

// Whether a variable or index is existed
func exist(v AstNode, scope *VariableScope) bool {
	switch v.(type) {
	case *NodeVariable:
		if _, ok := scope.Get(v.String()); ok {
			return true
		}

	case *NodeVarIndex:
		vi := v.(*NodeVarIndex)
		varVal := evalExpr(vi.Var, scope)
		indexVal := evalExpr(vi.Index, scope)
		switch varVal := varVal.(type) {
		case string, float64, int:
			idx := int(Interface2Float64(indexVal))
			if utf8.RuneCountInString(Interface2String(varVal)) > idx {
				return true
			}

		case ValueList:
			idx := int(Interface2Float64(indexVal))
			r := varVal
			if len(r) > idx {
				return true
			}

		case ValueMap:
			idx := Interface2String(indexVal)
			r := varVal
			if _, ok := r[idx]; ok {
				return true
			}
		}

	}

	return false
}

// Delete a or more variable or index
func del(v AstNode, scope *VariableScope) {
	switch v := v.(type) {
	case *NodeVariable:
		scope.enclosingScope.Del(v.String())

	case *NodeVarIndex:
		varVal := evalExpr(v.Var, scope)
		indexVal := evalExpr(v.Index, scope)
		switch varVal := varVal.(type) {
		case ValueList:
			idx := int(Interface2Float64(indexVal))
			if len(varVal) <= idx {
				return
			}

			// delete index
			varVal = append(varVal[:idx], varVal[idx+1:]...)

			// only delete a index from variable
			switch vr := v.Var.(type) {
			case *NodeVariable:
				// self scope no value
				if _, ok := scope.enclosingScope.Get(vr.Value); ok {
					scope.enclosingScope.Set(vr.Value, varVal)
				}
			}

		case ValueMap:
			idx := Interface2String(indexVal)
			if _, ok := varVal[idx]; !ok {
				return
			}

			// delete key
			delete(varVal, idx)

			// only delete a index from variable
			switch vr := v.Var.(type) {
			case *NodeVariable:
				// self scope no value
				if _, ok := scope.enclosingScope.Get(vr.Value); ok {
					scope.enclosingScope.Set(vr.Value, varVal)
				}
			}
		}
	}
}
