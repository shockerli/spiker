package spiker_test

import (
	"testing"

	"github.com/shockerli/spiker"
)

func TestNodeVariable_String(t *testing.T) {
	type fields struct {
		Ast   spiker.Ast
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"name", fields{Value: "name"}, "name"},
	}
	for _, tt := range tests {
		nv := spiker.NodeVariable{
			Ast:   tt.fields.Ast,
			Value: tt.fields.Value,
		}
		if got := nv.String(); got != tt.want {
			t.Errorf("%q. NodeVariable.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeString_String(t *testing.T) {
	type fields struct {
		Ast   spiker.Ast
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"abc", fields{Value: "abc"}, `"abc"`},
		{"true", fields{Value: "true"}, `"true"`},
		{"123", fields{Value: "123"}, `"123"`},
		{"12.34", fields{Value: "12.34"}, `"12.34"`},
	}
	for _, tt := range tests {
		str := spiker.NodeString{
			Ast:   tt.fields.Ast,
			Value: tt.fields.Value,
		}
		if got := str.String(); got != tt.want {
			t.Errorf("%q. NodeString.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeNumber_String(t *testing.T) {
	type fields struct {
		Ast   spiker.Ast
		Value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"123", fields{Value: 123}, `123`},
		{"12.3", fields{Value: 12.3}, `12.3`},
		{"+12.3", fields{Value: +12.3}, `12.3`},
		{"-12.3", fields{Value: -12.3}, `-12.3`},
		{"12.30", fields{Value: 12.30}, `12.3`},
		{"0.00", fields{Value: 0.00}, `0`},
	}
	for _, tt := range tests {
		num := spiker.NodeNumber{
			Ast:   tt.fields.Ast,
			Value: tt.fields.Value,
		}
		if got := num.String(); got != tt.want {
			t.Errorf("%q. NodeNumber.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeBool_String(t *testing.T) {
	type fields struct {
		Ast   spiker.Ast
		Value bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"true", fields{Value: true}, `true`},
		{"false", fields{Value: false}, `false`},
	}
	for _, tt := range tests {
		num := spiker.NodeBool{
			Ast:   tt.fields.Ast,
			Value: tt.fields.Value,
		}
		if got := num.String(); got != tt.want {
			t.Errorf("%q. NodeBool.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeList_String(t *testing.T) {
	type fields struct {
		Ast  spiker.Ast
		List []spiker.AstNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"scalar", fields{List: []spiker.AstNode{
			spiker.NodeBool{Value: true},
			spiker.NodeNumber{Value: 123},
			spiker.NodeString{Value: "abc"},
		}}, `[true, 123, "abc"]`},
		{"mixed", fields{List: []spiker.AstNode{
			spiker.NodeVariable{Value: "name"},
			spiker.NodeVarIndex{Var: spiker.NodeVariable{Value: "item"}, Index: spiker.NodeString{Value: "age"}},
			spiker.NodeList{},
			spiker.NodeList{List: []spiker.AstNode{spiker.NodeNumber{Value: 1}, spiker.NodeNumber{Value: 2}}},
			spiker.NodeMap{},
			spiker.NodeMap{Map: map[spiker.AstNode]spiker.AstNode{
				spiker.NodeString{Value: "name"}: spiker.NodeString{Value: "judy"},
				spiker.NodeString{Value: "age"}:  spiker.NodeNumber{Value: 18},
			}},
		}}, `[name, item["age"], [], [1, 2], [], ["age": 18, "name": "judy"]]`},
	}
	for _, tt := range tests {
		arr := spiker.NodeList{
			Ast:  tt.fields.Ast,
			List: tt.fields.List,
		}
		if got := arr.String(); got != tt.want {
			t.Errorf("%q. NodeList.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeMap_String(t *testing.T) {
	type fields struct {
		Ast spiker.Ast
		Map map[spiker.AstNode]spiker.AstNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"key-string", fields{Map: map[spiker.AstNode]spiker.AstNode{
			spiker.NodeString{Value: "name"}: spiker.NodeString{Value: "judy"},
			spiker.NodeString{Value: "age"}:  spiker.NodeNumber{Value: 18},
		}}, `["age": 18, "name": "judy"]`},
		{"key-number", fields{Map: map[spiker.AstNode]spiker.AstNode{
			spiker.NodeNumber{Value: 123}:  spiker.NodeString{Value: "judy"},
			spiker.NodeNumber{Value: 12.3}: spiker.NodeNumber{Value: 18},
		}}, `[12.3: 18, 123: "judy"]`},
	}
	for _, tt := range tests {
		nm := spiker.NodeMap{
			Ast: tt.fields.Ast,
			Map: tt.fields.Map,
		}
		if got := nm.String(); got != tt.want {
			t.Errorf("%q. NodeMap.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
