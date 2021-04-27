package spiker

import (
	"fmt"
	"unicode/utf8"
)

// Func builtin function type
type Func func(fnc *NodeFuncCallOp, scope *VariableScope) interface{}

// builtin function pool
var builtinMap = make(map[string]Func)

func init() {
	// init builtin function
	registerExport()
	registerLen()
	registerExist()
	registerDel()
	registerPrint()
}

// RegisterFunc register builtin function
func RegisterFunc(name string, fn Func) {
	if name != "" && fn != nil {
		builtinMap[name] = fn
	}
}

// return the expression value, and interrupt script
// Example: export(123)
func registerExport() {
	RegisterFunc("export", func(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
		if len(fnc.Params) != 1 {
			panic(fmt.Sprintf("export() expects 1 parameters, %d given", len(fnc.Params)))
		}

		panic(directiveExport{val: EvalExpr(fnc.Params[0], scope)})
	})
}

// return the length of expression
// Example: len("123")
func registerLen() {
	RegisterFunc("len", func(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
		if len(fnc.Params) != 1 {
			panic(fmt.Sprintf("len() expects 1 parameters, %d given", len(fnc.Params)))
		}

		val := EvalExpr(fnc.Params[0], scope)
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
	})
}

// whether a variable or index is existed
// Example: exist(var), exist(var[9]), exist(var[name])
func registerExist() {
	RegisterFunc("exist", func(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
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
			varVal := EvalExpr(v.Var, scope)
			indexVal := EvalExpr(v.Index, scope)
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
	})
}

// delete one or more variable or index
// Example: del(var), del(var["name"]), del(var[9])
func registerDel() {
	RegisterFunc("del", func(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
		if len(fnc.Params) < 1 {
			panic(fmt.Sprintf("del() expects least 1 parameters, %d given", len(fnc.Params)))
		}

		for _, node := range fnc.Params {
			switch v := node.(type) {
			case *NodeVariable:
				scope.enclosingScope.Del(v.Format())

			case *NodeVarIndex:
				varVal := EvalExpr(v.Var, scope)
				indexVal := EvalExpr(v.Index, scope)
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
	})
}

// print one or more expression value to the terminal
// Example: print(123)
func registerPrint() {
	RegisterFunc("print", func(fnc *NodeFuncCallOp, scope *VariableScope) interface{} {
		for _, v := range fnc.Params {
			fmt.Print(EvalExpr(v, scope))
		}
		return nil
	})
}
