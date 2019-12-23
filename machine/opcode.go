package machine

import (
	"fmt"
)

type opCode int

const (
	opHalt opCode = iota
	opSet
	opPush
	opPop
	opEq
	opGt
	opJmp
	opJt
	opJf
	opAdd
	opMult
	opMod
	opAnd
	opOr
	opNot
	opRmem
	opWmem
	opCall
	opRet
	opOut
	opIn
	opNoop
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
	opRmem: {"rmem", 2, true},
	opWmem: {"wmem", 2, false}, // false because the argument is a literal address
	opCall: {"call", 1, false},
	opRet:  {"ret", 0, false},
	opOut:  {"out", 1, false},
	opIn:   {"in", 1, true},
	opNoop: {"noop", 0, false},
}

func (op opCode) String() string {
	if o, ok := ops[op]; ok {
		return o.label
	}
	return fmt.Sprintf("UnknownOpCode(%d)", op)
}
