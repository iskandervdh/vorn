/*
Vorn is a simple interpreted C-like scripting language.
It can be used to run .vorn files or as a REPL when run without arguments.

Usage:

	vorn [flags] [path/to/file]

The flags are:

	--tokens
		Print the tokens of the input.

	-a, --ast
		Print the AST of the input.

	-h, --help
		Print this help message.

	-v, --version
		Print the version of Vorn.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/iskandervdh/vorn/constants"
	"github.com/iskandervdh/vorn/evaluator"
	"github.com/iskandervdh/vorn/lexer"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/parser"
	"github.com/iskandervdh/vorn/repl"
	"github.com/iskandervdh/vorn/token"
	"github.com/iskandervdh/vorn/version"
)

func runProgram(in io.Reader, out io.Writer) {
	// Create a new environment for the program
	env := object.NewEnvironment()

	// Read the file into a buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	// Create a new lexer and parser
	l := lexer.New(buf.String())
	p := parser.New(l, constants.TRACE)
	// Parse the program
	program := p.ParseProgram()

	// If there are any errors, print them and exit
	if len(p.Errors()) != 0 {
		fmt.Println("Error parsing program")
		parser.PrintErrors(out, p.Errors())
		return
	}

	// Create a new evaluator and evaluate the program
	e := evaluator.New()
	evaluated := e.Eval(program, env)

	// If the evaluated object is nil, something went wrong
	if evaluated == nil {
		io.WriteString(out, "Something went wrong while evaluating the program, got nil.\n")
		os.Exit(2)
	}

	// If the evaluated object is an error, print the error and exit
	if evaluated.Type() == object.ERROR_OBJ {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")

		os.Exit(1)
	}
}

func handleTokens(in io.Reader) {
	// Read the file into a buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	// Create a new lexer
	l := lexer.New(buf.String())

	// Print the first token
	currentToken := l.NextToken()
	fmt.Printf("%s ", currentToken.Type)

	// Print the rest of the tokens
	for token.TokenType(currentToken.Type) != token.EOF {
		if token.TokenType(currentToken.Type) == token.SEMICOLON {
			fmt.Println()
		}

		currentToken = l.NextToken()

		fmt.Printf("%s ", currentToken.Type)
	}

	fmt.Println()
}

func handleAST(in io.Reader, out io.Writer) {
	// Read the file into a buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	// Create a new lexer and parser
	l := lexer.New(buf.String())
	p := parser.New(l, constants.TRACE)
	// Parse the program
	program := p.ParseProgram()

	// If there are any errors, print them and exit
	if len(p.Errors()) != 0 {
		fmt.Println("Error parsing program")
		parser.PrintErrors(out, p.Errors())
		return
	}

	// Print the AST
	fmt.Println(program.String())
}

func printHelp() {
	fmt.Println(`Vorn is a simple interpreted C-like scripting language.
It can be used to run .vorn files or as a REPL when run without arguments.

Usage:

	vorn [flags] [path/to/file]

The flags are:

	--tokens
	    Print the tokens of the input.

	-a, --ast
	    Print the AST of the input.

	-v, --version
	    Print the version of Vorn.

	-h, --help
	    Print this help message.`)
}

func handleFlag(f string) {
	switch f {
	case "-h", "--help":
		printHelp()
		os.Exit(0)
	case "--tokens":
		if len(os.Args) < 3 {
			fmt.Println("Usage: vorn --tokens [file]")
			os.Exit(1)
		}

		file, err := os.OpenFile(os.Args[2], os.O_RDONLY, 0644)

		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", os.Args[2], err)
			os.Exit(1)
		}

		handleTokens(file)
		os.Exit(0)
	case "-a", "--ast":
		if len(os.Args) < 3 {
			fmt.Println("Usage: vorn -a [file]")
			os.Exit(1)
		}

		file, err := os.OpenFile(os.Args[2], os.O_RDONLY, 0644)

		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", os.Args[2], err)
			os.Exit(1)
		}

		handleAST(file, os.Stdout)
		os.Exit(0)
	case "-v", "--version":
		fmt.Printf("vorn %s\n", version.Version)
		os.Exit(0)
	default:
		fmt.Printf("Unknown flag %s\n", f)
		os.Exit(1)
	}
}

var traceFlag string

func init() {
	flag.Usage = printHelp
	flag.String("tokens", "", "Print the tokens of the input.")

	const (
		astDefault = ""
		astUsage   = "Print the AST of the input."
	)
	flag.StringVar(&traceFlag, "a", astDefault, astUsage+" (shorthand)")
	flag.StringVar(&traceFlag, "ast", astDefault, astUsage)

	const versionUsage = "Print the version of Vorn."
	flag.Bool("v", false, versionUsage+" (shorthand)")
	flag.Bool("version", false, versionUsage)

	const helpUsage = "Print this help message."
	flag.Bool("h", false, helpUsage+" (shorthand)")
	flag.Bool("help", false, helpUsage)
}

func main() {
	flag.Parse()

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
