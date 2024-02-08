package repl

import (
	"bufio"
	"fmt"
	"io"
	"uBasic/lexer"
	"uBasic/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	l := lexer.New("")
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		l = lexer.New(scanner.Text())
		for tok := l.NextToken(); tok.Kind != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
