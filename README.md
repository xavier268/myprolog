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

* dot (a dot (b dot (c dot (nil nil)))) *//This is the canonical form. dot has arity of 2 exactly.*
* [ a b c ] *// This is the bracket form of the same list*
* [ a | [ b c ]]*// This is the bar form of the same*
* [  ] *// Not the same as nil, its canonical form is dot(nil, nil)*

Note that the bar form can produce object that cannot be represented as lists. 
Lists are **always** terminated with dot(nil,nil).

* [ a | b ] *// dot (a b)

Imbrication of lists is allowed :

* [X | [ a [ b | c ] [ e f ]]]


## General solving algorithm

SOLVE attempts to solve the goal node on top of the goal stack.
For instance, trying to SOLVE f(X,Y).

SOLVE will try to UNIFY successiveley with one of the rules.

The first rule, being rr(X,55,6), cannot be unified.
The second rule ..
The next rule ...

If we reach the end of the rule list without being able to UNIFY, then we report failure and return.

Assuming rule f(g(Y,3),Y)~... have the right functor, UNIFY will 
first RESCOPE it as f(g(Y#1, 3),Y)~ ..., making sure variables using the same symbol are not the same.

UNIFY will compare the goal with the rescoped head of the rule, returning a new list of constraints.
These constraints are recorded in the solving context.

*For instance, to unify f(X,Y) with f(g(Y#1,3),Y#1), we would create the constraints :*
*bind (f(X,Y) , f(g(Y#1,3),Y#1)).*

Then, we recursively process the constraints, together with the ones we already have, until either we have a canonical expression, or we demonstrate impossibility.

*bind(X,g(Y#1,3))*
*bind(Y,Y#1)*

If processing the constraints does not fail, the initial goal is replaced by the rescoped body of the unifed rule.

Then, SOLVE iterates on the next goal on the stack.



