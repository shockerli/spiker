package spiker

// VariableScope variable scope
type VariableScope struct {
	scopeName      string
	scopeLevel     int
	vars           map[string]interface{}
	enclosingScope *VariableScope
}

// NewScopeTable return a new VariableScope
func NewScopeTable(scopeName string, scopeLevel int, scope *VariableScope) *VariableScope {
	vs := &VariableScope{}
	vs.vars = make(map[string]interface{})
	vs.scopeName = scopeName
	vs.scopeLevel = scopeLevel
	vs.enclosingScope = scope
	return vs
}

// NewScopeTable return a new VariableScope
func NewScopeTableCap(scopeName string, scopeLevel int, scope *VariableScope, cap int) *VariableScope {
	vs := &VariableScope{}
	vs.vars = make(map[string]interface{}, cap)
	vs.scopeName = scopeName
	vs.scopeLevel = scopeLevel
	vs.enclosingScope = scope
	return vs
}

// Set store variable values
func (scope *VariableScope) Set(variable string, val interface{}) {
	scope.vars[variable] = val
}

// Get fetch variable values
func (scope *VariableScope) Get(variable string) (interface{}, bool) {
	if val, ok := scope.vars[variable]; ok {
		return val, true
	}
	if scope.enclosingScope != nil {
		return scope.enclosingScope.Get(variable)
	}

	return nil, false
}

// Del delete a variable
func (scope *VariableScope) Del(variable string) {
	delete(scope.vars, variable)
}

// Clean clean all of the vars
func (scope *VariableScope) Clean() {
	scope.vars = make(map[string]interface{})
}
