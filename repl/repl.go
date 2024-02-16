package repl

import (
	"bufio"
	"fmt"
	"io"
	"uBasic/lexer"
	"uBasic/parser"
	"uBasic/sem"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l)
		file := p.ParseFile()

		if file != nil {
			_, err := sem.Check(file)
			if err != nil {
				fmt.Fprintln(out, err)
			} else {
				fmt.Fprintln(out, file.String())
			}
		} else if p.Errors() != nil {
			for _, msg := range p.Errors() {
				fmt.Fprintln(out, msg)
			}
		}

	}
}
