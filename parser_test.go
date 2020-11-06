package spiker_test

import (
	"io/ioutil"
	"testing"

	"github.com/shockerli/spiker"
)

var srcTests = []string{
	"testdata/value.src",
	"testdata/assign.src",
	"testdata/collect.src",
	"testdata/operator.src",
}

func readFile(file string) string {
	if src, err := ioutil.ReadFile(file); err == nil {
		return string(src)
	}

	return ""
}

func BenchmarkStatements(b *testing.B) {
	src := readFile("testdata/collect.src")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		lexer := spiker.NewLexer(src)
		p := spiker.Parser{Lexer: lexer}
		if _, err := p.Statements(); err != nil {
			b.Log(err)
			b.Fail()
		}
	}
}

func TestParser_Statements(t *testing.T) {
	for _, file := range srcTests {
		src := readFile(file)
		lexer := spiker.NewLexer(src)
		p := spiker.Parser{Lexer: lexer}

		if _, err := p.Statements(); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}
