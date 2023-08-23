package parser

//go:generate go install golang.org/x/tools/cmd/goyacc

//go:generate goyacc -o myparser.go -p "My" grammar.y
