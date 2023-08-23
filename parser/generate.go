package parser

//go:generate go install golang.org/x/tools/cmd/goyacc

//go:generate goyacc -o parser.go -p "my" grammar.y

//go:generate go fmt ./...
