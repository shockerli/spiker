package spiker_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/c5433137/spiker"
)

func TestRegisterFunc(t *testing.T) {
	type args struct {
		name string
		fn   spiker.Func
		code string
	}
	tests := []args{
		{"log", func(fnc *spiker.NodeFuncCallOp, scope *spiker.VariableScope) interface{} {
			for _, p := range fnc.Params {
				log.Println(spiker.EvalExpr(p, scope))
			}
			return nil
		}, `log(123, "abc")`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spiker.RegisterFunc(tt.name, tt.fn)
			if _, err := spiker.Execute(tt.code); err != nil {
				t.Errorf("register func failed, error = %v", err)
			}
		})
	}
}

func TestBuiltin_Export(t *testing.T) {
	type args struct {
		name   string
		code   string
		expect interface{}
	}
	tests := []args{
		{`export-var-string`, `name="jioby";export(name);`, "jioby"},
		{`export-var-int`, `age=18;export(age);`, float64(18)},
		{`export-list`, `export([1,2,3]);`, spiker.ValueList{float64(1), float64(2), float64(3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := spiker.Execute(tt.code)
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(res, tt.expect) {
				t.Errorf("want = %v, got = %v", tt.expect, res)
			}
		})
	}
}

func TestBuiltin_Len(t *testing.T) {
	type args struct {
		name   string
		code   string
		expect interface{}
	}
	tests := []args{
		{`len-var-string`, `name="jioby";len(name);`, 5},
		{`len-list`, `len([1,2,3]);`, 3},
		{`len-map`, `len([1:11,2:22,3:33]);`, 3},
		{`len-float`, `len(12.34);`, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := spiker.Execute(tt.code)
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}

			if !reflect.DeepEqual(res, tt.expect) {
				t.Errorf("want = %v, got = %v", tt.expect, res)
			}
		})
	}
}
