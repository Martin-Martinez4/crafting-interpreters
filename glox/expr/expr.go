package expr

import "go/token"

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
	VisitBinary(expr Binary) any
	VisitGrouping(expr Grouping) any
	VisitLiteral(expr Literal) any

	VisitUnary(expr Unary) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (l Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
	defaultStartEnd
}

func (u Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}
