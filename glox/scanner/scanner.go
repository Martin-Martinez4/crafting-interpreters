package scanner

import (
	"github.com/Martin-Martinez4/crafting-interpreters/glox/errorhandling"
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
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()

	}

	s.tokens = append(s.tokens, *token.NewToken(token.EOF, "", nil, s.line))
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
		break
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
		break
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
		break
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
		break
	case ',':
		s.addToken(token.COMMA, nil)
		break
	case '.':
		s.addToken(token.DOT, nil)
	case '=':
		s.addToken(token.EQUAL, nil)
		break
	case '+':
		s.addToken(token.PLUS, nil)
		break
	case ';':
		s.addToken(token.SEMICOLON, nil)
		break
	case '*':
		s.addToken(token.STAR, nil)
		break
	default:
		errorhandling.ReportAndExit(s.line, "", "unexpected character.")
		break
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
