package repl

import (
	"bufio"
	"fmt"
	"io"
	"uBasic/errors"
	"uBasic/eval"
	"uBasic/lexer"
	"uBasic/parser"
	"uBasic/sem"
	"uBasic/source"
)

const PROMPT = ">> "
const PROMPT_CONTINUE = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// accumulate as many lines as the user wants to enter
		// until the user enters an empty line
		text := ""
		fmt.Fprint(out, PROMPT)
		for {
			scanned := scanner.Scan()
			if !scanned {
				break
			}
			text += scanner.Text()
			if scanner.Text() == "" {
				break
			}
			fmt.Fprint(out, PROMPT_CONTINUE)
		}

		l := lexer.New(text)
		p := parser.New(l)
		file := p.ParseFile()

		if file == nil {
			if p.Errors() != nil {
				for _, err := range p.Errors() {
					e := err.(*errors.Error)
					e.Src = &source.Source{Input: text, Name: "repl"}
					fmt.Fprintln(out, err)
				}
			}
		} else {
			info, err := sem.Check(file)
			if err != nil {
				e := err.(*errors.Error)
				e.Src = &source.Source{Input: text, Name: "repl"}
				fmt.Fprintln(out, err)
			} else {
				env := eval.Define(info)
				eval.Eval(nil, file, env)
			}
		}
	}
}
