package spiker

import (
	"sort"
	"strconv"
)

// ValueList value list
type ValueList []interface{}

// ValueMap kv dict
type ValueMap map[string]interface{}

// NodeVariable variable node
type NodeVariable struct {
	Ast
	Value string
}

// Format .
func (nv NodeVariable) Format() string {
	return nv.Value
}

// NodeString string node
type NodeString struct {
	Ast
	Value string
}

// Format .
func (str NodeString) Format() string {
	return "\"" + str.Value + "\""
}

// NodeNumber number node
type NodeNumber struct {
	Ast
	Value float64
}

// Format .
func (num NodeNumber) Format() string {
	return strconv.FormatFloat(num.Value, 'f', -1, 64)
}

// NodeBool bool node
type NodeBool struct {
	Ast
	Value bool
}

// Format .
func (num NodeBool) Format() string {
	if num.Value {
		return "true"
	}
	return "false"
}

// NodeList list node
type NodeList struct {
	Ast
	List []AstNode
}

// Format .
func (arr NodeList) Format() string {
	f := "["
	for idx, as := range arr.List {
		if idx > 0 {
			f += ", "
		}
		f += as.Format()
	}
	f += "]"

	return f
}

// NodeMap map node
type NodeMap struct {
	Ast
	Map map[AstNode]AstNode
}

// Format .
func (nm NodeMap) Format() string {
	sortedKeys := make([]string, 0)
	sortMap := make(map[string]AstNode)
	for kn, kv := range nm.Map {
		sortedKeys = append(sortedKeys, kn.Format())
		sortMap[kn.Format()] = kv
	}
	sort.Strings(sortedKeys)

	f := "["
	for idx, key := range sortedKeys {
		if idx > 0 {
			f += ", "
		}
		f += key + ": " + sortMap[key].Format()
	}
	f += "]"

	return f
}
