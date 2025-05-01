package parser

import (
	"fmt"
	"reflect"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Interpreter struct {
	Statements  []Stmt
	environment *Environment
	globals     *Environment
	locals      map[string]int
}

type Return struct {
	value any
}

func NewInterpreter(statements []Stmt) *Interpreter {
	e := NewEnvironment(nil)
	e.define("clock", &clock{})

	i := &Interpreter{
		Statements:  statements,
		environment: e,
		globals:     e,
		locals:      map[string]int{},
	}

	return i
}

func (i *Interpreter) Interpret(statements []Stmt) {

	for _, s := range statements {
		s.Accept(i)

	}
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[fmt.Sprintf("%v", expr)] = depth
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", object)

}

func (i *Interpreter) VisitLiteral(expr *Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGrouping(expr *Grouping) any {
	return expr.Expression.Accept(i)
}

func (i *Interpreter) VisitUnary(expr *Unary) any {
	right := expr.Right.Accept(i)

	switch expr.Operator.Type {
	case token.MINUS:
		r, ok := right.(float64)
		if !ok {
			panic("VisitUnary MINUS type expr.Right was expected to be a float64 and was not")
		}
		return -r

	case token.BANG:
		r, ok := right.(bool)
		if !ok {
			panic("VisitUnary BANG type expr.Right was expected to be a bool and was not")
		}
		return !isTruthy(r)
	}

	return nil
}

func (i *Interpreter) VisitVariable(expr *Variable) any {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) lookUpVariable(name *token.Token, expr Expr) any {
	distance, ok := i.locals[fmt.Sprintf("%v", expr)]
	if !ok {
		return i.globals.Get(name)
	} else {
		return i.environment.getAt(distance, name.Lexeme)
	}
}

func (e *Environment) getAt(distance int, name string) any {
	v, ok := e.ancestor(distance).values[name]
	if !ok {
		panic("Not found")
	}
	return v
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}

func (i *Interpreter) visitVariableStmt(vstmt *VarStmt) any {
	var value any = nil
	if vstmt.initializer != nil {
		value = vstmt.initializer.Accept(i)
	}

	i.environment.define(vstmt.name.Lexeme, value)
	return nil
}

func (i *Interpreter) visitPrintStmt(pstmt *PrintStmt) any {
	value := pstmt.Expr.Accept(i)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) visitExpressionStmt(etmt *ExprStmt) any {
	etmt.Expr.Accept(i)
	return nil
}

func checkNumberOperand(operator token.Token, operand any) (float64, error) {
	f, ok := operand.(float64)
	if !ok {
		return 0, fmt.Errorf("%v Operand must be a number.", operator)
	}
	return f, nil
}

func (i *Interpreter) VisitBinary(expr *Binary) any {
	left := expr.Left.Accept(i)
	right := expr.Right.Accept(i)

	switch expr.Operator.Type {
	case token.MINUS:

		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)

		return l - r

	case token.SLASH:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)

		return l / r

	case token.STAR:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)

		return l * r

	case token.PLUS:

		switch right.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			l, _ := checkNumberOperand(*expr.Operator, left)
			r, _ := checkNumberOperand(*expr.Operator, right)
			return l + r

		case string:
			r, ok := right.(string)
			if !ok {
				panic("VisitBinary PLUS type expr.Right was expected to be a string and was not")
			}
			l, ok := left.(string)
			if !ok {
				panic("VisitBinary PLUS type expr.Left was expected to be a string and was not")
			}

			return l + r

		default:
			panic("VisitBinary PLUS type expr.Right was expected to be a string or number and was not")
		}

	case token.GREATER:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)
		return l > r

	case token.GREATER_EQUAL:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)
		return l >= r

	case token.LESS:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)
		return l < r

	case token.LESS_EQUAL:
		l, _ := checkNumberOperand(*expr.Operator, left)
		r, _ := checkNumberOperand(*expr.Operator, right)
		return l <= r

	case token.BANG_EQUAL:
		return !isEqual(left, right)

	case token.EQUAL_EQUAL:
		return isEqual(left, right)

	default:
		return nil
	}

}

