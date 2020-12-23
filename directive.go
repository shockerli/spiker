package spiker

// directive `return` statement
type directiveReturn struct {
	hasVal bool
	val    interface{}
}

// directive `export` statement
type directiveExport struct {
	val interface{}
}

// directive `break` statement
type directiveBreak struct{}

// directive `continue` statement
type directiveContinue struct{}
