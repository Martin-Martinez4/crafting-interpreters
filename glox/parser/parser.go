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
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		panic("Boo")
	}

	var initializer Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return &VarStmt{name: name, initializer: initializer}

}

func (p *Parser) statement() Stmt {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &PrintStmt{value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return &ExprStmt{expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		e, ok := expr.(*Variable)
		if !ok {
			panic(equals.Lexeme + " invalid assignment target.")
		}

		return NewAssignExpr(e.name, value)

	} else {

		return expr
	}

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

	return p.primary()
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

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, fmt.Errorf("%v: %s", message)
}

func parserErro(t *token.Token, message string) {
	if t.Type == token.EOF {
		errorhandling.ReportError(t.Line, " at end", message)
	} else {
		errorhandling.ReportError(t.Line, fmt.Sprintf(" at '%v'", t.Lexeme), message)
	}
}
