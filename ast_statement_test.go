package spiker_test

import (
	"testing"

	"github.com/shockerli/spiker"
)

func TestNodeIf_String(t *testing.T) {
	type fields struct {
		Ast    spiker.Ast
		Expr   spiker.AstNode
		Body   []spiker.AstNode
		ElseIf *spiker.NodeIf
		Else   []spiker.AstNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"only-if", fields{Expr: spiker.NodeBool{Value: true}}, `if (true) {
}`},
		{"with-body", fields{Expr: spiker.NodeNumber{Value: 123}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "judy"}},
			spiker.NodeIf{Expr: spiker.NodeString{Value: "abc"}, Body: []spiker.AstNode{
				spiker.NodeBinaryOp{Left: spiker.NodeVariable{Value: "count"}, Op: spiker.SymbolAssignAdd, Right: spiker.NodeNumber{Value: 1}},
			}},
		}}, `if (123) {
    name = "judy";
    if ("abc") {
        count += 1;
    }
}`},
		{"with-else-if", fields{Expr: spiker.NodeBinaryOp{
			Left: spiker.NodeVariable{Value: "age"}, Op: spiker.SymbolGTR, Right: spiker.NodeNumber{Value: 18},
		}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "judy"}},
		}, ElseIf: &spiker.NodeIf{Expr: spiker.NodeBinaryOp{
			Left: spiker.NodeVariable{Value: "age"}, Op: spiker.SymbolLTE, Right: spiker.NodeNumber{Value: 18},
		}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "tom"}},
		}}}, `if (age > 18) {
    name = "judy";
} else if (age <= 18) {
    name = "tom";
}`},
		{"with-else", fields{Expr: spiker.NodeBinaryOp{
			Left: spiker.NodeVariable{Value: "age"}, Op: spiker.SymbolGTR, Right: spiker.NodeNumber{Value: 18},
		}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "judy"}},
		}, Else: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "tom"}},
		}}, `if (age > 18) {
    name = "judy";
} else {
    name = "tom";
}`},
	}
	for _, tt := range tests {
		ifs := spiker.NodeIf{
			Ast:    tt.fields.Ast,
			Expr:   tt.fields.Expr,
			Body:   tt.fields.Body,
			ElseIf: tt.fields.ElseIf,
			Else:   tt.fields.Else,
		}
		if got := ifs.String(); got != tt.want {
			t.Errorf("%q. NodeIf.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeWhile_String(t *testing.T) {
	type fields struct {
		Ast  spiker.Ast
		Expr spiker.AstNode
		Body []spiker.AstNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"while", fields{Expr: spiker.NodeBool{Value: true}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "judy"}},
		}}, `while (true) {
    name = "judy";
}`},
		{"with-if", fields{Expr: spiker.NodeBool{Value: true}, Body: []spiker.AstNode{
			spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: spiker.NodeString{Value: "judy"}},
			spiker.NodeIf{Expr: spiker.NodeString{Value: "abc"}, Body: []spiker.AstNode{
				spiker.NodeBinaryOp{Left: spiker.NodeVariable{Value: "count"}, Op: spiker.SymbolAssignAdd, Right: spiker.NodeNumber{Value: 1}},
			}},
		}}, `while (true) {
    name = "judy";
    if ("abc") {
        count += 1;
    }
}`},
	}
	for _, tt := range tests {
		nws := spiker.NodeWhile{
			Ast:  tt.fields.Ast,
			Expr: tt.fields.Expr,
			Body: tt.fields.Body,
		}
		if got := nws.String(); got != tt.want {
			println(got)
			t.Errorf("%q. NodeWhile.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeContinue_String(t *testing.T) {
	type fields struct {
		Ast spiker.Ast
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"return", fields{}, spiker.SymbolContinue.String()},
	}
	for _, tt := range tests {
		nc := spiker.NodeContinue{
			Ast: tt.fields.Ast,
		}
		if got := nc.String(); got != tt.want {
			t.Errorf("%q. NodeContinue.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeBreak_String(t *testing.T) {
	type fields struct {
		Ast spiker.Ast
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"return", fields{}, spiker.SymbolBreak.String()},
	}
	for _, tt := range tests {
		nb := spiker.NodeBreak{
			Ast: tt.fields.Ast,
		}
		if got := nb.String(); got != tt.want {
			t.Errorf("%q. NodeBreak.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestNodeReturn_String(t *testing.T) {
	type fields struct {
		Ast spiker.Ast
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"return", fields{}, spiker.SymbolReturn.String()},
	}
	for _, tt := range tests {
		nr := spiker.NodeReturn{
			Ast: tt.fields.Ast,
		}
		if got := nr.String(); got != tt.want {
			t.Errorf("%q. NodeReturn.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
