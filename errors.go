package spiker

import "fmt"

// error `return` statement
type errReturn struct {
	hasVal bool
	val    interface{}
}

func (e errReturn) Error() string {
	return fmt.Sprintf("return value: %v", e.val)
}

// error `export` statement
type errExport struct {
	val interface{}
}

func (e errExport) Error() string {
	return fmt.Sprintf("export value: %v", e.val)
}

// error `break` statement
type errBreak struct{}

func (e errBreak) Error() string {
	return "break statement"
}

// error `continue` statement
type errContinue struct{}

func (e errContinue) Error() string {
	return "continue statement"
}
