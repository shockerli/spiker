package spiker

import (
	"fmt"
	"unicode/utf8"
)

type builtinFunc func(fnc *NodeFuncCallOp, scope *VariableScope) interface{}

var builtinMap map[string]builtinFunc

func init() {
	builtinMap = map[string]builtinFunc{
		"export": export,
		"len":    length,
		"exist":  exist,
		"del":    del,
		"print":  prints,
	}
}

// return the expression value, and interrupt script
func export(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	if len(fnc.Params) != 1 {
		panic(fmt.Sprintf("export() expects 1 parameters, %d given", len(fnc.Params)))
	}

	panic(directiveExport{val: evalExpr(fnc.Params[0], scope)})
}

// return the length of expression
func length(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	if len(fnc.Params) != 1 {
		panic(fmt.Sprintf("len() expects 1 parameters, %d given", len(fnc.Params)))
	}

	val := evalExpr(fnc.Params[0], scope)
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

// whether a variable or index is existed
func exist(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	if len(fnc.Params) != 1 {
		panic(fmt.Sprintf("exist() expects 1 parameters, %d given", len(fnc.Params)))
	}
	switch fnc.Params[0].(type) {
	case *NodeVarIndex, *NodeVariable:
	default:
		panic("exist() expects parameter variable or index")
	}

	switch v := fnc.Params[0].(type) {
	case *NodeVariable:
		if _, ok := scope.Get(v.Format()); ok {
			return true
		}

	case *NodeVarIndex:
		varVal := evalExpr(v.Var, scope)
		indexVal := evalExpr(v.Index, scope)
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

// delete one or more variable or index
func del(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	if len(fnc.Params) < 1 {
		panic(fmt.Sprintf("del() expects least 1 parameters, %d given", len(fnc.Params)))
	}

	for _, node := range fnc.Params {
		switch v := node.(type) {
		case *NodeVariable:
			scope.enclosingScope.Del(v.Format())

		case *NodeVarIndex:
			varVal := evalExpr(v.Var, scope)
			indexVal := evalExpr(v.Index, scope)
			switch varVal := varVal.(type) {
			case ValueList:
				idx := int(Interface2Float64(indexVal))
				if len(varVal) <= idx {
					continue
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
					continue
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

	return nil
}

// print one or more expression value to the terminal
func prints(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
	for _, v := range fnc.Params {
		fmt.Print(evalExpr(v, scope))
	}
	return nil
}
