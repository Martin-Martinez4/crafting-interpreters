package parser

import "time"

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) any
	arity() int
}

type Clock struct {
}

func (c *Clock) arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, arguments []any) any {
	return time.Now().Unix()
}

func (c *Clock) String() string {
	return "<native fn 'clock'> prints the time"
}

func (f *FunctionStmt) Call(interpreter *Interpreter, arguments []any) any {
	env := NewEnvironment(interpreter.globals)

	for i := 0; i < len(f.params); i++ {
		env.define(f.params[i].Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.body, env)
	return nil
}

func (f *FunctionStmt) arity() int {
	return len(f.params)
}

func (f *FunctionStmt) String() string {
	return "<fn " + f.name.Lexeme + ">"
}
