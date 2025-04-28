package parser

import "github.com/Martin-Martinez4/crafting-interpreters/glox/token"

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

func (e *Environment) Get(name *token.Token) any {
	v, ok := e.values[name.Lexeme]
	if !ok {

		if e.enclosing != nil {

			e := e.enclosing.Get(name)
			if !ok {

				panic("undefined variable '" + name.Lexeme + "'.")
			} else {
				return e
			}
		} else {
			panic("undefined variable '" + name.Lexeme + "'.")
		}

	}
	return v
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name *token.Token, value any) {
	_, ok := e.values[name.Lexeme]
	if !ok {

		if e.enclosing != nil {
			e.enclosing.Assign(name, value)
		} else {

			panic("undefined variable '" + name.Lexeme + "'.")
		}

	}

	e.define(name.Lexeme, value)
}
