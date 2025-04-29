package parser

import "github.com/Martin-Martinez4/crafting-interpreters/glox/token"

type Stmt interface {
	Accept(StmtVisitor) any
}

type StmtVisitor interface {
	visitPrintStmt(*PrintStmt) any
	visitExpressionStmt(*ExprStmt) any
	visitVariableStmt(*VarStmt) any
	visitBlockStmt(*BlockStmt) any
	visitIfStmt(*IfStmt) any
	visitWhileStmt(*WhileStmt) any
	visitFunctionStmt(*FunctionStmt) any
	visitReturnStmt(*ReturnStmt) any
}

type PrintStmt struct {
	Expr
}

func (p *PrintStmt) Accept(v StmtVisitor) any {
	return v.visitPrintStmt(p)
}

type ExprStmt struct {
	Expr
}

func (e *ExprStmt) Accept(v StmtVisitor) any {
	return v.visitExpressionStmt(e)
}

type VarStmt struct {
	name        *token.Token
	initializer Expr
}

func (vs *VarStmt) Accept(v StmtVisitor) any {
	return v.visitVariableStmt(vs)
}

type BlockStmt struct {
	statments []Stmt
}

func (vb *BlockStmt) Accept(v StmtVisitor) any {
	return v.visitBlockStmt(vb)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *IfStmt) Accept(v StmtVisitor) any {
	return v.visitIfStmt(i)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (w *WhileStmt) Accept(v StmtVisitor) any {
	return v.visitWhileStmt(w)
}

type FunctionStmt struct {
	name   *token.Token
	params []*token.Token
	body   []Stmt
}

func (f *FunctionStmt) Accept(v StmtVisitor) any {
	return v.visitFunctionStmt(f)
}

type ReturnStmt struct {
	keyword *token.Token
	value   Expr
}

func (r *ReturnStmt) Accept(v StmtVisitor) any {
	v.visitReturnStmt(r)
	return nil
}
