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

const helpRepl = "'x':exit, 's' status, 'n' next solution, 'e' or <space>: enter new query or rule, 'r' : explain rules, 'z': zero the rules, 'h' for help"

// main repl
func Repl() {
	st := NewState(nil) // create initial, non nil, state
	fmt.Println(helpRepl)
	for {
		st = Solve(st, solHandlr) // display all solution, return nil on empty solutions.
		if st == nil {
			st = NewState(nil) // prevent exit, rules are kept
		}
	}
}

// solution handler
func solHandlr(st *State) *State {

	if st == nil {
		fmt.Println("No (more) solutions.")
		return NewState(nil)
	} else {
		fmt.Println("Ok.")
	}

	if len(st.Constraints) == 0 {
		// fmt.Println("\nNo constraints.")
	} else {
		fmt.Println("Solution :", FilterSolutions(st.Constraints))
	}
	for {
		switch readCharRawMode() {
		case 's': // print states
			fmt.Print(st)
			fmt.Printf("\nKnown rules :\n%s-------------------\n", ListDBRules())
			return st
		case 'n':
			return st.Parent

		case 'e', ' ': // enter new query or rules
			st = acceptRuleQuery(st)
			return st
		case 'x':
			os.Exit(0)
		case 'r':
			fmt.Printf("Rules used :\n%s\n", st.RulesHistory())
			// loop
		case 0:
			os.Exit(1)
		case 'z': // reset
			ResetDB()
			return nil
		case 'h':
			fmt.Println(helpRepl) // loop
		case '\n', '\r', '\t':
		// Ignoren
		// loop
		default:
			//fmt.Println(help) // loop
			// loop
		}

	}
}

// Accept and parse new rule or query as input from stdin, add input as goals.
func acceptRuleQuery(st *State) *State {

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
func readCharRawMode() rune {
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
