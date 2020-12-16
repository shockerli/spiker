package spiker_test

import (
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