func (i *Interpreter) VisitCall(expr *CallExpr) any {
	callee := expr.callee.Accept(i)

	arguments := []any{}
	for _, arg := range expr.arguments {

		arguments = append(arguments, arg.Accept(i))
	}

	c, ok := callee.(LoxCallable)
	if ok {
		if len(arguments) != c.arity() {

			panic(fmt.Sprintf("%v Expected %d arguments but got %d", expr.paren, c.arity(), len(arguments)))
		}
		return c.Call(i, arguments)
	} else {
		panic(fmt.Sprintf("tried to call uncallable object %s; can only call functions and classes", reflect.TypeOf(callee)))
	}
}
func (i *Interpreter) VisitGet(expr *Get) any {
	obj := expr.object.Accept(i)

	o, ok := obj.(*LoxInstance)
	if !ok {
		panic("only instances have properties")
	}

	return o.Get(expr.name)
}

func (i *Interpreter) VisitAssign(expr *Assign) any {
	value := expr.value.Accept(i)

	distance, ok := i.locals[fmt.Sprintf("%v", expr)]
	if !ok {
		i.globals.Assign(expr.name, value)

	} else {
		i.environment.AssignAt(distance, expr.name, value)
	}

	i.environment.Assign(expr.name, value)
	return value
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}

	b, ok := object.(bool)
	if !ok {
		return true
	}
	return b
}

func (i *Interpreter) VisitLogical(expr *Logical) any {
	left := expr.left.Accept(i)

	if expr.operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return expr.right.Accept(i)
}

func (i *Interpreter) VisitSet(expr *Set) any {
	object := expr.object.Accept(i)

	o, ok := object.(*LoxInstance)
	if !ok {
		panic(expr.name.Lexeme + " only instances have fields.")
	}

	value := expr.value.Accept(i)
	o.Set(expr.name, value)
	return nil
}

func (i *Interpreter) visitBlockStmt(block *BlockStmt) any {
	i.executeBlock(block.statments, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) visitIfStmt(ifStmt *IfStmt) any {
	if isTruthy(ifStmt.condition.Accept(i)) {
		ifStmt.thenBranch.Accept(i)
	} else if ifStmt.elseBranch != nil {
		ifStmt.elseBranch.Accept(i)
	}

	return nil
}

func (i *Interpreter) visitClassStmt(cStmt *ClassStmt) any {
	i.environment.define(cStmt.name.Lexeme, nil)

	methods := map[string]*Function{}
	for _, m := range cStmt.methods {
		methods[m.name.Lexeme] = NewFunciton(m, i.environment, m.name.Lexeme == "this")
	}

	class := NewClass(cStmt.name.Lexeme, methods)
	i.environment.Assign(cStmt.name, class)
	return nil
}

func (i *Interpreter) visitWhileStmt(while *WhileStmt) any {
	for isTruthy(while.condition.Accept(i)) {
		while.body.Accept(i)
	}
	return nil
}

func (i *Interpreter) visitFunctionStmt(fun *FunctionStmt) any {
	f := NewFunciton(fun, i.environment, false)
	i.environment.define(fun.name.Lexeme, f)
	return nil
}

func (i *Interpreter) visitReturnStmt(r *ReturnStmt) any {

	// var v Expr = nil
	// if r.value != nil {
	// 	e, ok := r.value.Accept(i).(Expr)
	// 	if !ok {
	// 		panic("cast to Expr failed in visitReturnStmt")
	// 	}
	// 	v = e
	// }

	// return &ReturnStmt{value: v}
	var value any
	if r.value != nil {
		value = r.value.Accept(i)
	}
	panic(Return{value: value})

}

func (i *Interpreter) VisitThis(expr *This) any {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) any {
	prev := i.environment
	defer func() {
		i.environment = prev
	}()

	i.environment = env
	for _, statement := range statements {

		statement.Accept(i)
	}

	return nil
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else {
		aType := reflect.TypeOf(a)
		bType := reflect.TypeOf(b)
		if aType != bType {
			return false
		}
		return reflect.DeepEqual(a, b)

	}
}
