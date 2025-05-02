package parser

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type functionType int

const (
	none functionType = iota
	function
	initializer
	method
)

type Resolver struct {
	*Interpreter
	*scopes
	currentFunction functionType
	currentClass    classType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		scopes:          &scopes{},
		currentFunction: none,
		currentClass:    NONE,
	}
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
	for i := len(*r.scopes) - 1; i >= 0; i-- {
		s := (*r.scopes)[i]
		if _, defined := (*s)[name.Lexeme]; defined {
			depth := len(*r.scopes) - 1 - i
			r.Interpreter.Resolve(expr, depth)
			// s.use(name.Lexeme)
			return
		}
	}
}

func (r *Resolver) ResolveStmts(stmts []Stmt) {
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
	r.ResolveStmts(bs.statments)
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

	if len(*r.scopes) > 0 {
		if declared, defined := (*r.scopes.peek())[expr.name.Lexeme]; declared && !defined {
			panic("cannot read local variable in its own initializer")
		}
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
	// r.define(fs.name)
	r.resolveFunction(fs, function)
	return nil
}

func (r *Resolver) resolveFunction(fs *FunctionStmt, ft functionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ft

	r.beginScope()
	for _, param := range fs.params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmts(fs.body)
	r.endScope()

	r.currentFunction = enclosingFunction
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
	if r.currentFunction == none {
		panic("cannot return from top-level code.")
	}
	if stmt.value != nil {
		if r.currentFunction == initializer {
			panic("cannot return from init.")

		}
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

func (r *Resolver) VisitLogical(expr *Logical) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) VisitUnary(expr *Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) visitClassStmt(stmt *ClassStmt) any {
	cc := r.currentClass
	r.currentClass = CLASS

	r.declare(stmt.name)
	if stmt.superclass != nil {
		r.currentClass = SUBCLASS
		r.resolveExpr(stmt.superclass)
	}
	r.define(stmt.name)

	if stmt.superclass != nil && stmt.name.Lexeme == stmt.superclass.name.Lexeme {
		panic("A class cannot inherit from itself")
	}

	if stmt.superclass != nil {
		r.currentClass = SUBCLASS
		r.resolveExpr(stmt.superclass)
		r.beginScope()
		(*r.scopes.peek())["super"] = true
	}

	r.beginScope()
	(*r.scopes.peek())["this"] = true

	for _, m := range stmt.methods {

		if m.name.Lexeme == "this" {
			r.resolveFunction(m, initializer)
		} else {

			r.resolveFunction(m, method)
		}

	}

	r.endScope()

	if stmt.superclass != nil {
		r.endScope()
	}

	r.currentClass = cc

	return nil
}

func (r *Resolver) VisitGet(expr *Get) any {
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) VisitSet(expr *Set) any {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) VisitThis(expr *This) any {
	if r.currentClass == NONE {
		panic("Cannot use 'this' outside of a class.")
	}
	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) VisitSuper(expr *Super) any {
	if r.currentClass == NONE {
		panic("cannot use 'super' outside of a class")
	} else if r.currentClass != SUBCLASS {
		panic("cannot use 'super' in a class with no superclass")
	}
	r.resolveLocal(expr, expr.keyword)
	return nil
}
