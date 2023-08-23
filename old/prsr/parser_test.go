package prsr

import (
	"fmt"
	"testing"

	"github.com/xavier268/myprolog/config"
	"github.com/xavier268/myprolog/tknz"
)

func TestParser(t *testing.T) {

	parserData := []struct {
		input  string
		err    bool
		output string
	}{
		{"toto xavier .", false, "[toto xavier]"},
		{"123 .", false, "[123]"},
		{"123 2.", true, "[123 2]"}, // period captured by number !
		{"123 2 .", false, "[123 2]"},
		{"123.", true, "[123]"},

		// should fail
		{"toto xavier", true, "[toto xavier]"}, // missing final period
		{"123", true, "[123]"},                 // missing final period
	}

	for _, pd := range parserData {

		tkn := tknz.NewTokenizerString(config.New(), pd.input)
		tt, err := ParseTerm(tkn)
		if (err != nil) != pd.err {
			t.Fatalf("unexpected error %v for parserData %v", err, pd)
		}
		out := fmt.Sprintf("%v", tt)
		if out != pd.output {
			t.Fatalf("unexpected output %v for parserData %v", out, pd)

		}

	}

}
