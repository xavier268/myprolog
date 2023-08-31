package parser

import (
	"fmt"
	"io"
	"math"
	"strings"
	"text/scanner"
)

// Lexer for prolog

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type myLex struct {
	s       scanner.Scanner // golang scanner
	LastErr []error         // errors emitted
}

// Required to satisfy interface.
func (lx *myLex) Error(s string) {
	lx.LastErr = append(lx.LastErr, fmt.Errorf("error in %s, line %d : %v", lx.s.Filename, lx.s.Line, s))
	fmt.Println(lx.LastErr)
}

// The myLexer interface is implemented by myLex, defined by generated parser.
var _ myLexer = new(myLex)

// NewLexer from io.Reader
func NewLexer(input io.Reader, sourcename string) *myLex {
	ps := new(myLex)
	ps.s.Mode = scanner.GoTokens
	ps.s.Init(input)
	ps.s.Filename = sourcename
	return ps
}

// The parser calls this method to get each new token.
// This implementation returns token extracted from scanner
func (lx *myLex) Lex(lval *mySymType) int {

	lval.list = []Term{} // lexer never return a list
	lval.value = nil     // reset the value to nil

	tk := lx.s.Scan()
	switch tk {
	case scanner.EOF:
		return eof // 0 for the lexer
	case scanner.Ident:
		txt := lx.s.TokenText()
		if txt == "_" {
			lval.value = Underscore{}
			return '_'
		}
		if txt[0] >= 'A' && txt[0] <= 'Z' {
			lval.value = Variable{
				Name: txt,
				Nsp:  0,
			}
			return VARIABLE
		}
		if len(txt) > 1 && txt[0] == '_' {
			lx.Error(fmt.Sprintf("Identifier cannot begin with underscore : %v", txt))
			return eof
		}
		lval.value = Atom{
			Value: txt,
		}
		return ATOM

	case scanner.Comment: // ignore comments
		return lx.Lex(lval)

	case scanner.RawString, scanner.String:
		lval.value = String{
			Value: lx.s.TokenText()[1 : len(lx.s.TokenText())-1],
		}
		return STRING

	case scanner.Char: // handled as a STRING with 1 char.
		lval.value = String{
			Value: lx.s.TokenText()[1 : len(lx.s.TokenText())-1],
		}
		return STRING

	case scanner.Float: // float are only accepted in the %f format. They are converted to rationales, normalized.
		var num, den int
		var z float64
		txt := strings.ToUpper(lx.s.TokenText())
		if strings.Contains(txt, "E") {
			lx.Error(fmt.Sprintf("Exponents are not allowed in numbers : %v", txt))
			return LEXERROR
		}
		_, err := fmt.Sscanf(txt, "%f", &z)
		if err != nil {
			lx.Error(fmt.Sprintf("Expected an decimal number like %f but got %v instead", 10.455, txt))
			return LEXERROR
		}
		before, after, _ := strings.Cut(txt, ".")
		fmt.Sscanf(before+after, "%d", &num)
		pow := len(after) // always > 0
		den = int(math.Pow10(pow) + 0.001)
		lval.value = Number{
			Num:        num,
			Den:        den,
			Normalized: false,
		}.Normalize()
		return NUMBER

	case scanner.Int: // Rational CAN be represented as 4/3. We need to check for this.
		var num, den int
		txt := strings.ToUpper(lx.s.TokenText())
		if strings.Contains(txt, "E") {
			lx.Error(fmt.Sprintf("Exponents are not allowed in integers : %v", txt))
			return LEXERROR
		}
		fmt.Sscanf(txt, "%d", &num)
		if lx.s.Peek() == '/' {
			_ = lx.s.Scan()                 // eat /
			peek := lx.s.Peek()             // look for the denominator
			if peek >= '0' && peek <= '9' { // peek gets the next unicode, not the next token !
				lx.s.Scan() // get the denominator
				txt = strings.ToUpper(lx.s.TokenText())
				if strings.Contains(txt, "E") {
					lx.Error(fmt.Sprintf("Exponents are not allowed in integers : %v", txt))
					return LEXERROR
				}
				fmt.Sscanf(txt, "%d", &den)
				lval.value = Number{
					Num:        num,
					Den:        den,
					Normalized: false,
				}.Normalize()
				return NUMBER
			} else {
				lx.Error("Expected an number in the form of a rational, like a/b but got a/ instead")
				return LEXERROR
			}
		}
		// no / available, it is a normal integer
		lval.value = Number{
			Num:        num,
			Den:        1,
			Normalized: true,
		}
		return NUMBER

	case '(', ')', '[', ']', ',', ';', '.', '|', '-': // single char tokens recognized by parser, cannot begin a multichar operator.
		// yylval is not set for these.
		return int(tk)

	case ':': // possibly multichar token.
		if lx.s.Peek() == '-' {
			lx.s.Scan()
			return OPRULE
		}
		// a single ':' is recognized as an atom.
		lval.value = Atom{
			Value: ":",
		}
		return ATOM

	case '?': // possibly a query operator.
		if lx.s.Peek() == '-' {
			lx.s.Scan()
			return OPQUERY
		}
		// a single '?' is recognized as an atom.
		lval.value = Atom{
			Value: "?",
		}
		return ATOM

	default:
		lx.Error(fmt.Sprintf("Unknown token type %v :  %v", tk, lx.s.TokenText()))
		return LEXERROR
	}

}
