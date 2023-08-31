package main

import (
	"fmt"

	"github.com/xavier268/myprolog/parser"

	"github.com/xavier268/myprolog/solver"
)

func main() {
	fmt.Println("Welcome to prolog !")
	fmt.Printf("Version parser:%s - solver:%s\n", parser.VERSION, solver.VERSION)
	fmt.Println("(c) 2022, 2023 Xavier Gandillot (aka xavier268)")
	solver.Repl()
}
