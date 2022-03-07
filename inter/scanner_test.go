package inter

import (
	"fmt"
	"testing"
)

func TestScanner101(t *testing.T) {
	src := `
	hello world // comment 
	next sentence 2+3; 2x /* comment 
	kjh multiline */ x2 <= :-  "double 	quoted" `
	fmt.Printf("%q\n", src)
	tzr := NewTokenizerString(src)
	for tk := tzr.Next(); tk != ""; tk = tzr.Next() {
		fmt.Printf("\t>%s<\n", tk)
	}
}

func TestScannerFile(t *testing.T) {
	fn := "scanner.go"
	tzr := NewTokenizerFile(fn)
	tk := tzr.Next()
	fmt.Printf("Fisrt token from file %s is : %s\n", fn, tk)
	if tk != "package" {
		fmt.Println("It should have been 'package'")
		t.FailNow()
	}
}
func TestScanner(t *testing.T) {

	// Test definition
	table := [...]struct {
		input  string
		expect []string
	}{
		{"Hello world", []string{"Hello", "world"}},
		{"Hello \n   \nworld\n", []string{"Hello", "world"}},
		{"\nHello \n   \nworld\n", []string{"Hello", "world"}},
		{"\n\nHello \n   \nworld", []string{"Hello", "world"}},
		{"1+2+3", []string{"1", "+", "2", "+", "3"}},
		{"1+2+3// comment", []string{"1", "+", "2", "+", "3"}},
		{"// comment\n1+2+3// comment", []string{"1", "+", "2", "+", "3"}},
		{"1+/* comment */2+3// comment", []string{"1", "+", "2", "+", "3"}},
		{"1<=2", []string{"1", "<", "=", "2"}},
		{"1<=2", []string{"1", "<", "=", "2"}},
		{"1 <= 2", []string{"1", "<", "=", "2"}},
		{"1:-2", []string{"1", ":", "-", "2"}},
		{"1: -2", []string{"1", ":", "-", "2"}},
		{"1 :-2", []string{"1", ":", "-", "2"}},
		{"1:- 2", []string{"1", ":", "-", "2"}},
		{" un \"deux trois   \" quatre", []string{"un", "\"deux trois   \"", "quatre"}},
	}

	// Loop table and compare results to expectation.
	for _, ts := range table {
		tzr := NewTokenizerString(ts.input)
		var got []string
		for i, tk := 0, tzr.Next(); tk != ""; tk, i = tzr.Next(), i+1 {
			got = append(got, tk)
			if ts.expect[i] != tk {
				t.Errorf("%s does not match %s\n", ts.expect[i], tk)
				t.Fatalf("Expected :\n%v\nGot :\n%v\n", ts.expect, got)
			}

		}
		if len(ts.expect) != len(got) {
			t.Errorf("Length do not match : %d and %d\n", len(ts.expect), len(got))
			t.Fatalf("Expected :\n%v\nGot :\n%v\n", ts.expect, got)
		}
	}

}
