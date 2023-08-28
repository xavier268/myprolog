// Code generated by goyacc -o parser.go -p my grammar.y. DO NOT EDIT.

//line grammar.y:4

package parser

import __yyfmt__ "fmt"

//line grammar.y:7

import ()

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

//line grammar.y:28
type mySymType struct {
	yys int
	// define the SymType structure
	list  []Term // list of Terms
	value Term   // single Term
}

const OPRULE = 57346
const OPQUERY = 57347
const ATOM = 57348
const STRING = 57349
const NUMBER = 57350
const VARIABLE = 57351
const LEXERROR = 57352

var myToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'('",
	"')'",
	"'.'",
	"','",
	"';'",
	"'['",
	"']'",
	"'|'",
	"'_'",
	"OPRULE",
	"OPQUERY",
	"ATOM",
	"STRING",
	"NUMBER",
	"VARIABLE",
	"LEXERROR",
	"'-'",
}

var myStatenames = [...]string{}

const myEofCode = 1
const myErrCode = 2
const myInitialStackSize = 16

//line grammar.y:166

//line yacctab:1
var myExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const myPrivate = 57344

const myLast = 85

var myAct = [...]int8{
	8, 36, 13, 19, 40, 12, 32, 48, 10, 10,
	10, 26, 20, 4, 7, 7, 7, 37, 26, 15,
	30, 34, 33, 38, 10, 47, 16, 25, 31, 43,
	21, 23, 27, 24, 42, 28, 41, 29, 44, 26,
	26, 46, 45, 35, 38, 17, 3, 10, 39, 11,
	25, 22, 9, 21, 23, 27, 24, 5, 28, 10,
	18, 2, 25, 1, 0, 21, 23, 27, 24, 10,
	28, 0, 25, 0, 0, 21, 23, 27, 24, 0,
	28, 14, 6, 0, 6,
}

var myPact = [...]int16{
	-1, -1000, -1, -1000, 1, 13, -1000, 41, -1000, -1000,
	50, -1000, 31, 12, 21, -1000, 0, 38, -1000, 7,
	37, 41, -1000, -1000, -1000, -1000, -1000, -1000, -13, -1000,
	1, 1, -1000, 23, 33, -1000, 16, -1000, 60, 15,
	-1000, -1000, -1000, -1000, -1000, -1000, -3, -1000, -1000,
}

var myPgo = [...]int8{
	0, 63, 61, 3, 46, 5, 2, 81, 57, 0,
	52, 1, 51,
}

var myR1 = [...]int8{
	0, 1, 2, 2, 4, 4, 4, 4, 8, 5,
	5, 6, 6, 7, 7, 9, 9, 9, 3, 3,
	11, 11, 11, 11, 11, 11, 10, 10, 10, 10,
	12, 12,
}

var myR2 = [...]int8{
	0, 1, 1, 2, 3, 2, 3, 4, 1, 1,
	3, 1, 3, 1, 1, 4, 3, 1, 1, 3,
	1, 1, 1, 1, 1, 1, 2, 3, 5, 4,
	1, 2,
}

var myChk = [...]int16{
	-1000, -1, -2, -4, 14, -8, -7, 15, -9, -10,
	9, -4, -5, -6, -7, 6, 13, 4, 10, -3,
	-11, 15, -12, 16, 18, 12, -9, 17, 20, 6,
	8, 7, 6, -5, -3, 5, -11, 10, 7, 11,
	17, -5, -6, 6, 5, -3, -11, 10, 10,
}

var myDef = [...]int8{
	0, -2, 1, 2, 0, 0, 8, 13, 14, 17,
	0, 3, 0, 9, 11, 5, 0, 0, 26, 0,
	18, 20, 21, 22, 23, 24, 25, 30, 0, 4,
	0, 0, 6, 0, 0, 16, 18, 27, 0, 0,
	31, 10, 12, 7, 15, 19, 0, 29, 28,
}

var myTok1 = [...]int8{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 5, 3, 3, 7, 20, 6, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 8,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 9, 3, 10, 3, 12, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 11,
}

var myTok2 = [...]int8{
	2, 3, 13, 14, 15, 16, 17, 18, 19,
}

var myTok3 = [...]int8{
	0,
}

var myErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	myDebug        = 0
	myErrorVerbose = false
)

type myLexer interface {
	Lex(lval *mySymType) int
	Error(s string)
}

type myParser interface {
	Parse(myLexer) int
	Lookahead() int
}

type myParserImpl struct {
	lval  mySymType
	stack [myInitialStackSize]mySymType
	char  int
}

func (p *myParserImpl) Lookahead() int {
	return p.char
}

func myNewParser() myParser {
	return &myParserImpl{}
}

const myFlag = -1000

