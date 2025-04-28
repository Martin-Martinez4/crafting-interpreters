package token

const (
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "["
	RIGHT_BRACE = "]"
	COMMA       = ","
	DOT         = "."
	MINUS       = "-"
	PLUS        = "+"
	SEMICOLON   = ";"
	SLASH       = "/"
	STAR        = "*"

	BANG          = "!"
	BANG_EQUAL    = "!="
	EQUAL         = "="
	EQUAL_EQUAL   = "=="
	GREATER       = ">"
	GREATER_EQUAL = ">="
	LESS          = "<"
	LESS_EQUAL    = "<="

	IDENTIFIER = "IDENT"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// Keywords
	FUN    = "FUN"
	VAR    = "VAR"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	AND    = "AND"
	CLASS  = "CLASS"
	FOR    = "FOR"
	OR     = "OR"
	PRINT  = "PRINT"
	SUPER  = "SUPER"
	THIS   = "THIS"
	WHILE  = "WHILE"

	EOF = "EOF"

	// may have to change to NULL
	NIL = "NIL"
)
