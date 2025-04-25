package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	args := os.Args[1:]

	lenArgs := len(args)

	if lenArgs > 1 {
		println("Usage: glox [script]")
	} else if lenArgs == 1 {
		runFile(args[1])
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
