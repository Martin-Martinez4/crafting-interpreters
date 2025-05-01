package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

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

	expression := parser.NewBinaryExpr(
		parser.NewUnaryExpr(minusOp, parser.NewLiteralExpr(123)),
		starOp,
		parser.NewGroupingExpr(parser.NewLiteralExpr(45.67)),
	)

	astp := parser.AstPrinter{}
	fmt.Println(astp.Print(expression))

	scanner := scanner.NewScanner("22 + 2 / (2 * 8) + 10 /4")
	scanner.ScanTokens()

	// p := parser.NewParser(scanner.GetTokens())
	// stmts := p.Parse()

	args := os.Args[1:]

	lenArgs := len(args)

	// fmt.Println(args)

	if lenArgs > 1 {
		println("Usage: glox [script]")
	} else if lenArgs == 1 {
		runFile(args[0])
	} else {
		runPrompt(os.Stdin, os.Stdout)
	}
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
	// astp := parser.AstPrinter{}
	scanner := scanner.NewScanner(source)
	scanner.ScanTokens()

	p := parser.NewParser(scanner.GetTokens())

	// for i := 0; i < len(scanner.GetTokens()); i++ {
	// 	fmt.Printf("line %d: '%s' %s\n", scanner.GetTokens()[i].Line, scanner.GetTokens()[i].Lexeme, scanner.GetTokens()[i].Type)
	// }

	stmts := p.Parse()

	i := parser.NewInterpreter(stmts)

	r := parser.NewResolver(i)
	r.ResolveStmts(stmts)

	i.Interpret(stmts)

}

func runPrompt(in io.Reader, out io.Writer) {
	PROMPT := "->"
	s := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := s.Scan()
		if !scanned {
			break
		}

		run(s.Text())
	}
}
