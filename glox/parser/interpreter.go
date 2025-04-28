package parser

import (
	"fmt"
	"reflect"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Interpreter struct {
	Statements  []Stmt
	environment *Environment
}

func (i *Interpreter) Interpret(statements []Stmt) {

	i.environment = NewEnvironment()

	for _, s := range statements {
		s.Accept(i)

		// switch s.(type) {
		// case *PrintStmt:
		// 	a := s.(*PrintStmt)
		// 	a.Accept(i)

		// }
	}
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
	return expr.Accept(i)
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
	return i.environment.Get(expr.name)
}

func (i *Interpreter) visitVariableStmt(vstmt *VarStmt) error {
	var value any = nil
	if vstmt.initializer != nil {
		value = vstmt.initializer.Accept(i)
	}

	i.environment.define(vstmt.name.Lexeme, value)
	return nil
}

func (i *Interpreter) visitPrintStmt(pstmt *PrintStmt) error {
	value := pstmt.Expr.Accept(i)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) visitExpressionStmt(etmt *ExprStmt) error {
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

func (i *Interpreter) VisitAssign(expr *Assign) any {
	value := expr.value.Accept(i)
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
