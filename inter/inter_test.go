package inter

import "testing"

const p1 = `
fils(X,Y)			~	pere(Y,X).
fils(jean jacques).
fils(jean paul).
age(jean 30).
age(jacques 40).
parent(X,Y) ~ pere(X,Y);mere(X,Y).
testPrint 			~ print("hello world").
`

func TestParseRulesVisual(t *testing.T) {

	tzr := NewTokenizerString(p1)
	in := NewInter()
	err := in.ParseRules(tzr)
	if err != nil {
		in.Dump()
		t.Fatal(err)
	}

}
