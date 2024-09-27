package main

import (
	"bufio"
	"fmt"
	"lispy/lexer"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, _ := reader.ReadString('\n')
		_, items := lexer.Lex("REPL", line)

		for item := range items {
			fmt.Printf("%v\n", item)
		}
	}
}
