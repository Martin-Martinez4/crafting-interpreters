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
