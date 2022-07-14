package node

import (
	"fmt"
	"testing"
)

func TestNilNode(t *testing.T) {
	var n *Node
	if n.String() != "" || n.Dump() != "<Nil>" {
		t.Fatal("failed nil test")
	}
	n.Add()
	n.Add(NewNode("toto")) // will be ignored, since n is nil
	fmt.Println(n.Dump())
	if n != nil {
		t.Fatal("n should still be nil !")
	}
}
func TestNodeNewCloneEquallStringDump(t *testing.T) {

	data := []string{ // creation-name, String(), Dump()

		"lkj", "lkj", "<Str>lkj",
		"", "", "<Nil>",
		"123", "123", "<Num>123",
		"123.4", "123.4", "<Num>123.4",
		"-56.4", "-56.4", "<Num>-56.4",
		"-.4", "-0.4", "<Num>-0.4",
		"Xkjh", "Xkjh", "<Var>Xkjh",
		"_", "_", "<Uds>_",
		"_Mlkj55", "_Mlkj55", "<Str>_Mlkj55", // Not an Underscore, nor a Variable, but a String !
		"m55", "m55", "<Str>m55",
		"55m", "55m", "<Str>55m",
		"$rule", "$rule", "<Str>$rule",
	}

	var last *Node
	for i := 0; i+2 < len(data); i += 3 {

		nd := NewNode(data[i])
		if nd.String() != data[i+1] {
			fmt.Printf("%d : %q\n", i, data[i:i+3])
			t.Fatal("unexpected String() value", nd.String())
		}
		if nd.Dump() != data[i+2] {
			fmt.Printf("%d : %q\n", i, data[i:i+3])
			t.Fatal("unexpected Dump() value", nd.Dump())
		}
		// Verify a clone is equal
		n2 := nd.Clone()
		if !nd.Equal(n2) || !n2.Equal(nd) {
			fmt.Printf("%d : %q\nnd:%s\nn2:%s\n", i, data[i:i+3], nd, n2)
			t.Fatal("failed clone should equal test")
		}

		// Verify not equal to last try ?
		if nd.Equal(last) || last.Equal(nd) {
			fmt.Printf("%d : %q\nnd:%s\nn2:%s\n", i, data[i:i+3], nd, n2)
			t.Fatal("failed should never equal last")
		}
		last = nd.Clone()
	}
}

func TestAdd(t *testing.T) {

	n1 := NewNode("n1")
	n2 := NewNode("n2")
	n3 := NewNode("n3")
	n4 := NewNode("n4")

	n1.Add(n2, n3)
	n3.Add(n4)
	n1.Add(n3)

	n11 := n1.Clone()
	if !n1.Equal(n11) {
		t.Fatalf("should be equal n11")
	}
	n4.Add(nil) // legal, will modify n1, not equal anymore to n11
	if n1.Equal(n11) {
		t.Fatalf("should not be equal, n11 after adding nil")
	}

	n12 := n1.Clone()
	if !n12.Equal(n1) || !n1.Equal(n12) {
		t.Fatalf("Failed equal of clone\n   n1: %s\nclone:%s\n", n1.Dump(), n1.Clone().Dump())
	}

}

func TestWalkVisual(t *testing.T) {

	n1 := NewNode("n1")
	n2 := NewNode("n2")
	n3 := NewNode("n3")
	n4 := NewNode("n4")

	n1.Add(n2, n3)
	n3.Add(n4)
	n1.Add(n3)
	n2.Add(nil)

	var count = new(int)
	n1.Walk(func(load any) error { *count = *count + 1; fmt.Println(load); return nil })
	if *count != 6 { // nil node is nevr called (no load), hence 6 and not 7
		fmt.Println(n1.Dump())
		t.Fatalf("Unexpected node count : %d", *count)
	}
}
