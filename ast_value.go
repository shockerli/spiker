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

func (nv NodeVariable) String() string {
	return nv.Value
}

// NodeString string node
type NodeString struct {
	Ast
	Value string
}

func (str NodeString) String() string {
	return "\"" + str.Value + "\""
}

// NodeNumber number node
type NodeNumber struct {
	Ast
	Value float64
}

func (num NodeNumber) String() string {
	return strconv.FormatFloat(num.Value, 'f', -1, 64)
}

// NodeBool bool node
type NodeBool struct {
	Ast
	Value bool
}

func (num NodeBool) String() string {
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

func (arr NodeList) String() string {
	f := "["
	for idx, as := range arr.List {
		if idx > 0 {
			f += ", "
		}
		f += as.String()
	}
	f += "]"

	return f
}

// NodeMap map node
type NodeMap struct {
	Ast
	Map map[AstNode]AstNode
}

func (nm NodeMap) String() string {
	sortedKeys := make([]string, 0)
	sortMap := make(map[string]AstNode)
	for kn, kv := range nm.Map {
		sortedKeys = append(sortedKeys, kn.String())
		sortMap[kn.String()] = kv
	}
	sort.Strings(sortedKeys)

	f := "["
	for idx, key := range sortedKeys {
		if idx > 0 {
			f += ", "
		}
		f += key + ": " + sortMap[key].String()
	}
	f += "]"

	return f
}
