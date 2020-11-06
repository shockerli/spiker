package spiker_test

import (
	"fmt"

	"github.com/shockerli/spiker"
)

// ExampleExecute
func ExampleExecute() {
	src := `
a = 100;
a = len("abc");
if (a > 1) {
    b = 10;
}
export(b + "10");
`
	fmt.Println(spiker.Execute(src))
	// Output: 20 <nil>
}
