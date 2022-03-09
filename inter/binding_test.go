package inter

import "testing"

func TestBinding(t *testing.T) {
	b := NewBinder()
	n1, n2, n3, n4 := new(Node), new(Node), new(Node), new(Node)

	if b.Get(n1) != nil {
		t.FailNow()
	}
	if b.Get(n2) != nil {
		t.FailNow()
	}
	if b.Get(n3) != nil {
		t.FailNow()
	}
	if b.Get(n4) != nil {
		t.FailNow()
	}
	err := b.Pop()
	if err == nil {
		t.Fatal("Poping new binder should fail")
	}
	b.Set(n1, n2)
	if b.Get(n1) != n2 {
		t.FailNow()
	}
	b.Push()
	if b.Get(n1) != n2 {
		t.FailNow()
	}
	b.Set(n3, n4)
	if b.Get(n1) != n2 {
		t.FailNow()
	}
	if b.Get(n3) != n4 {
		t.FailNow()
	}
	err = b.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if b.Get(n1) != n2 {
		t.FailNow()
	}
	if b.Get(n3) == n4 {
		t.FailNow()
	}
}
