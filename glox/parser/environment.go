package parser

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    map[string]any{},
		enclosing: enclosing,
	}
}

func (e *Environment) Get(name *token.Token) any {
	v, ok := e.values[name.Lexeme]
	if ok {

		return v

	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic("boom")
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name *token.Token, value any) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.define(name.Lexeme, value)
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
	} else {

		panic("undefined variable '" + name.Lexeme + "'.")
	}

}

func (e *Environment) AssignAt(distance int, name *token.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}
