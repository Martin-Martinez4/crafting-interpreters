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
	VisitLogical(expr *Logical) any
	VisitCall(expr *CallExpr) any
	VisitGet(expr *Get) any
	VisitSet(expr *Set) any
	VisitThis(expr *This) any
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

type CallExpr struct {
	callee    Expr
	paren     *token.Token
	arguments []Expr
	defaultStartEnd
}

func NewCallExpr(callee Expr, paren *token.Token, arguments []Expr) *CallExpr {
	return &CallExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func (ce *CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCall(ce)
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

type This struct {
	keyword *token.Token
	defaultStartEnd
}

func (t *This) Accept(visitor ExprVisitor) any {
	return visitor.VisitThis(t)
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

type Logical struct {
	left     Expr
	operator *token.Token
	right    Expr
	defaultStartEnd
}

func NewLogical(left Expr, operator *token.Token, right Expr) *Logical {
	return &Logical{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (l *Logical) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogical(l)
}

type Get struct {
	object Expr
	name   *token.Token
	defaultStartEnd
}

func NewGet(object Expr, name *token.Token) any {
	return &Get{
		object: object,
		name:   name,
	}
}

func (g *Get) Accept(visitor ExprVisitor) any {
	return visitor.VisitGet(g)
}

type Set struct {
	object Expr
	name   *token.Token
	value  Expr
	defaultStartEnd
}

func (s *Set) Accept(visitor ExprVisitor) any {
	return visitor.VisitSet(s)
}
