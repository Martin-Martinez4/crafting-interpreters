package scanner

import (
	"fmt"
	"testing"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

func TestSimpleTokenInput(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		name           string
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{"assign/equal sign", token.EQUAL, "="},
		{"plus", token.PLUS, "+"},
		{"Left Paren", token.LEFT_PAREN, "("},
		{"Right Paren", token.RIGHT_PAREN, ")"},
		{"Left Brace", token.LEFT_BRACE, "{"},
		{"Right Brace", token.RIGHT_BRACE, "}"},
		{"Comma", token.COMMA, ","},
		{"Semicolon", token.SEMICOLON, ";"},
		{"End of File", token.EOF, ""},
	}

	s := NewScanner(input)
	s.scanTokens()

	for i, tt := range tests {
		tok := s.tokens[i]

		fmt.Println(i, tok.Lexeme)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {

			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedLexeme, tok.Lexeme)
		}
	}
}
