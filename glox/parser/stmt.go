package parser

import "github.com/Martin-Martinez4/crafting-interpreters/glox/token"

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitPrintStmt(*PrintStmt) error
	visitExpressionStmt(*ExprStmt) error
	visitVariableStmt(*VarStmt) error
	visitBlockStmt(*BlockStmt) error
	visitIfStmt(*IfStmt) error
	visitWhileStmt(*WhileStmt) error
	visitFunctionStmt(*FunctionStmt) error
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

type BlockStmt struct {
	statments []Stmt
}

func (vb *BlockStmt) Accept(v StmtVisitor) error {
	return v.visitBlockStmt(vb)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *IfStmt) Accept(v StmtVisitor) error {
	return v.visitIfStmt(i)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (w *WhileStmt) Accept(v StmtVisitor) error {
	return v.visitWhileStmt(w)
}

type FunctionStmt struct {
	name   *token.Token
	params []*token.Token
	body   []Stmt
}

func (f *FunctionStmt) Accept(v StmtVisitor) error {
	return v.visitFunctionStmt(f)
}
