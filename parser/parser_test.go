package parser

import (
	"fmt"
)

func ExampleParseString() {

	tdata := []string{
		// valid
		"un(deux,trois).",
		"un.",
		"un(deux).",
		"un(2,3).",
		"empty().",
		"un(deux(),trois).",

		"[].",
		"[2].",
		"[2,3].",
		"[2,3,4].",

		"[2|3].",
		"[4|].",
		"[4|X].",

		// invalid
		"un,deux.",
		"2(a).",
	}

	for _, d := range tdata {

		r, err := ParseString(d, "<"+d+">")
		fmt.Println(d+"\t\terror : ", err)
		for _, v := range r {
			fmt.Println("\tstring:\t", v.String())
			fmt.Println("\tpretty:\t", v.Pretty())
		}

		// Output:

	}
}
