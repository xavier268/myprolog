
[![Go Reference](https://pkg.go.dev/badge/github.com/xavier268/myprolog.svg)](https://pkg.go.dev/github.com/xavier268/myprolog) 
[![Go Report Card](https://goreportcard.com/badge/github.com/xavier268/myprolog)](https://goreportcard.com/report/github.com/xavier268/myprolog)
# myprolog
A Prolog (the langage) implementation for educational purposes

# work in progress



* simplify VarIsVar with VarIsVar : do not try to be clever, just substitute in the rhs if we can (assumming VarIsVar is ordered)
    * with  Y = X,  
        * Z = Y => Z = X
        * Z = X => Z = X
    * VarIsVar substitution is dangerous since it can grow indefinitely the constraint set !

* find a way to handle diff ?
    * do we need constraints, or can we handle it as predicate only ?
    * if diff constraints, can we have a single constraint, capitalising on any other constraint(s)
