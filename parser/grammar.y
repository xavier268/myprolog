// expression grammar for MyProlog


%{
 
 // Autogenerated file - DO NOT EDIT
 package parser


 import (
 
 )

// init global variables
func init() {
    // Set verbose error
    myErrorVerbose = true
    // set debug level
    myDebug = 0
}

// where the parse results are available
var lastParseResult []Term

// To keep the compiler happy ...
 var _ = __yyfmt__.Printf

%}

%union{
    // define the SymType structure
    list []Term // list of Terms
    value  Term // single Term
}

%type <list> top phrases params
%type <value> phrase disjterms conjterms conjterm 
%type <value> compterm number list param


%token <value> '(' ')' '.' ',' ';' '[' ']' '|' '_'
%token <value> OPRULE OPQUERY // :-  and ?-
%token <value> ATOM STRING INTEGER FLOAT VARIABLE

%% 

top:
    phrases                             { 
                                        $$ = $1
                                        // save final result in the provided variable, results
                                        lastParseResult = $$
                                        }

phrases:
    phrase                              { $$ = append( $$ , $1 )}
    | phrases phrase                    { $$ = append( $1 , $2 )}

phrase: 
    OPQUERY disjterms '.'                { 
                                        $$ = &CompoundTerm{  
                                                Functor : "query", 
                                                Children: []Term{$2},
                                                }                                                
                                        }
    | conjterm '.'                      {  // implicit OPRULE
                                        $$ = &CompoundTerm{
                                                Functor : "rule", 
                                                Children: []Term{ $1 } ,
                                                };
                                        }
    | conjterm OPRULE '.'               { 
                                        $$ = &CompoundTerm{
                                                Functor : "rule",    
                                                Children: []Term{ $1} ,
                                                };
                                        }
    | conjterm OPRULE disjterms '.'     { 
                                        $$ = &CompoundTerm{
                                                Functor : "rule",    
                                                Children: []Term{ $1, $3},
                                                }
                                        }

disjterms:
    conjterms                           { $$ = $1 }
    | conjterms ';' disjterms           { $$ = &CompoundTerm{
                                                Functor : "or",
                                                Children : []Term{ $1, $3},
                                                }
                                        }

conjterms:
    conjterm                            { $$ = $1 }
    | conjterm ',' conjterms            { $$ = &CompoundTerm {
                                                Functor : "and",
                                                Children: []Term{$1 , $3},
                                                }
                                        }

conjterm:
    ATOM                                { $$ = $1 }
    | compterm                          { $$ = $1 }

compterm:
    ATOM '(' params ')'                 { $$ = &CompoundTerm {
                                                Functor : $1.String(),
                                                Children : $3,
                                                }
                                        }
    | ATOM '(' ')'                      {
                                            $$ = &CompoundTerm {
                                                Functor : $1.String(),
                                            }
                                        }
    | list                              { 
                                            $$ = $1 
                                        }

params:                                
    param                               { $$ = []Term{ $1 } }
    | param ',' params                  { 
                                            $$ = append([]Term{$1}, $3...)                                           
                                        }

param:
    ATOM                                { $$ = $1 }
    | number                            { $$ = $1 }
    | STRING                            { $$ = $1 }
    | VARIABLE                          { $$ = $1 }
    | '_'                               { $$ = $1 }
    | compterm                          { $$ = $1 }

number:
    INTEGER                             { $$ = $1 }
    | FLOAT                             { $$ = $1 }


list: 
    '[' ']'                             { $$ = &CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{},
                                            } 
                                        }
    | '['  params ']'                   {
                                            $$ = newList( $2 ... )                                            
                                        }
    | '[' param '|' param ']'           { $$ = &CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{$2 , $4} ,
                                            } 
                                        }
    | '[' param '|' ']'                 {
                                            $$ = &CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{$2} ,
                                            } 
                                        }


%%

