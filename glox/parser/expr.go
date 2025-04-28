package parser

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
	StartLine() int
	EndLine() int
}

// default startline and endline
type defaultStartEnd struct{}

func (d defaultStartEnd) StartLine() int {
	// TODO implement me
	panic("implement me")
}

func (d defaultStartEnd) EndLine() int {
	// TODO implement me
	panic("implement me")
}

type ExprVisitor interface {
	VisitBinary(expr *Binary) any
	VisitGrouping(expr *Grouping) any
	VisitLiteral(expr *Literal) any

	VisitUnary(expr *Unary) any
	VisitVariable(expr *Variable) any
	VisitAssign(expr *Assign) any
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
	defaultStartEnd
}

func NewBinaryExpr(left Expr, operator *token.Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
	defaultStartEnd
}

func NewGroupingExpr(expression Expr) *Grouping {
	return &Grouping{
		Expression: expression,
	}
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
	defaultStartEnd
}

func NewLiteralExpr(value any) *Literal {
	return &Literal{
		Value: value,
	}
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator *token.Token
	Right    Expr
	defaultStartEnd
}

func NewUnaryExpr(operator *token.Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}

type Variable struct {
	name *token.Token
	defaultStartEnd
}

func NewVariableExpr(name *token.Token) *Variable {
	return &Variable{name: name}
}

func (v *Variable) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariable(v)
}

type Assign struct {
	name  *token.Token
	value Expr
	defaultStartEnd
}

func NewAssignExpr(name *token.Token, value Expr) *Assign {
	return &Assign{
		name:  name,
		value: value,
	}
}

func (a *Assign) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssign(a)
}
