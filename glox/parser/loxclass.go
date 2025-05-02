package parser

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type classType int

const (
	NONE classType = iota
	CLASS
	SUBCLASS
)

type Class struct {
	name       string
	fields     map[string]any
	methods    map[string]*Function
	superclass *Class
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

	if m, ok := li.findMethod(name.Lexeme); ok {
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

func NewClass(name string, superclass *Class, methods map[string]*Function) *Class {
	return &Class{
		name:       name,
		fields:     map[string]any{},
		methods:    methods,
		superclass: superclass,
	}
}

func (c *Class) findMethod(name string) (*Function, bool) {
	var v *Function
	var ok bool

	v, ok = c.methods[name]
	if ok {
		return v, ok
	}

	if c.superclass != nil {
		v, ok = c.superclass.methods[name]
	}

	return v, ok
}

func (lc *Class) String() string {
	return lc.name
}

func (lc *Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewLoxInstance(lc)

	initializer, ok := lc.findMethod("this")
	if ok {
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (lc *Class) arity() int {
	initializer, ok := lc.findMethod("this")
	if !ok {
		return 0
	}
	return initializer.arity()
}
