package machine

import (
	"fmt"
)

type opCode int

const (
	opHalt opCode = 0
	opSet         = 1
	opPush        = 2
	opPop         = 3
	opEq          = 4
	opGt          = 5
	opJmp         = 6
	opJt          = 7
	opJf          = 8
	opAdd         = 9
	opMult        = 10
	opMod         = 11
	opAnd         = 12
	opOr          = 13
	opNot         = 14
	opCall        = 17
	opRet         = 18
	opOut         = 19
	opNoop        = 21
)

type op struct {
	label  string
	args   int
	writes bool
}

var ops = map[opCode]op{
	opHalt: {"halt", 0, false},
	opSet:  {"set", 2, true},
	opPush: {"push", 1, false},
	opPop:  {"pop", 1, true},
	opEq:   {"eq", 3, true},
	opGt:   {"gt", 3, true},
	opJmp:  {"jmp", 1, false},
	opJt:   {"jt", 2, false},
	opJf:   {"jf", 2, false},
	opAdd:  {"add", 3, true},
	opMult: {"mult", 3, true},
	opMod:  {"mod", 3, true},
	opAnd:  {"and", 3, true},
	opOr:   {"or", 3, true},
	opNot:  {"and", 2, true},
	opCall: {"call", 1, false},
	opRet:  {"ret", 0, false},
	opOut:  {"out", 1, false},
	opNoop: {"noop", 0, false},
}

func (op opCode) String() string {
	if o, ok := ops[op]; ok {
		return o.label
	}
	return fmt.Sprintf("UnknownOpCode(%d)", op)
}
