# myprolog
A Prolog (the langage) implementation for educationnal purposes


## Types

Variables starts by an uppercase letter, or is the special underscore Variable, _

Numbers are recognized because they start with a digit.

## Syntax parsing

Commas (,) are ignored.
Comments and white spaces are ignored.

Variables or Numbers cannot be functors of compound object.
Therefore, we maintain a symbol list for these.

Subtrees NOT containing Variables are identified as "constant". 
They will require no rescoping for new rule contexts.
