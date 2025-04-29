package parser

import (
	"fmt"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/errorhandling"
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		// statements = append(statements, p.statement())
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(token.FUN) {
		return p.function("function")
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string) *FunctionStmt {
	name, err := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		panic(err.Error())
	}

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after"+kind+" name.")
	if err != nil {
		panic(err.Error())
	}

	parameters := []*token.Token{}

	if !p.check(token.RIGHT_PAREN) {
		if len(parameters) >= 255 {
			panic(fmt.Sprintf("%v cannot have more than 255 parameters.", p.peek()))
		}
		tt, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
		if err != nil {
			panic(err.Error())
		}
		parameters = append(parameters, tt)

		for p.match(token.COMMA) {
			if len(parameters) >= 255 {
				panic(fmt.Sprintf("%v cannot have more than 255 parameters.", p.peek()))
			}
			tt, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				panic(err.Error())
			}
			parameters = append(parameters, tt)
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		panic(err.Error())
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before"+kind+" body.")
	if err != nil {
		panic(err.Error())
	}
	body := p.block()
	return &FunctionStmt{name: name, params: parameters, body: body}

}

func (p *Parser) varDeclaration() Stmt {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		panic(err.Error())
	}

	var initializer Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		panic(err.Error())
	}
	return &VarStmt{name: name, initializer: initializer}

}

func (p *Parser) statement() Stmt {
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return &BlockStmt{p.block()}
	}

	return p.expressionStatement()
}
func (p *Parser) whileStatement() Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after while.")
	if err != nil {
		panic(err.Error())
	}

	condition := p.expression()

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after while condition.")
	if err != nil {
		panic(err.Error())
	}

	body := p.statement()

	return &WhileStmt{condition: condition, body: body}

}

func (p *Parser) forStatement() Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()

	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()

	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()

	if increment != nil {
		body = &BlockStmt{statments: []Stmt{body, &ExprStmt{Expr: increment}}}
	}

	if condition == nil {
		condition = &Literal{Value: true}
	}
	body = &WhileStmt{body: body, condition: condition}

	if initializer != nil {
		body = &BlockStmt{statments: []Stmt{initializer, body}}
	}

	return body

}

func (p *Parser) ifStatement() Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after if.")
	if err != nil {
		panic(err.Error())
	}

	condition := p.expression()

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		panic(err.Error())
	}

	thenBranch := p.statement()
	var elseBranch Stmt = nil

	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil

	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	_, err := p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		panic(err.Error())
	}

	return &ReturnStmt{keyword: keyword, value: value}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		panic(err.Error())
	}
	return &PrintStmt{value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		panic(err.Error())
	}
	return &ExprStmt{expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		e, ok := expr.(*Variable)
		if !ok {
			panic(equals.Lexeme + " invalid assignment target.")
		}

		return NewAssignExpr(e.name, value)

	}

	return expr

}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) match(types ...token.TokenType) bool {
	for i := 0; i < len(types); i++ {
		t := types[i]
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) peek() *token.Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnaryExpr(operator, right)
	}

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}

	if !p.check(token.RIGHT_PAREN) {
		if len(arguments) >= 255 {
			panic("called function with more than 254 arguments")
		}
		arguments = append(arguments, p.expression())

		for p.match(token.COMMA) {
			arguments = append(arguments, p.expression())

		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		panic(err.Error())
	}

	return NewCallExpr(callee, paren, arguments)
}

func (p *Parser) primary() Expr {
	if p.match(token.FALSE) {
		return NewLiteralExpr(false)
	}
	if p.match(token.TRUE) {
		return NewLiteralExpr(true)
	}
	if p.match(token.NIL) {
		return NewLiteralExpr(nil)
	}

	if p.match(token.NUMBER, token.STRING) {
		return NewLiteralExpr(p.previous().Literal)
	}

	if p.match(token.IDENTIFIER) {
		return NewVariableExpr(p.previous())
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		_, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			errorhandling.ReportAndExit(0, "", "Expect ')' after expression.")
		}
		return NewGroupingExpr(expr)
	}

	panic(fmt.Sprintf("error at line %d: unknown token '%s'", p.peek().Line, p.peek().Literal))
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		panic(err.Error())
	}
	return statements
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, fmt.Errorf("%v: %s", t, message)
}

func parserErro(t *token.Token, message string) {
	if t.Type == token.EOF {
		errorhandling.ReportError(t.Line, " at end", message)
	} else {
		errorhandling.ReportError(t.Line, fmt.Sprintf(" at '%v'", t.Lexeme), message)
	}
}
