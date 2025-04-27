package parser

import (
	"fmt"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/ast"
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

func (p *Parser) Parse() ast.Expr {
	return p.expression()
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinaryExpr(expr, operator, right)
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

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnaryExpr(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpr(false)
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpr(true)
	}
	if p.match(token.NIL) {
		return ast.NewLiteralExpr(nil)
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpr(p.previous().Literal)
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		_, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			errorhandling.ReportAndExit(0, "", "Expect ')' after expression.")
		}
		return ast.NewGroupingExpr(expr)
	}

	// may cause problems
	return nil
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
