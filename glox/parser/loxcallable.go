package parser

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) any
	arity() int
}

type Function struct {
	declaration *FunctionStmt
	closure     *Environment
	isInit      bool
}

func NewFunciton(declaration *FunctionStmt, closure *Environment, isInit bool) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
		isInit:      isInit,
	}
}

func (f *Function) bind(instance *LoxInstance) *Function {
	env := NewEnvironment(f.closure)
	env.define("this", instance)
	return NewFunciton(f.declaration, env, f.isInit)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) (returnVal any) {
	defer func() {
		if err := recover(); err != nil {
			if f.isInit {
				returnVal = f.closure.getAt(0, "this")
				return
			}
			if v, ok := err.(Return); ok {
				if f.isInit {
					returnVal = f.closure.getAt(0, "this")
					return
				}

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
