package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Martin-Martinez4/crafting-interpreters/glox/ast"
	"github.com/Martin-Martinez4/crafting-interpreters/glox/parser"
	"github.com/Martin-Martinez4/crafting-interpreters/glox/scanner"
	"github.com/Martin-Martinez4/crafting-interpreters/glox/token"
)

func main() {

	// expression := ast.Binary{
	// 	Left: ast.Unary{
	// 		Operator: *token.NewToken(token.MINUS, "-", nil, 1),
	// 		Right: ast.Literal{
	// 			Value: 123,
	// 		},
	// 	},
	// 	Operator: *token.NewToken(token.STAR, "*", nil, 1),
	// 	Right: ast.Grouping{
	// 		Expression: ast.Literal{Value: 45.67},
	// 	},
	// }

	minusOp := token.NewToken(token.MINUS, "-", nil, 1)
	starOp := token.NewToken(token.STAR, "*", nil, 1)

	expression := ast.NewBinaryExpr(
		ast.NewUnaryExpr(minusOp, ast.NewLiteralExpr(123)),
		starOp,
		ast.NewGroupingExpr(ast.NewLiteralExpr(45.67)),
	)

	astp := ast.AstPrinter{}
	fmt.Println(astp.Print(expression))

	scanner := scanner.NewScanner("22 + 2 / (2 * 8) + 10 /4")
	scanner.ScanTokens()

	p := parser.NewParser(scanner.GetTokens())
	exprs := p.Parse()

	fmt.Println(astp.Print(exprs))

	// args := os.Args[1:]

	// lenArgs := len(args)

	// if lenArgs > 1 {
	// 	println("Usage: glox [script]")
	// } else if lenArgs == 1 {
	// 	runFile(args[1])
	// } else {
	// 	runPrompt(os.Stdin, os.Stdout)
	// }
}

func runFile(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	run(string(f))

	return nil
}

func run(source string) {
	fmt.Println("Coming Soon")
}

func runPrompt(in io.Reader, out io.Writer) {
	PROMPT := "->"
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		run(scanner.Text())
	}
}
