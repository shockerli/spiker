package spiker_test

import (
	"testing"

	"github.com/shockerli/spiker"
)

func BenchmarkTransform(b *testing.B) {
	src := readFile("testdata/collect.src")
	lexer := spiker.NewLexer(src)
	p := spiker.Parser{Lexer: lexer}
	stmts, err := p.Statements()
	if err != nil {
		b.Log(err)
		b.Fail()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := spiker.Transform(stmts); err != nil {
			b.Log(err)
			b.Fail()
		}
	}
}
