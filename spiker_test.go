package spiker_test

import (
	"reflect"
	"testing"

	"github.com/shockerli/spiker"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{`10;`, float64(10)},
		{`1 + 2 - 3 * 4 / 5;`, 0.6},
		{`1 + "234";`, float64(235)},
		{`1 + "2.34";`, 3.34},
		{`1 + "0.234";`, 1.234},
		{`1 + "234a";`, "1234a"},
		{`1 + "0.234a";`, "10.234a"},
		{`1 + "abc234";`, "1abc234"},
		{`123 in "abc234";`, false},
		{`23 in "abc234";`, true},
		{`2 in [1,2,3];`, true},
		{`!a`, true},
		{`~8`, -9},
		{`
add = (a, b) -> {
  return a + b;
};

c = add(1, 2);
export(c);
`, float64(3)},

{`
add = (a, b) -> {
  return a + b;
};

a = 0;
b = 3;
while (true) {
	a += add(a, 3);
	if (a > 10) {
		break;
	}
}
export(a);
`, float64(21)},

{`
n2 = x -> x * x;
a = n2(5);
export(a);
`, float64(25)},
	}

	for index, tt := range tests {
		val, err := spiker.Execute(tt.input)
		if val != tt.expect {
			t.Errorf("test[%d], expected = %v, got = %v", index, tt.expect, val)
		} else if err != nil {
			t.Error(err.Error())
		}
	}
}

func BenchmarkExecute(b *testing.B) {
	src := readFile("testdata/collect.src")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := spiker.Execute(src); err != nil {
			b.Log(err)
			b.Fail()
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`a=1;b+=2;`, `a = 1;
b += 2;`},
	}
	for index, tt := range tests {
		val, err := spiker.Format(tt.input)
		println(val)
		println(tt.expect)
		if val != tt.expect {
			t.Errorf("test[%d], expected: %v, got: %v",
				index,
				tt.expect,
				val,
			)
		} else if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestExecuteWithScope(t *testing.T) {
	var scopes = spiker.NewScopeTable("demo", 1, nil)
	scopes.Set("a", 3)
	scopes.Set("b", 4)

	tests := []struct {
		name    string
		code    string
		scope   *spiker.VariableScope
		wantVal interface{}
		wantErr bool
	}{
		{"nil", "a * b", nil, nil, true},
		{"mul", "a * b", scopes, float64(12), false},
		{"add", "a + b", scopes, float64(7), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := spiker.ExecuteWithScope(tt.code, tt.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("ExecuteWithScope() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
