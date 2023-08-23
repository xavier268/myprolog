package parser

import (
	"fmt"
)

func ExampleParseString_valid1() {

	tdata := []string{
		// valid 1
		"un(deux,trois).",
		"un.",
		"un(deux).",
		"un(2,3).",
		"empty().",
		"un(deux(),trois).",
	}

	runex(tdata)

	// Output:
	// un(deux,trois).
	// -- OK
	//                 (string)      :-(un(deux, trois))
	//                 (pretty)      :-(un(deux, trois))
	// un.
	// -- OK
	//                 (string)      :-(un)
	//                 (pretty)      :-(un)
	// un(deux).
	// -- OK
	//                 (string)      :-(un(deux))
	//                 (pretty)      :-(un(deux))
	// un(2,3).
	// -- OK
	//                 (string)      :-(un(2, 3))
	//                 (pretty)      :-(un(2, 3))
	// empty().
	// -- OK
	//                 (string)      :-(empty())
	//                 (pretty)      :-(empty())
	// un(deux(),trois).
	// -- OK
	//                 (string)      :-(un(deux(), trois))
	//                 (pretty)      :-(un(deux(), trois))
}

func ExampleParseString_lists1() {

	tdata := []string{

		"[].",
		"[2].",
		"[2,3].",
		"[2,3,4].",

		"[[2,3],4].",
		"[2,[3,4]].",
	}

	runex(tdata)

	// Output:
	// 	[].
	// -- OK
	//                 (string)      :-(dot())
	//                 (pretty)      :-([])
	// [2].
	// -- OK
	//                 (string)      :-(dot(2, dot()))
	//                 (pretty)      :-([2])
	// [2,3].
	// -- OK
	//                 (string)      :-(dot(2, dot(3, dot())))
	//                 (pretty)      :-([2, 3])
	// [2,3,4].
	// -- OK
	//                 (string)      :-(dot(2, dot(3, dot(4, dot()))))
	//                 (pretty)      :-([2, 3, 4])
	// [[2,3],4].
	// -- OK
	//                 (string)      :-(dot(dot(2, dot(3, dot())), dot(4, dot())))
	//                 (pretty)      :-([[2, 3], 4])
	// [2,[3,4]].
	// -- OK
	//                 (string)      :-(dot(2, dot(dot(3, dot(4, dot())), dot())))
	//                 (pretty)      :-([2, [3, 4]])
}

func ExampleParseString_list2() {

	tdata := []string{

		"[[2,3],4].",
		//"[2,[3,4]].",

		// "[2|3].",
		// "[4|].",
		// "[4|X].",

		// "dot(1,dot(2,3)).",                   // not a list
		// "dot(1,dot(dot(4,dot(5,dot())),3)).", // not a list, but contains a list

	}

	runex(tdata)

	// Output:
	// [[2,3],4].
	// -- OK
	// (string)      :-(dot(dot(2, dot(3, dot())), dot(4, dot())))
	// (pretty)      :-([[2, 3], 4])

}

func ExampleParseString_invalid1() {

	tdata := []string{

		"un,deux.",
		"2(a).",
	}

	runex(tdata)

	// Output:
	// error in <un,deux.>, line 1 : syntax error: unexpected ',', expecting '.' or OPRULE
	// un,deux.
	// -- error :  Parse error
	// error in <2(a).>, line 1 : syntax error: unexpected INTEGER, expecting '[' or OPQUERY or ATOM
	// 2(a).
	// -- error :  Parse error

}

func runex(tdata []string) {

	for _, d := range tdata {

		r, err := ParseString(d, "<"+d+">")
		fmt.Println(d)
		if err != nil {
			fmt.Println("-- error : ", err)
		} else {
			fmt.Println("-- OK")
		}
		for _, v := range r {
			fmt.Println("\t\t(string)     ", v.String())
			fmt.Println("\t\t(pretty)     ", v.Pretty())
		}
	}
}
