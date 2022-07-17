
# myprolog
A Prolog (the langage) implementation for educational purposes

## How to use 

Launch the interpreter :

``` 

cd cli
go run .

```

Then, at the command promt, enter some rules or facts and query them.

```

ok>
father(john, paul).
father(john, henry).
?father(john,A).

Results : [A = paul, ]
Ok> 

```


## Basic syntax

The general idea is that the program defines a number of **rules**  (facts are a special case of rules) and then try to find a solution the **queries** you formulate, binding query **variables** to potential solutions.

Variables start with a capital letter A-Z. Anything else can be used as facts, rules or queries (except for limited reserved keywords).

Some examples : 

```
// First, let's define some rules ... 

grandfather(A,B) :- father(A, S) , father (S B).  // we define the concept of grandfather.
grandfather(paul , john).                         // We know that Paul is John's grand father.

// There can be multiple definitions of a grandfather ...
grandfather(X,Y) :- father(X,D) , mother (D, Y).

// or, shall we define the notion of parent ?
parent(A,B) :- mother(A,B) ; father(A,B).         // Did you notice the semi colon to specify an alternative ?

/*  Then, we can query the known rules and facts 
    Notice how queries always start with a question mark (?).
    ... And notice this as a block comment ;-) */

// Give me someone who's grandfathr is paul ?
? grandfather (paul, S). // We will get a response like : S = john.

// Anyone having a grand father ?
? grandfather (GF, _).  // We will get a response like : GF = paul.
                        // Dis you notice the use of the undesrcore _ variable, to mean anything ?

// Same as above, but without the underscore 
? grandfather (GF, S).  / We will get a response like : GF = paul, S = john.

```

## More about the syntax

Spaces and commas are not significant nor needed as long there is no ambiguity for the scanner.

The _ (underscore) variable is a special variable that can match anything.

Strings can be quoted or not. If not quoted, the scanner used mainly the Golang syntax. For instance, 3x is the same a 3 , x. 

Any non variable is a legal object name. Redefine + or £ if you so wish ...

## Reserved words

( to be completed )

Have fun !



