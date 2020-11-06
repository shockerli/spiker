package spiker

type tokenRegistry struct {
	symTable map[Symbol]*Token
}

func (reg *tokenRegistry) token(sym Symbol, value string, line int, col int) *Token {
	return &Token{
		sym:          sym,
		value:        value,
		line:         line,
		col:          col,
		bindingPower: reg.symTable[sym].bindingPower,
		nud:          reg.symTable[sym].nud,
		led:          reg.symTable[sym].led,
		std:          reg.symTable[sym].std,
	}
}

func (reg *tokenRegistry) defined(sym Symbol) bool {
	if _, ok := reg.symTable[sym]; ok {
		return true
	}
	return false
}

func (reg *tokenRegistry) register(sym Symbol, bp int, nud nudFn, led ledFn, std stdFn) {
	if val, ok := reg.symTable[sym]; ok {
		if nud != nil && val.nud == nil {
			val.nud = nud
		}
		if led != nil && val.led == nil {
			val.led = led
		}
		if std != nil && val.std == nil {
			val.std = std
		}
		if bp > val.bindingPower {
			val.bindingPower = bp
		}
	} else {
		reg.symTable[sym] = &Token{bindingPower: bp, nud: nud, led: led, std: std}
	}
}

// an infix Token has two children, the exp on the left and the one that follows
func (reg *tokenRegistry) infix(sym Symbol, bp int) {
	reg.register(sym, bp, nil, func(t *Token, p *Parser, left *Token) *Token {
		t.children = append(t.children, left)
		t.children = append(t.children, p.expression(t.bindingPower))
		return t
	}, nil)
}

func (reg *tokenRegistry) infixLed(sym Symbol, bp int, led ledFn) {
	reg.register(sym, bp, nil, led, nil)
}

func (reg *tokenRegistry) infixRight(sym Symbol, bp int) {
	reg.register(sym, bp, nil, func(t *Token, p *Parser, left *Token) *Token {
		t.children = append(t.children, left)
		t.children = append(t.children, p.expression(t.bindingPower-1))
		return t
	}, nil)
}

func (reg *tokenRegistry) infixRightLed(sym Symbol, bp int, led ledFn) {
	reg.register(sym, bp, nil, led, nil)
}

// a prefix Token has a single children, the expression that follows
func (reg *tokenRegistry) prefix(sym Symbol) {
	reg.register(sym, 0, func(t *Token, p *Parser) *Token {
		t.children = append(t.children, p.expression(100))
		return t
	}, nil, nil)
}

func (reg *tokenRegistry) prefixNud(sym Symbol, nud nudFn) {
	reg.register(sym, 0, nud, nil, nil)
}

func (reg *tokenRegistry) stmt(sym Symbol, std stdFn) {
	reg.register(sym, 0, nil, nil, std)
}

func (reg *tokenRegistry) symbol(sym Symbol) {
	reg.register(sym, 0, func(t *Token, p *Parser) *Token { return t }, nil, nil)
}

func (reg *tokenRegistry) consumable(sym Symbol) {
	reg.register(sym, 0, nil, nil, nil)
}
