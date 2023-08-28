package solver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavier268/myprolog/mytest"
)

// Data set of constraints to test
var cData = []Constraint{

	// Testing all zero values
	VarEQ{},
	VarIsAtom{},
	VarIsString{},
	VarIsVar{},
	// VarIsCompoundTerm{},

	VarEQ{
		V:     Variable{Name: "X", Nsp: 0},
		Value: Number{Num: 5, Den: 2, Normalized: false},
	},

	VarLT{
		V:     Variable{Name: "X", Nsp: 0},
		Value: Number{Num: 5, Den: 2, Normalized: false},
	},
	VarGT{
		V:     Variable{Name: "X", Nsp: 0},
		Value: Number{Num: 5, Den: 2, Normalized: false},
	},
	VarLTE{
		V:     Variable{Name: "X", Nsp: 0},
		Value: Number{Num: 5, Den: 2, Normalized: false},
	},
	VarGTE{
		V:     Variable{Name: "X", Nsp: 0},
		Value: Number{Num: 5, Den: 2, Normalized: false},
	},
	VarINT{
		V: Variable{Name: "X", Nsp: 0},
	},

	VarIsAtom{
		V: Variable{Name: "X", Nsp: 2},
		A: Atom{
			Value: "foo",
		},
	},

	VarIsString{
		V: Variable{
			Name: "X",
			Nsp:  2,
		},
		S: String{Value: "hello world"},
	},

	VarIsVar{
		V: Variable{
			Name: "X",
			Nsp:  2,
		},
		W: Variable{
			Name: "X",
			Nsp:  2,
		},
	},

	VarIsVar{ // The checked version will show a different order
		V: Variable{
			Name: "X",
			Nsp:  2,
		},
		W: Variable{
			Name: "Z",
			Nsp:  2,
		},
	},

	VarIsVar{
		V: Variable{
			Name: "Z",
			Nsp:  2,
		},
		W: Variable{
			Name: "X",
			Nsp:  2,
		},
	},
	VarIsVar{ // The checked version will show a different order
		V: Variable{
			Name: "X",
			Nsp:  2,
		},
		W: Variable{
			Name: "Z",
			Nsp:  1,
		},
	},

	VarIsVar{
		V: Variable{
			Name: "Z",
			Nsp:  1,
		},
		W: Variable{
			Name: "X",
			Nsp:  2,
		},
	},
}

func TestConstraintsCheck(t *testing.T) {
	var err error
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "\n========= Single constraint test =========\n")

	for i, c := range cData {

		fmt.Fprintln(sb)

		fmt.Fprintf(sb, "\n%d\t(original)\t%v", i, c.String())
		fmt.Fprintf(sb, "\n%d\t(raw form)\t%#v", i, c)

		c, err = c.Check()
		fmt.Fprintf(sb, "\n%d\t(checked)\t%v", i, c)
		if err != nil {
			fmt.Fprintf(sb, ", error : %v", err)
		}
	}

	mytest.Verify(t, sb.String(), "constraint_test")
}
