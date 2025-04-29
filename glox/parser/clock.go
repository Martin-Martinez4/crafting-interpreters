package parser

import "time"

type clock struct {
}

func (c *clock) arity() int {
	return 0
}

func (c *clock) Call(interpreter *Interpreter, arguments []any) (returnVal any) {
	return float64(time.Now().UnixMilli())
}

func (c *clock) String() string {
	return "<native fn 'clock'> prints the time in milliseconds"
}
