// expression grammar for MyProlog


%{

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

// To keep the compiler happy ...
 var _ = __yyfmt__.Printf

%}

%union{
    // define the SymType structure
    list []Term // list of Terms
    value  Term // single Term
}

%type <list> top phrases params
%type <value> phrase disjterms conjterms conjterm head
%type <value> compterm list param number contraint


%token <value> '(' ')' '.' ',' ';' '[' ']' '|' '_' '=' '!'
%token <value> OPRULE OPQUERY // :-  and ?-
%token <value> ATOM STRING NUMBER VARIABLE 
%token <value> LEXERROR

%% 

top:
    phrases                             { 
                                        $$ = $1
                                        // save final result in the provided variable, results
                                        mylex.(*myLex).LastResult = $$
                                        }

phrases:
    phrase                              { $$ = append( $$ , $1 )}
    | phrases phrase                    { $$ = append( $1 , $2 )}

phrase: 
    OPQUERY disjterms '.'                { 
                                        $$ = CompoundTerm{  
                                                Functor : "query", 
                                                Children: []Term{$2},
                                                }                                                
                                        }
    | head '.'                      {  // implicit OPRULE
                                        $$ = CompoundTerm{
                                                Functor : "rule", 
                                                Children: []Term{ $1 } ,
                                                };
                                        }
    | head OPRULE '.'               { 
                                        $$ = CompoundTerm{
                                                Functor : "rule",    
                                                Children: []Term{ $1} ,
                                                };
                                        }
    | head OPRULE disjterms '.'     { 
                                        $$ = CompoundTerm{
                                                Functor : "rule",    
                                                Children: []Term{ $1, $3},
                                                }
                                        }


head:  
    ATOM                                { $$ = $1 }
    | compterm                          { $$ = $1 }
    


disjterms:
    conjterms                           { $$ = $1 }
    | conjterms ';' disjterms           { $$ = CompoundTerm{
                                                Functor : "or",
                                                Children : []Term{ $1, $3},
                                                }
                                        }

conjterms:
    conjterm                            { $$ = $1 }
    | conjterm ',' conjterms            { $$ =CompoundTerm {
                                                Functor : "and",
                                                Children: []Term{$1 , $3},
                                                }
                                        }

conjterm:
    ATOM                                { $$ = $1 }
    | compterm                          { $$ = $1 }
    | contraint                         { $$ = $1 }
    | '!'                               { $$ = Atom {"cut"}}

contraint:
    param '=' param                     { $$ = CompoundTerm {
                                                Functor:  "eq" ,
                                                Children: []Term{ $1, $3},
                                                }
                                        }
    

compterm:
    ATOM '(' params ')'                 { $$ = CompoundTerm {
                                                Functor : $1.String(),
                                                Children : $3,
                                                }
                                        }
    | ATOM '(' ')'                      {
                                            $$ = CompoundTerm {
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


list: 
    '[' ']'                             { $$ = CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{},
                                            } 
                                        }
    | '['  params ']'                   {
                                            $$ = newList( $2 ... )                                            
                                        }
    | '[' param '|' param ']'           { $$ = CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{$2 , $4} ,
                                            } 
                                        }
    | '[' param '|' ']'                 {
                                            $$ = CompoundTerm{
                                            Functor : "dot",
                                            Children : []Term{$2} ,
                                            } 
                                        }

                                        
number :        NUMBER                  { $$ = $1 }
    |   '-' NUMBER                      { 
                                            $$ = Number {
                                                Num : -$2.(Number).Num, 
                                                Den: $2.(Number).Den,
                                                }
                                        }

%%

