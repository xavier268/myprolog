// package parser contains lexer, parser and ast definitions.
// run go generate ./... to compile new parser

package parser

//go:generate go get golang.org/x/tools/cmd/goyacc

//go:generate go install golang.org/x/tools/cmd/goyacc

//go:generate goyacc -o parser.go -p "my" grammar.y

//go:generate go mod tidy

//go:generate go fmt ./...

const VERSION = "0.8.6"
