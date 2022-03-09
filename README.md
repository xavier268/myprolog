# myprolog
A Prolog (the langage) implementation for educationnal purposes


## Types

Variables starts by an uppercase letter, or an underscore, _
Binding to variables starting with underscore can be ignored, as _ in golang.

Numbers are recognized because they start with a digit.

## Syntax parsing

Scanner respect the golang grammar.

Commas (,) are ignored.
Comments and white spaces are ignored.

Variables, Numbers, (,),[,], cannot be functors of compound object.
We maintain a symbol list for the non functors symbols.

Subtrees NOT containing Variables are identified as "constant". 
It is a different notion than non-functors.
They will require no rescoping for new rule contexts.

## rules

Rules can have the following forms :

* ~ ( a b ) . *// This is the canonical form**
* f(X,Y) ~ f(1, Y) , f(X,2) .  *//commas are optionals. The final period is mandatory*
* f(X,Y) ~ one ; two .  *//this form is short hand to define two rules with same head*
* f(X,Y). *//final period required. This is a head only rule. Same as f(Y,Y) ~ .*

## Lists

Lists can take the following forms :

* dot (a dot (b dot (c dot (d nil)))) *//This is the canonical form. dot has arity of 2 exactly.*
* [ a b c d ] *// This is the bracket form.*
* [ a | RestOfTheList ] or [ a | [ b c d ]]*// This is the bar form*
* [  ] *// Same as nil, its canonical form*

Imbrication of lists is allowed :

* [X | [ a [ b | c ] [ e f ]]]
