package node

import "fmt"

// contains reserved keywords
var keywords = []string{
	"dot",
	"rule",
	"slash",
}

type Keyword struct {
	name string
}

func (kw *Keyword) String() string {
	if kw == nil {
		return " nil"
	}
	return kw.name
}

var errNotFound = fmt.Errorf("keyword not found")

func Reserved(kw string) (Keyword, error) {
	for _, k := range keywords {
		if k == kw {
			return Keyword{kw}, nil
		}
	}
	return Keyword{}, errNotFound
}
