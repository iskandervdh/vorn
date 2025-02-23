package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/iskandervdh/vorn/constants"
	"github.com/iskandervdh/vorn/evaluator"
	"github.com/iskandervdh/vorn/lexer"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/parser"
)

// The REPL prompt to show the user.
const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	// Create a new scanner that reads from the input stream.
	scanner := bufio.NewScanner(in)
	// Create a new environment for the REPL that will persist between input lines.
	env := object.NewEnvironment()

	for {
		// Print the REPL prompt and wait for the user to enter a line.
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		// If the scanner encountered an error or the user entered an EOF character, stop the REPL.
		if !scanned {
			return
		}

		// Read the line that the user entered.
		line := scanner.Text()
		// Create a new lexer and parser for the input line.
		l := lexer.New(line)
		p := parser.New(l, constants.TRACE)
		// Parse the current line
		program := p.ParseProgram()

		// If the parser encountered any errors, print them and continue to the next line.
		if len(p.Errors()) != 0 {
			parser.PrintErrors(out, p.Errors())
			continue
		}

		// Otherwise evaluate the current line and print the result.
		e := evaluator.New()
		evaluated := e.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
