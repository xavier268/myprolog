package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	var tdata = []string{

		// rules implicit
		"un(deux,trois).",
		"un.",
		"un(deux).",
		"un(2,3).",
		"empty().",
		"un(deux(),trois).",

		// rules explicit
		"a(X):-un(X).",
		"a(X):-un(X,Y).",
		"a(X):-un(X,Y,Z);b(Y);c(X).",
		"a(X):-.",

		// queries
		"?- 3.", // invalid syntax
		"?- test.",
		"?- un(deux,X).",
		"?- un(deux,X), trois(X).",

		// list
		"[].",
		"[2].",
		"[2,3].",
		"[2,3,4].",

		// sublist
		"[[2,3],4].",
		"[2,[3,4]].",

		// list and pairs
		"[2|3].",
		"[4|].",
		"[4|X].",

		// non list
		"dot(1,dot(2,3)).",                   // not a list
		"dot(1,dot(dot(4,dot(5,dot())),3)).", // not a list, but contains a list

		// invalid
		"un,deux.",
		"2(a).",
		" [|2].",
		"[|]",
		"a(b,,).",
	}

	res := run(tdata)

	verify(t, res, "parser_test_want.txt")

}

func run(tdata []string) string {
	sb := new(strings.Builder)
	for i, d := range tdata {

		r, err := ParseString(d, fmt.Sprintf("test # %d <%s>", i, d))
		fmt.Fprintf(sb, "\n%d \t\t<%s>\n", i, d)
		fmt.Fprintf(sb, "%d\t\terr=%v\n", i, err)
		for _, v := range r {
			fmt.Fprintf(sb, "%d\t\t(string)    %s\n", i, v.String())
			fmt.Fprintf(sb, "%d\t\t(pretty)    %s\n", i, v.Pretty())
		}
	}
	return sb.String()
}

func verify(t *testing.T, content string, filename string) {
	filename, _ = filepath.Abs(filename)
	fmt.Println("Verifying parse results against file : ", filename)
	check, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("File not found, create it as a reference for future test. Make sure you manually review it !")
		os.WriteFile(filename, []byte(content), 0644)
		return
	}
	sc := string(check)
	if sc != content {
		for i, c := range content {
			if i >= len([]rune(sc)) || c != ([]rune(sc)[i]) {
				i1 := i - 160
				i2 := i + 160
				if i1 <= 0 {
					i1 = 0
				}
				if i2 > len(content) {
					i2 = len(content)
				}
				fmt.Printf("Parser result differ from reference file\n")
				fmt.Printf("\n============================ got ==============================\n%s%s%s%s\n",
					content[i1:i], START_RED, content[i:i2], END_RED)
				if i2 >= len([]rune(sc)) {
					i2 = len([]rune(sc))
				}
				fmt.Printf("\n============================ want==============================\n%s%s%s%s\n",
					sc[i1:i], START_RED, sc[i:i2], END_RED)

				t.Fatalf("parser result differs from reference file")
			}
		}
	}
}
