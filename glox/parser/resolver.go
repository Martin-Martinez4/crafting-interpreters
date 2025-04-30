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
