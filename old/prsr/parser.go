package yacc

import (
	"fmt"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/tknz"
)

var infix = map[string]bool{
	",":  true,
	";":  true,
	"-?": true,
	":-": true,
}

// Parse a Term as a list of Term.
// Term must be terminated by a period (.), which will not be parsed.
func ParseTerm(tknzr *tknz.Tokenizer) ([]config.Term, error) {

	var tt []config.Term
	var err error

	// Read the tokens into a slice of Term, until but excluding the first '.'.
	for {
		tk := tknzr.Next()
		if tk == "" {
			return tt, fmt.Errorf("unexpected end of input while parsing a Term : missing final period")
		}
		if tk == "." {
			break // continue to next processing step
		}
		// only handle Integer at this stage, no floats ...
		var i int
		_, err = fmt.Sscanf(tk, "%d", &i)
		if err == nil {
			tt = append(tt, &config.Integer{Value: i})
			continue
		}

		if len(tk) >= 2 && tk[0] == '"' && tk[len(tk)-1] == '"' {
			tt = append(tt, &config.String{Value: tk[1 : len(tk)-1]})
		}

		if tk == "_" {
			tt = append(tt, &config.Underscore{})
		}

		// by default, assume everything else is Atom at this point
		tt = append(tt, &config.Atom{Value: tk})
	}

	for tt, err, changed := processFunctors(tt); err == nil && changed; tt, err, changed = processFunctors(tt) {
		// Now, process functors and parenthesis, creating corresponding CompoundTerms
	}

	return tt, nil
}

// tt is modified in place to implement the compound terms with functors.
// idempotent.
func processFunctors(tt []config.Term) (newtt []config.Term, err error, changed bool) {

	// find outside indexes
	left, right := -1, -1 // first left and last right parenthesis
	for i, t := range tt {
		if a, ok := t.(*config.Atom); left < 0 && ok && a.Value == "(" { // first left parenthesis
			left = i
		}
		if a, ok := t.(*config.Atom); ok && a.Value == ")" { // last right parenthesis
			right = i
		}
	}
	// no parenthesis found
	if left == -1 && right == -1 {
		return tt, nil, false
	}
	// parenthesis do not match
	if left < 0 || right < 0 || left >= right {
		return tt, fmt.Errorf("parenthesis do not match"), false
	}
	// find functor
	if left <= 0 {
		return tt, fmt.Errorf("missing functor"), false
	}
	functor, ok := (tt[left-1]).(*config.Atom)
	if !ok {
		return tt[:left-1], fmt.Errorf("invalid functor"), false
	}

	// build the compound Term
	c := &config.CompoundTerm{
		Functor:  functor,
		Children: tt[left+1 : right],
	}
	tt[left] = c
	tt = append(tt[:left+1], tt[right+1:]...)
	return tt, nil, true
}

func processDisjonction(tt []config.Term) (newtt []config.Term, err error, changed bool) {

	where := []int{}

	for i, t := range tt {
		if a, ok := t.(*config.Atom); ok && a.Value == ";" {
			where = append(where, i)
		}
	}

	if len(where) == 0 {
		return tt, nil, false
	}

}
