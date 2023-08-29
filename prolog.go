package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xavier268/myprolog/parser"
	"github.com/xavier268/myprolog/solver"
)

// main repl
func main() {
	fmt.Println("Welcome to prolog !")
	st := solver.NewState(nil) // create initial state
	for {
		fmt.Println("Enter multiple lines, empty line start executing.")
		input := readLinesUntilEmptyLine()

		terms, err := parser.ParseString(input, "keyboard")
		if err != nil {
			fmt.Println(err)
			// break
		}
		fmt.Println("Ok.")

		fmt.Println("\nterms = ", terms)

		st.Goals = append(st.Goals, terms...)

		st = solver.Solve(st, solHandlr)
		if st == nil {
			fmt.Println("Bye !")
			return
		}
	}
}

// naive (debugging) solution handler
func solHandlr(st *solver.State) *solver.State {
	var c rune
	if st == nil {
		fmt.Println("No more solution.")
		return nil
	}

	if len(st.Constraints) == 0 {
		fmt.Println("\nNo constraints.")
	} else {
		fmt.Println("Solution found:", st.Constraints)
	}
	fmt.Println("Enter 'x' to exit, 'd' for debug, 'a' to abort, otherwise to continue")
	fmt.Scanf("%c", &c)
	switch c {
	case 'x':
		return nil
	case 'd':
		fmt.Printf("\nDEBUG\n%#v\n", st)
		return st
	case 'a':
		st.Abort()
		return st
	default:
		return st.Parent
	}

}

// basic interactive entry.
func readLinesUntilEmptyLine() string {
	fmt.Println(">")
	sb := new(strings.Builder)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		line := scanner.Text()

		if line == "" {
			break
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	return sb.String()
}
