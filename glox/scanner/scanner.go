package scanner

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]token.Token, 10),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken(){
	byte c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
		break
	case ')'
		s.addToken(token.RIGHT_PAREN)
		break
	case '{':
		s.addToken(token.LEFT_BRACE)
		break
	case '}':
		s.addToken(token.RIGHT_BRACE)
		break
	case ',':
		s.addToken(token.COMMA)
		break
	case '.':
		s.addToken(token.DOT)
		break
	case '+':
		s.addToken(token.PLUS)
		break
	case ';':
		s.addToken(token.SEMICOLON)
		break
	case '*':
		s.addToken(token.STAR)
		break
	}
}

func (s *Scanner) advance() byte {
	return s.source[s.current++]
}

func (s *Scanner) addToken(type string, literal any){
	text := s.source[s.start, s.current]
	tokens := append(tokens, token.NewToken(type, text, literal, s.line))
}