func myTokname(c int) string {
	if c >= 1 && c-1 < len(myToknames) {
		if myToknames[c-1] != "" {
			return myToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func myStatname(s int) string {
	if s >= 0 && s < len(myStatenames) {
		if myStatenames[s] != "" {
			return myStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func myErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !myErrorVerbose {
		return "syntax error"
	}

	for _, e := range myErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + myTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(myPact[state])
	for tok := TOKSTART; tok-1 < len(myToknames); tok++ {
		if n := base + tok; n >= 0 && n < myLast && int(myChk[int(myAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if myDef[state] == -2 {
		i := 0
		for myExca[i] != -1 || int(myExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; myExca[i] >= 0; i += 2 {
			tok := int(myExca[i])
			if tok < TOKSTART || myExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if myExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += myTokname(tok)
	}
	return res
}

func mylex1(lex myLexer, lval *mySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(myTok1[0])
		goto out
	}
	if char < len(myTok1) {
		token = int(myTok1[char])
		goto out
	}
	if char >= myPrivate {
		if char < myPrivate+len(myTok2) {
			token = int(myTok2[char-myPrivate])
			goto out
		}
	}
	for i := 0; i < len(myTok3); i += 2 {
		token = int(myTok3[i+0])
		if token == char {
			token = int(myTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(myTok2[1]) /* unknown char */
	}
	if myDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", myTokname(token), uint(char))
	}
	return char, token
}

func myParse(mylex myLexer) int {
	return myNewParser().Parse(mylex)
}

func (myrcvr *myParserImpl) Parse(mylex myLexer) int {
	var myn int
	var myVAL mySymType
	var myDollar []mySymType
	_ = myDollar // silence set and not used
	myS := myrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	mystate := 0
	myrcvr.char = -1
	mytoken := -1 // myrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		mystate = -1
		myrcvr.char = -1
		mytoken = -1
	}()
	myp := -1
	goto mystack

ret0:
	return 0

ret1:
	return 1

mystack:
	/* put a state and value onto the stack */
	if myDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", myTokname(mytoken), myStatname(mystate))
	}

	myp++
	if myp >= len(myS) {
		nyys := make([]mySymType, len(myS)*2)
		copy(nyys, myS)
		myS = nyys
	}
	myS[myp] = myVAL
	myS[myp].yys = mystate

mynewstate:
	myn = int(myPact[mystate])
	if myn <= myFlag {
		goto mydefault /* simple state */
	}
	if myrcvr.char < 0 {
		myrcvr.char, mytoken = mylex1(mylex, &myrcvr.lval)
	}
	myn += mytoken
	if myn < 0 || myn >= myLast {
		goto mydefault
	}
	myn = int(myAct[myn])
	if int(myChk[myn]) == mytoken { /* valid shift */
		myrcvr.char = -1
		mytoken = -1
		myVAL = myrcvr.lval
		mystate = myn
		if Errflag > 0 {
			Errflag--
		}
		goto mystack
	}

mydefault:
	/* default state action */
	myn = int(myDef[mystate])
	if myn == -2 {
		if myrcvr.char < 0 {
			myrcvr.char, mytoken = mylex1(mylex, &myrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if myExca[xi+0] == -1 && int(myExca[xi+1]) == mystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			myn = int(myExca[xi+0])
			if myn < 0 || myn == mytoken {
				break
			}
		}
		myn = int(myExca[xi+1])
		if myn < 0 {
			goto ret0
		}
	}
	if myn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			mylex.Error(myErrorMessage(mystate, mytoken))
			Nerrs++
			if myDebug >= 1 {
				__yyfmt__.Printf("%s", myStatname(mystate))
				__yyfmt__.Printf(" saw %s\n", myTokname(mytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for myp >= 0 {
				myn = int(myPact[myS[myp].yys]) + myErrCode
				if myn >= 0 && myn < myLast {
					mystate = int(myAct[myn]) /* simulate a shift of "error" */
					if int(myChk[mystate]) == myErrCode {
						goto mystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if myDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", myS[myp].yys)
				}
				myp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if myDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", myTokname(mytoken))
			}
			if mytoken == myEofCode {
				goto ret1
			}
			myrcvr.char = -1
			mytoken = -1
			goto mynewstate /* try again in the same state */
		}
	}

	/* reduction by production myn */
	if myDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", myn, myStatname(mystate))
	}

	mynt := myn
	mypt := myp
	_ = mypt // guard against "declared and not used"

	myp -= int(myR2[myn])
	// myp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if myp+1 >= len(myS) {
		nyys := make([]mySymType, len(myS)*2)
		copy(nyys, myS)
		myS = nyys
	}
	myVAL = myS[myp+1]

	/* consult goto table to find next state */
	myn = int(myR1[myn])
	myg := int(myPgo[myn])
	myj := myg + myS[myp].yys + 1

	if myj >= myLast {
		mystate = int(myAct[myg])
	} else {
		mystate = int(myAct[myj])
		if int(myChk[mystate]) != -myn {
			mystate = int(myAct[myg])
		}
	}
	// dummy call; replaced with literal code
	switch mynt {

	case 1:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:46
		{
			myVAL.list = myDollar[1].list
			// save final result in the provided variable, results
			lastParseResult = myVAL.list
		}
	case 2:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:53
		{
			myVAL.list = append(myVAL.list, myDollar[1].value)
		}
	case 3:
		myDollar = myS[mypt-2 : mypt+1]
//line grammar.y:54
		{
			myVAL.list = append(myDollar[1].list, myDollar[2].value)
		}
	case 4:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:57
		{
			myVAL.value = CompoundTerm{
				Functor:  "query",
				Children: []Term{myDollar[2].value},
			}
		}
	case 5:
		myDollar = myS[mypt-2 : mypt+1]
//line grammar.y:63
		{ // implicit OPRULE
			myVAL.value = CompoundTerm{
				Functor:  "rule",
				Children: []Term{myDollar[1].value},
			}
		}
	case 6:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:69
		{
			myVAL.value = CompoundTerm{
				Functor:  "rule",
				Children: []Term{myDollar[1].value},
			}
		}
	case 7:
		myDollar = myS[mypt-4 : mypt+1]
//line grammar.y:75
		{
			myVAL.value = CompoundTerm{
				Functor:  "rule",
				Children: []Term{myDollar[1].value, myDollar[3].value},
			}
		}
	case 8:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:82
		{
			myVAL.value = myDollar[1].value
		}
	case 9:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:87
		{
			myVAL.value = myDollar[1].value
		}
	case 10:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:88
		{
			myVAL.value = CompoundTerm{
				Functor:  "or",
				Children: []Term{myDollar[1].value, myDollar[3].value},
			}
		}
	case 11:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:95
		{
			myVAL.value = myDollar[1].value
		}
	case 12:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:96
		{
			myVAL.value = CompoundTerm{
				Functor:  "and",
				Children: []Term{myDollar[1].value, myDollar[3].value},
			}
		}
	case 13:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:103
		{
			myVAL.value = myDollar[1].value
		}
	case 14:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:104
		{
			myVAL.value = myDollar[1].value
		}
	case 15:
		myDollar = myS[mypt-4 : mypt+1]
//line grammar.y:107
		{
			myVAL.value = CompoundTerm{
				Functor:  myDollar[1].value.String(),
				Children: myDollar[3].list,
			}
		}
	case 16:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:112
		{
			myVAL.value = CompoundTerm{
				Functor: myDollar[1].value.String(),
			}
		}
	case 17:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:117
		{
			myVAL.value = myDollar[1].value
		}
	case 18:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:122
		{
			myVAL.list = []Term{myDollar[1].value}
		}
	case 19:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:123
		{
			myVAL.list = append([]Term{myDollar[1].value}, myDollar[3].list...)
		}
	case 20:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:128
		{
			myVAL.value = myDollar[1].value
		}
	case 21:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:129
		{
			myVAL.value = myDollar[1].value
		}
	case 22:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:130
		{
			myVAL.value = myDollar[1].value
		}
	case 23:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:131
		{
			myVAL.value = myDollar[1].value
		}
	case 24:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:132
		{
			myVAL.value = myDollar[1].value
		}
	case 25:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:133
		{
			myVAL.value = myDollar[1].value
		}
	case 26:
		myDollar = myS[mypt-2 : mypt+1]
//line grammar.y:137
		{
			myVAL.value = CompoundTerm{
				Functor:  "dot",
				Children: []Term{},
			}
		}
	case 27:
		myDollar = myS[mypt-3 : mypt+1]
//line grammar.y:142
		{
			myVAL.value = newList(myDollar[2].list...)
		}
	case 28:
		myDollar = myS[mypt-5 : mypt+1]
//line grammar.y:145
		{
			myVAL.value = CompoundTerm{
				Functor:  "dot",
				Children: []Term{myDollar[2].value, myDollar[4].value},
			}
		}
	case 29:
		myDollar = myS[mypt-4 : mypt+1]
//line grammar.y:150
		{
			myVAL.value = CompoundTerm{
				Functor:  "dot",
				Children: []Term{myDollar[2].value},
			}
		}
	case 30:
		myDollar = myS[mypt-1 : mypt+1]
//line grammar.y:158
		{
			myVAL.value = myDollar[1].value
		}
	case 31:
		myDollar = myS[mypt-2 : mypt+1]
//line grammar.y:159
		{
			myVAL.value = Number{
				Num: -myDollar[2].value.(Number).Num,
				Den: myDollar[2].value.(Number).Den,
			}
		}
	}
	goto mystack /* stack new state and value */
}
