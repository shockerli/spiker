package spiker_test

import (
	"testing"

	"github.com/shockerli/spiker"
)

func TestEvaluator(t *testing.T) {
	type args struct {
		nodeList []spiker.AstNode
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{"assign-string", args{[]spiker.AstNode{
			&spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: &spiker.NodeString{Value: "judy"}},
		}}, "judy", false},
		{"assign-number", args{[]spiker.AstNode{
			&spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "age"}, Op: spiker.SymbolAssign, Expr: &spiker.NodeNumber{Value: 18}},
		}}, 18, false},
		{"assign-add-string", args{[]spiker.AstNode{
			&spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssignAdd, Expr: &spiker.NodeString{Value: "judy"}},
		}}, "judy", false},
	}

	for _, tt := range tests {
		gotRes, err := spiker.Evaluator(tt.args.nodeList)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Evaluator() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if spiker.Interface2String(gotRes) != spiker.Interface2String(tt.wantRes) {
			t.Errorf("%q. Evaluator() = %v, want %v", tt.name, gotRes, tt.wantRes)
		}
	}
}

func TestEvaluateWithScope(t *testing.T) {
	type args struct {
		nodeList []spiker.AstNode
		scope    *spiker.VariableScope
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr bool
	}{
		{"assign-string", args{[]spiker.AstNode{
			&spiker.NodeAssignOp{Var: spiker.NodeVariable{Value: "name"}, Op: spiker.SymbolAssign, Expr: &spiker.NodeString{Value: "judy"}},
		}, func() *spiker.VariableScope {
			s := spiker.NewScopeTable("assign-string", 1, nil)
			s.Set("name", "tom")
			return s
		}()}, "judy", false},
	}

	for _, tt := range tests {
		gotRes, err := spiker.EvaluateWithScope(tt.args.nodeList, tt.args.scope)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. EvaluateWithScope() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if spiker.Interface2String(gotRes) != spiker.Interface2String(tt.wantRes) {
			t.Errorf("%q. EvaluateWithScope() = %v, want %v", tt.name, gotRes, tt.wantRes)
		}
	}
}
