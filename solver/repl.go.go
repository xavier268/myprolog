package solver

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xavier268/myprolog/parser"
	"golang.org/x/term"
)

// Solution Handler would display solutions, if any.
// It should be a function that takes a state and returns a state.
// It tmodifies the state according to the the next expected task to do : look fore more solutions, stop, ...
var _ SolutionHandler = solHandlr

// main repl
func Repl() {
	st := NewState(nil) // create initial, non nil, state
	for {

		if st == nil {
			fmt.Println("There are no solutions")
			st = NewState(nil) // recreate empty state
		}

		st = Solve(st, solHandlr)
	}
}

// Accept new rule or query as input, add input as goals.
func AcceptRuleQuery(st *State) *State {

	if st == nil {
		return nil
	}

	fmt.Println("Enter new rules or queries, terminated by empty line :")
	input := readLinesUntilEmptyLine()
	terms, err := parser.ParseString(input, "keyboard entry")
	if err != nil {
		fmt.Println(err)
	}
	st.Goals = append(st.Goals, terms...)
	return st
}

// naive (debugging) solution handler
func solHandlr(st *State) *State {

	if st == nil {
		fmt.Println("No (more) solutions.")
		return st
	} else {
		fmt.Println("Ok.")
	}

	if len(st.Constraints) == 0 {
		// fmt.Println("\nNo constraints.")
	} else {
		fmt.Println("Solution :", FilterSolutions(st.Constraints))
	}
	fmt.Println("Enter 'x' to exit, 's' to dump state, 'n' for next solution, 'e' to enter new query or rules, 'r' to see the rules used")
	for {
		switch ReadCharRawMode() {
		case 's': // print states
			for ss := st; ss != nil; ss = ss.Parent {
				fmt.Println(ss)
			}
			return st
		case 'n':
			return st.Parent
		case 'e': // enter new query or rules
			return AcceptRuleQuery(st)
		case 'x':
			os.Exit(0)
		case 'r':
			fmt.Printf("Rules used :\n%s\n", st.RulesHistory())
			// loop
		case 0:
			os.Exit(1)
		default: // loop
		}

	}
}

// basic interactive entry.
func readLinesUntilEmptyLine() string {

	sb := new(strings.Builder)
	scanner := bufio.NewScanner(os.Stdin)
	line := 0
	for {
		line = line + 1
		fmt.Printf("%d>", line)
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

// Read a byte in raw mode.
func ReadCharRawMode() rune {
	b := make([]byte, 1)
	oldstate, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer term.Restore(int(os.Stdin.Fd()), oldstate)

	_, err = os.Stdin.Read(b)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return rune(b[0])
}
