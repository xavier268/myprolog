package solver

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xavier268/myprolog/parser"
	"golang.org/x/term"
)

const helpRepl = "'x':exit, 's' status, 'n' next solution, 'e' or <space>: enter new query or rule, 'r' : explain rules, 'z': zero the rules, 'h' for help"

// convenience entry point, same as NewSession().Repl()
func Repl() {
	NewSession().Repl()
}

// main repl loop
func (sess *Session) Repl() {
	st := NewState(nil) // create initial, non nil, state
	sess = st.Session   // save session for future resets accross states
	fmt.Println(helpRepl)
	for {
		st = Solve(st, sess.interactiveSolutionHandler) // display all solutions, return nil on empty solutions.
		//fmt.Println("DEBUG - Solve returned with state : ", st)
		if st == nil {
			st = NewStateWithSession(sess) // re-create new state , keeping known rules.
		}
	}
}

// solution handler
func (sess *Session) interactiveSolutionHandler(st *State) *State {

	if st == nil {
		fmt.Println("No (more) solutions.")
		return NewStateWithSession(sess) // do not loose the rules already entered !
	} else {
		fmt.Println("Ok.")
	}

	if len(st.Constraints) == 0 {
		// fmt.Println("\nNo constraints.")
	} else {
		fmt.Println("Solution, if ", FilterSolutions(st.Constraints))
	}
	for {
		switch readCharRawMode() {
		case 's': // print states
			fmt.Print(st)
			fmt.Printf("\nKnown rules :\n%s-------------------\n", sess.ListRules())
			// loop
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
		case 'z': // reset both current state and database
			fmt.Printf("%sCAUTION : ALL RULES AND CONTEXT WILL BE DESTROYEDBE ZEROED%s\nPlease confirm that is what you really want ('y')\n", RED, RESET)
			if readCharRawMode() == 'y' {
				sess.ResetSession()
				st = nil
				fmt.Printf("%sALL RULES AND CONTEXT HAVE BEEN DESTROYED%s\n", RED, RESET)
				return nil
			}
			fmt.Println("Request was canceled ...")
			// else, loop ...
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
