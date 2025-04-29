package parser

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) any
	arity() int
}

type Function struct {
	declaration *FunctionStmt
	closure     *Environment
}

func NewFunciton(declaration *FunctionStmt, closure *Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (returnVal any) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(Return); ok {

				returnVal = v.value
				return
			}
			panic(err)
		}
	}()

	env := NewEnvironment(f.closure)

	for i, v := range f.declaration.params {
		env.define(v.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.body, env)

	return nil
}

func (f *Function) arity() int {
	return len(f.declaration.params)
}

func (f *Function) String() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}
