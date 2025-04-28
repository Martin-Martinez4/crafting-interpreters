package parser

import "github.com/Martin-Martinez4/crafting-interpreters/glox/token"

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitPrintStmt(*PrintStmt) error
	visitExpressionStmt(*ExprStmt) error
	visitVariableStmt(*VarStmt) error
}

type PrintStmt struct {
	Expr
}

func (p *PrintStmt) Accept(v StmtVisitor) error {
	return v.visitPrintStmt(p)
}

type ExprStmt struct {
	Expr
}

func (e *ExprStmt) Accept(v StmtVisitor) error {
	return v.visitExpressionStmt(e)
}

type VarStmt struct {
	name        *token.Token
	initializer Expr
}

func (vs *VarStmt) Accept(v StmtVisitor) error {
	return v.visitVariableStmt(vs)
}
