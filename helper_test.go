package spiker_test

import (
	"testing"

	"github.com/shockerli/spiker"
)

type TestParseNumberExpected struct {
	number float64
	err    bool
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		TestParseNumberExpected
	}{
		{"abc", TestParseNumberExpected{0, false}},
		{"abc12.43", TestParseNumberExpected{0, false}},
		{"12.43", TestParseNumberExpected{12.43, true}},
		{"12.43abc", TestParseNumberExpected{12.43, true}},
		{"0.1abc", TestParseNumberExpected{0.1, true}},
		{"-12.43", TestParseNumberExpected{-12.43, true}},
		{"-12.430", TestParseNumberExpected{-12.43, true}},
		{"-abc", TestParseNumberExpected{0, false}},
		{".abc", TestParseNumberExpected{0, false}},
		{"0.abc", TestParseNumberExpected{0, true}},
		{"-0.abc", TestParseNumberExpected{-0, true}},
		{"-12.43abc", TestParseNumberExpected{-12.43, true}},
		{"abc-12.43", TestParseNumberExpected{0, false}},
	}

	for index, tt := range tests {
		num, err := spiker.ParseNumber(tt.input)
		if tt.TestParseNumberExpected.number != num || (err == nil) != tt.TestParseNumberExpected.err {
			t.Errorf("test[%d], wrong operator precedence, expected = %g, got = %g",
				index, tt.TestParseNumberExpected.number, num)
		}
	}
}

func BenchmarkParseNumber(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = spiker.ParseNumber("-12.43abc")
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"abc", false},
		{"abc12.43", false},
		{"12.43", true},
		{"12.430", true},
		{"12.43abc", false},
		{"0.1abc", false},
		{"-12.43", true},
		{"-12.4300", true},
		{"-abc", false},
		{".abc", false},
		{"0.abc", false},
		{"-0.abc", false},
		{"-12.43abc", false},
		{"abc-12.43", false},
	}

	for index, tt := range tests {
		isN := spiker.IsNumber(tt.input)
		if isN != tt.expect {
			t.Errorf("test[%d], wrong operator precedence, expected = %t, got = %t",
				index, tt.expect, isN)
		}
	}
}

func TestIsTrue(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect bool
	}{
		{"abc", true},
		{"123", true},
		{"", false},
		{"0", true},
		{"-0.12", true},
		{123, true},
		{-123, true},
		{12.01, true},
		{-12.01, true},
		{0, false},
		{0.00, false},
		{-0, false},
		{true, true},
		{false, false},
		{spiker.ValueList{1, 2, 3}, true},
		{spiker.ValueList{spiker.ValueList{}}, true},
		{spiker.ValueList{}, false},
		{spiker.ValueMap{}, false},
		{spiker.ValueMap{"t": "123"}, true},
		{spiker.ValueMap{"t": spiker.ValueMap{"b": 123}}, true},
	}

	for index, tt := range tests {
		val := spiker.IsTrue(tt.input)
		if val != tt.expect {
			t.Errorf("test[%d], wrong operator precedence, expected = %t, got = %t",
				index, tt.expect, val)
		}
	}
}

func TestInterface2String(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect string
	}{
		{"abc", "abc"},
		{"123", "123"},
		{"", ""},
		{"0", "0"},
		{"-0.12", "-0.12"},
		{123, "123"},
		{-123, "-123"},
		{12.01, "12.01"},
		{-12.01, "-12.01"},
		{0, "0"},
		{0.00, "0"},
		{-0, "0"},
		{true, "1"},
		{false, ""},
		{spiker.ValueList{1, 2, 3}, "[1,2,3]"},
		{spiker.ValueList{spiker.ValueList{}}, "[[]]"},
		{spiker.ValueList{}, "[]"},
		{spiker.ValueMap{}, "{}"},
		{spiker.ValueMap{"t": "123"}, `{"t":"123"}`},
		{spiker.ValueMap{"t": spiker.ValueMap{"b": 123}}, `{"t":{"b":123}}`},
	}

	for index, tt := range tests {
		val := spiker.Interface2String(tt.input)
		if val != tt.expect {
			t.Errorf("test[%d], wrong operator precedence, expected = %s, got = %s",
				index, tt.expect, val)
		}
	}
}

func TestInterface2Float64(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect float64
	}{
		{"abc", 0},
		{"123", 123},
		{"", 0},
		{"0", 0},
		{"-0.12", -0.12},
		{123, 123},
		{-123, -123},
		{12.01, 12.01},
		{-12.01, -12.01},
		{0, 0},
		{0.00, 0},
		{-0, 0},
		{true, 1},
		{false, 0},
		{spiker.ValueList{1, 2, 3}, 0},
		{spiker.ValueList{spiker.ValueList{}}, 0},
		{spiker.ValueList{}, 0},
		{spiker.ValueMap{}, 0},
		{spiker.ValueMap{"t": "123"}, 0},
		{spiker.ValueMap{"t": spiker.ValueMap{"b": 123}}, 0},
	}

	for index, tt := range tests {
		val := spiker.Interface2Float64(tt.input)
		if val != tt.expect {
			t.Errorf("test[%d], wrong operator precedence, expected = %f, got = %f",
				index, tt.expect, val)
		}
	}
}
