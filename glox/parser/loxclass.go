package parser

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type classType int

const (
	NONE classType = iota
	CLASS
)

type Class struct {
	name    string
	fields  map[string]any
	methods map[string]*Function
}

type LoxInstance struct {
	*Class
}

func NewLoxInstance(c *Class) *LoxInstance {
	return &LoxInstance{
		c,
	}
}

func (li *LoxInstance) Get(name *token.Token) any {
	if v, ok := li.fields[name.Lexeme]; ok {
		return v
	}

	if m, ok := li.methods[name.Lexeme]; ok {
		return m.bind(li)
	}

	panic("undefined property '" + name.Lexeme + "'.")
}

func (li *LoxInstance) Set(name *token.Token, value any) {
	if li.fields == nil {
		li.fields = make(map[string]interface{})
	}
	li.fields[name.Lexeme] = value
}

func (li *LoxInstance) String() string {
	return li.name + " instance"
}

func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{
		name:    name,
		fields:  map[string]any{},
		methods: methods,
	}
}

func (lc *Class) String() string {
	return lc.name
}

func (lc *Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewLoxInstance(lc)

	initializer, ok := lc.methods["this"]
	if ok {
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (lc *Class) arity() int {
	initializer, ok := lc.methods["this"]
	if !ok {
		return 0
	}
	return initializer.arity()
}
