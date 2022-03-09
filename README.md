# myprolog
A Prolog (the langage) implementation for educationnal purposes


## Types

Variables starts by an uppercase letter, or an underscore, _
Binding to variables starting with underscore care ingnored, as _ in golang.

Numbers are recognized because they start with a digit.

## Syntax parsing

Scanner respect the golang grammar.

Commas (,) are ignored.
Comments and white spaces are ignored.

Variables or Numbers cannot be functors of compound object.
Therefore, we maintain a symbol list for these.

Subtrees NOT containing Variables are identified as "constant". 
They will require no rescoping for new rule contexts.

## rules

Rules can have the following forms :

* f(X,Y) ~ f(1, Y) , f(X,2) .  *//commas are optionals. The final period is mandatory*
* f(X,Y) ~ one ; two .  *//this form is short hand to define two rules with same head*
* f(X,Y). *//final period required. This is a head only rule. Same as f(Y,Y) ~ .*