package inter

import (
	"testing"
)

func TestPContextSet(t *testing.T) {
	b := NewPContext()

	n1, n2, n3 := new(Node), new(Node), new(Node)
	n1.name, n2.name, n3.name = "1", "2", "3"

	if b.Get(EQ, n1) != nil {
		t.FailNow()
	}
	if b.Get(EQ, n2) != nil {
		t.FailNow()
	}
	if b.Get(EQ, n3) != nil {
		t.FailNow()
	}

	b1 := b.Set(EQ, n1, n2)
	if b.Get(EQ, n1) == n2 {
		t.FailNow()
	}
	if b1.Get(EQ, n1) != n2 {
		t.FailNow()
	}

	b2 := b1.Set(EQ, n1, n3)
	if b.Get(EQ, n1) == n2 {
		t.FailNow()
	}
	if b1.Get(EQ, n1) != n2 {
		t.FailNow()
	}
	if b2.Get(EQ, n1) != n3 {
		t.FailNow()
	}
	b1.dump()
	b2.dump()

}

func TestPContextUnset(t *testing.T) {
	b := NewPContext()
	n1, n2, n3 := new(Node), new(Node), new(Node)
	n1.name, n2.name, n3.name = "1", "2", "3"

	// set in bb
	bb := b.Set(EQ, n1, n2)
	if bb.Get(EQ, n1) == nil {
		t.FailNow()
	}

	// unset in cc
	cc := bb.Set(EQ, n1, nil)
	if cc.Get(EQ, n1) != nil {
		t.FailNow()
	}

	// but still there in bb
	if bb.Get(EQ, n1) == nil {
		t.FailNow()
	}
	b.dump()
	bb.dump()
	cc.dump()
}
