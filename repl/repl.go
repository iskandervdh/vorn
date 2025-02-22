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

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l, constants.TRACE)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			parser.PrintErrors(out, p.Errors())
			continue
		}

		e := evaluator.New()
		evaluated := e.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
