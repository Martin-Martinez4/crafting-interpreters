package parser

import (
	"fmt"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Resolver struct {
	Interpreter
	*scopes
}

type scope = map[string]bool
type scopes []*scope

func (s *scopes) peek() *scope {
	return (*s)[len(*s)-1]
}

func (s *scopes) push(scope *scope) {
	*s = append(*s, scope)
}

func (s *scopes) pop() {
	*s = (*s)[:len(*s)-1]
}

func (r *Resolver) resolveLocal(expr Expr, name *token.Token) {
	for i, sc := range *r.scopes {
		_, ok := (*sc)[name.Lexeme]
		if ok {
			r.Resolve(expr, len(*r.scopes)-1-i)
			return
		}
	}

	return
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[fmt.Sprintf("%v", expr)] = depth
}

func (r *Resolver) resolveStmts(stmts []Stmt) {
	for _, s := range stmts {
		r.resolveStmt(s)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.push(&scope{})
}

func (r *Resolver) endScope() {
	r.pop()
}

func (r *Resolver) declare(name *token.Token) {
	if len(*r.scopes) == 0 {
		return
	}

	scope := r.peek()

	_, ok := (*scope)[name.Lexeme]
	if ok {
		return
	}

	(*scope)[name.Lexeme] = false
}

func (r *Resolver) define(name *token.Token) {
	if len(*r.scopes) == 0 {
		return
	}

	(*r.peek())[name.Lexeme] = true
}

func (r *Resolver) visitBlockStmt(bs *BlockStmt) any {
	r.beginScope()
	r.resolveStmts(bs.statments)
	r.endScope()
	return nil
}

func (r *Resolver) visitVariableStmt(vs *VarStmt) any {
	r.declare(vs.name)
	if vs.initializer != nil {
		r.resolveExpr(vs.initializer)
	}
	r.define(vs.name)
	return nil
}

func (r *Resolver) VisitVariable(expr *Variable) any {
	v, _ := (*r.scopes.peek())[expr.name.Lexeme]
	if (len(*r.scopes) > 0) && v == false {
		panic("cannot read local variable in its own initializer")
	}

	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) VisitAssign(expr *Assign) any {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) visitFunctionStmt(fs *FunctionStmt) any {
	r.declare(fs.name)
	r.define(fs.name)
	r.resolveFunction(fs)
	return nil
}

func (r *Resolver) resolveFunction(fs *FunctionStmt) {
	r.beginScope()
	for _, param := range fs.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(fs.body)
	r.endScope()
}

func (r *Resolver) visitExpressionStmt(stmt *ExprStmt) any {
	r.resolveExpr(stmt.Expr)
	return nil
}

func (r *Resolver) visitIfStmt(stmt *IfStmt) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(stmt *PrintStmt) any {
	r.resolveExpr(stmt.Expr)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt *ReturnStmt) any {
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
	return nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStmt) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil
}

func (r *Resolver) VisitBinary(expr *Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCall(expr *CallExpr) any {
	r.resolveExpr(expr.callee)

	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGrouping(expr *Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteral(expr *Literal) any {
	return nil
}
