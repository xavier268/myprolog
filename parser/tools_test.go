package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Verify content against reference file.
// If no reference file found, creates it.
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
