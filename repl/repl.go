package repl

import (
	"bufio"
	"fmt"
	"github.com/Neal-C/interpreter-in-go/lexer"
	"github.com/Neal-C/interpreter-in-go/token"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		monkeyLexer := lexer.New(line)

		for tok := monkeyLexer.NextToken(); tok.Type != token.EOF; tok = monkeyLexer.NextToken() {
			fmt.Printf("%+v \n", tok)
		}
	}

}
