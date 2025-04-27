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
	s.ScanTokens()

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
func TestSimpleTokenInputSkipWhiteSpace(t *testing.T) {
	input := `= +   ()
	{}	,	;`

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
	s.ScanTokens()

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
func TestDoubleCharInputs(t *testing.T) {
	input := `<===!=(){},	;`

	tests := []struct {
		name           string
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{"assign/equal sign", token.LESS_EQUAL, "<="},
		{"plus", token.EQUAL_EQUAL, "=="},
		{"plus", token.BANG_EQUAL, "!="},
		{"Left Paren", token.LEFT_PAREN, "("},
		{"Right Paren", token.RIGHT_PAREN, ")"},
		{"Left Brace", token.LEFT_BRACE, "{"},
		{"Right Brace", token.RIGHT_BRACE, "}"},
		{"Comma", token.COMMA, ","},
		{"Semicolon", token.SEMICOLON, ";"},
		{"End of File", token.EOF, ""},
	}

	s := NewScanner(input)
	s.ScanTokens()

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

func TestDoubleCharInputsSkipWhiteSpace(t *testing.T) {
	input := `<= ==
	!=		(){},	;`

	tests := []struct {
		name           string
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{"assign/equal sign", token.LESS_EQUAL, "<="},
		{"plus", token.EQUAL_EQUAL, "=="},
		{"plus", token.BANG_EQUAL, "!="},
		{"Left Paren", token.LEFT_PAREN, "("},
		{"Right Paren", token.RIGHT_PAREN, ")"},
		{"Left Brace", token.LEFT_BRACE, "{"},
		{"Right Brace", token.RIGHT_BRACE, "}"},
		{"Comma", token.COMMA, ","},
		{"Semicolon", token.SEMICOLON, ";"},
		{"End of File", token.EOF, ""},
	}

	s := NewScanner(input)
	s.ScanTokens()

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

func TestSkipOneLineComment(t *testing.T) {
	input := `// test
	=  // test
	+ / ()`

	tests := []struct {
		name           string
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{"assign/equal sign", token.EQUAL, "="},
		{"plus sign", token.PLUS, "+"},
		{"slash", token.SLASH, "/"},
		{"Left Paren", token.LEFT_PAREN, "("},
		{"Right Paren", token.RIGHT_PAREN, ")"},
	}

	s := NewScanner(input)
	s.ScanTokens()

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

func TestStringLiterals(t *testing.T) {
	input := `"Hello, world!"`

	tests := []struct {
		name            string
		expectedType    token.TokenType
		expectedLexeme  string
		expectedLiteral any
	}{
		{"string: Hello, world!", token.STRING, `"Hello, world!"`, "Hello, world!"},
	}

	s := NewScanner(input)
	s.ScanTokens()

	for i, tt := range tests {
		tok := s.tokens[i]

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {

			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Literal != tt.expectedLiteral {

			t.Fatalf("tests[%d, %s] - Literal wrong, expected: %q got : %q", i, tt.name, tt.expectedLiteral, tok.Literal)
		}
	}

	input = `"+=(){}"`

	tests = []struct {
		name            string
		expectedType    token.TokenType
		expectedLexeme  string
		expectedLiteral any
	}{
		{"string: +=(){}", token.STRING, "\"+=(){}\"", "+=(){}"},
	}

	s = NewScanner(input)
	s.ScanTokens()

	for i, tt := range tests {
		tok := s.tokens[i]

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {

			t.Fatalf("tests[%d, %s] - Lexeme wrong, expected: %q got : %q", i, tt.name, tt.expectedLexeme, tok.Lexeme)
		}
		if tok.Literal != tt.expectedLiteral {

			t.Fatalf("tests[%d, %s] - Literal wrong, expected: %q got : %q", i, tt.name, tt.expectedLiteral, tok.Literal)
		}
	}
}
