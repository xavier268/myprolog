package node

import (
	"fmt"
	"testing"
)

func TestNilNode(t *testing.T) {
	var n *Node
	if n.String() != " nil" {
		fmt.Println(n)
		t.Fatal("failed nil test")
	}
	n.Add()
	if n != nil {
		t.Fatal("n should still be nil !")
	}
}
func TestNodeNewCloneEquallStringDump(t *testing.T) {

	data := []string{ // creation-name, String()

		"lkj", " lkj",
		"", " nil",
		"123", " 123",
		"123.4", " 123.4",
		"-56.4", " -56.4",
		"-.4", " -0.4",
		"Xkjh", " Xkjh",
		"_", " _",
		"_Mlkj55", " _Mlkj55", // Not an Underscore, nor a Variable, but a String !
		"m55", " m55",
		"55m", " 55m",
		"$rule", " $rule",
		"rule", " rule",
	}

	var last *Node
	for i := 0; i+1 < len(data); i += 2 {

		nd := NewNode(data[i])
		if nd.String() != data[i+1] {
			fmt.Printf("%d : %q\n", i, data[i:i+2])
			t.Fatalf("unexpected String() value <%s>", nd.String())
		}

		// Verify a clone is equal
		n2 := nd.Clone()
		if !nd.Equal(n2) || !n2.Equal(nd) {
			fmt.Printf("%d : %q\nnd:%s\nn2:%s\n", i, data[i:i+2], nd, n2)
			t.Fatal("failed clone should equal test")
		}

		// Verify not equal to last try ?
		if nd.Equal(last) || last.Equal(nd) {
			fmt.Printf("%d : %q\nnd:%s\nn2:%s\n", i, data[i:i+2], nd, n2)
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
		t.Fatalf("Failed equal of clone\n   n1: %s\nclone:%s\n", n1, n1.Clone())
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
	n1.Walk(func(nn *Node) error {
		*count = *count + 1
		fmt.Println(nn.GetLoad())
		return nil
	})
	fmt.Println(n1)
	if *count != 7 { // f should be called also on nil node , hence 7 and not 6
		fmt.Println(n1)
		t.Fatalf("Unexpected node count : %d", *count)
	}
}
