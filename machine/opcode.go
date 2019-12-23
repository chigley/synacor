package machine

import "fmt"

type opCode int

const (
	opHalt opCode = 0
	opSet         = 1
	opPush        = 2
	opEq          = 4
	opJmp         = 6
	opJt          = 7
	opJf          = 8
	opAdd         = 9
	opOut         = 19
	opNoop        = 21
)

type op struct {
	args   int
	writes bool
}

var ops = map[opCode]op{
	opSet:  {2, true},
	opPush: {1, false},
	opEq:   {3, true},
	opJmp:  {1, false},
	opJt:   {2, false},
	opJf:   {2, false},
	opAdd:  {3, true},
	opOut:  {1, false},
}

func (op opCode) String() string {
	switch op {
	case opHalt:
		return "halt"
	case opSet:
		return "set"
	case opPush:
		return "push"
	case opEq:
		return "eq"
	case opJmp:
		return "jmp"
	case opJt:
		return "jt"
	case opJf:
		return "jf"
	case opAdd:
		return "add"
	case opOut:
		return "out"
	case opNoop:
		return "noop"
	default:
		return fmt.Sprintf("UnknownOpCode(%d)", op)
	}
}
