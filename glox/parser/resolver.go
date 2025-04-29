package parser

type Resolver struct {
	Interpreter
}

func (r *Resolver) visitBlockStmt(bs *BlockStmt) any {
	r.beginScope()
	r.resolveStmts(bs.statments)
	r.endScope()
	return nil
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
