package mytest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Verify content against reference file.
// If no reference file found, creates it.
func Verify(t *testing.T, content string, filename string) {
	filename, _ = filepath.Abs(filename)
	fmt.Println(GREEN+"Verifying test results against file : "+RESET, filename)
	check, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(RED + "File not found, create it as a reference for future test. Make sure you manually review it !" + RESET)
		os.WriteFile(filename, []byte(content), 0644)
		return
	}
	sc := string(check)
	if sc != content { // we know it is different, lets try to show where ...
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

				fmt.Printf("\n===============================================================\n%s : Results differ from reference file", t.Name())
				fmt.Printf("\n============================ got ==============================\n%s%s%s%s\n",
					content[i1:i], RED, content[i:i2], RESET)
				if i2 >= len([]rune(sc)) {
					i2 = len([]rune(sc))
				}
				fmt.Printf("\n============================ want==============================\n%s%s%s%s\n",
					sc[i1:i], RED, sc[i:i2], RESET)

				t.Fatalf("Result differs from reference file in %s", t.Name())
			}
		}

		// If we reach here, it means conet is matching ref, but both files are different.
		// There must be extra reference ? Let's show it.

		i1 := len(content) - 160
		if len(content) <= 160 {
			i1 = 0
		}
		i := len(content)
		i2 := len(content) + 160
		if i2 > len(sc) {
			i2 = len(sc)
		}

		fmt.Printf("\n===============================================================\n%s : Results differ from reference file", t.Name())
		fmt.Printf("\n============================ got ==============================\n\n")
		fmt.Println()
		fmt.Printf("\n============================ want==============================\n%s%s%s%s\n",
			sc[i1:i], RED, sc[i:i2], RESET)

		t.Fatalf("Result differs from reference file in %s", t.Name())

	}
}
