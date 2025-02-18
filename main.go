package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/iskandervdh/vorn/evaluator"
	"github.com/iskandervdh/vorn/lexer"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/parser"
	"github.com/iskandervdh/vorn/repl"
	"github.com/iskandervdh/vorn/version"
)

func runProgram(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	l := lexer.New(buf.String())
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Error parsing program")
		parser.PrintErrors(out, p.Errors())
		return
	}

	e := evaluator.New()
	evaluated := e.Eval(program, env)

	if evaluated.Type() == object.ERROR_OBJ {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")

		os.Exit(1)
	}
}

func handleFlag(flag string) {
	switch flag {
	case "-v", "--version":
		fmt.Printf("vorn %s\n", version.Version)
		os.Exit(0)
	case "-h", "--help":
		fmt.Println("Usage: vorn [file]")
		os.Exit(0)
	default:
		fmt.Printf("Unknown flag %s\n", flag)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("vorn %s\n", version.Version)
		repl.Start(os.Stdin, os.Stdout)

		return
	}

	if strings.HasPrefix(os.Args[1], "-") {
		handleFlag(os.Args[1])

		return
	}

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", os.Args[1], err)
		return
	}

	runProgram(file, os.Stdout)
}
