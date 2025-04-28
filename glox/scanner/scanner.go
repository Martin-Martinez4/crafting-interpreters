package scanner

import (
	"fmt"
	"strconv"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/errorhandling"
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

var keywords = map[string]token.TokenType{
	"var": token.VAR,

	"and":   token.AND,
	"or":    token.OR,
	"if":    token.IF,
	"else":  token.ELSE,
	"true":  token.TRUE,
	"false": token.FALSE,
	"while": token.WHILE,
	"for":   token.FOR,

	"fun":    token.FUN,
	"return": token.RETURN,

	"class": token.CLASS,
	"this":  token.THIS,
	"super": token.SUPER,

	"print": token.PRINT,
	"nil":   token.NIL,
}

type Scanner struct {
	source   string
	tokens   []token.Token
	start    int
	current  int
	line     int
	keywords *(map[string]token.TokenType)
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:   source,
		tokens:   make([]token.Token, 0),
		start:    0,
		current:  0,
		line:     1,
		keywords: &keywords,
	}
}

func (s *Scanner) GetTokens() []token.Token {
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()

	}

	s.tokens = append(s.tokens, *token.NewToken(token.EOF, "", nil, s.line))
}

func (s *Scanner) skipWhiteSpace() byte {
	cc := s.advance()
	for cc == ' ' || cc == '\t' || cc == '\r' {
		s.start = s.current
		cc = s.advance()
	}

	return cc
}

func (s *Scanner) scanToken() {
	// c := s.advance()
	c := s.skipWhiteSpace()

	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)

	case ')':
		s.addToken(token.RIGHT_PAREN, nil)

	case '{':
		s.addToken(token.LEFT_BRACE, nil)

	case '}':
		s.addToken(token.RIGHT_BRACE, nil)

	case ',':
		s.addToken(token.COMMA, nil)

	case '.':
		s.addToken(token.DOT, nil)

	case '+':
		s.addToken(token.PLUS, nil)
	case '-':
		s.addToken(token.MINUS, nil)

	case ';':
		s.addToken(token.SEMICOLON, nil)

	case '*':
		s.addToken(token.STAR, nil)

	case '!':
		if s.match('=') {

			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case '=':
		if s.match('=') {

			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}

	case '<':
		if s.match('=') {

			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}

	case '>':
		if s.match('=') {

			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}

	case '\n':
		s.line++

	case '"':
		s.handleString()

	default:
		if isDigit(c) {
			for isDigit(s.peek()) {
				s.advance()
			}

			if s.peek() == '.' && isDigit(s.peekNext()) {
				s.advance()

				for isDigit(s.peek()) {
					s.advance()
				}
			}
			number := (s.source[s.start:s.current])
			f, err := strconv.ParseFloat(number, 64)
			if err != nil {

			}

			nt := token.NewToken(token.NUMBER, number, f, s.line)
			s.tokens = append(s.tokens, *nt)

		} else if IsAlpha(c) {

			for isAlphaNumeric(s.peek()) {
				s.advance()
			}
			text := s.source[s.start:s.current]
			t, ok := (*s.keywords)[text]
			if !ok {
				ident := (s.source[s.start:s.current])
				s.tokens = append(s.tokens, *token.NewToken(token.IDENTIFIER, ident, "", s.line))
			} else {

				s.tokens = append(s.tokens, *token.NewToken(t, text, "", s.line))

			}
		} else {
			fmt.Errorf("unknown character '%v' at line %d", string(c), s.line)
		}

	}

}

func (s *Scanner) advance() byte {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) addToken(t token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *token.NewToken(t, text, literal, s.line))

}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) charAt(index int) byte {
	return s.source[index]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	} else {
		return s.charAt(s.current)
	}
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	} else {
		return s.charAt(s.current + 1)
	}
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errorhandling.ReportAndExit(s.line, "", "unterminated string")
		return
	}
	s.advance()

	value := s.source[(s.start + 1):(s.current - 1)]
	// fmt.Println(value)
	s.addToken(token.STRING, value)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b == '_')
}

func isAlphaNumeric(b byte) bool {
	return isDigit(b) || IsAlpha(b)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	kw, ok := (*s.keywords)[text]
	if !ok {
		s.addToken(token.IDENTIFIER, nil)
	} else {
		s.addToken(kw, nil)
	}

	s.addToken(token.IDENTIFIER, nil)
}

func (s *Scanner) handleNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	f, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		errorhandling.ReportAndExit(s.line, "", "tried to treat a non-number as a number")
	}
	s.addToken(token.NUMBER, f)
}
