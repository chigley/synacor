package machine

import "fmt"

type opCode int

const (
	opHalt opCode = 0
	opAdd         = 9
	opOut         = 19
	opNoop        = 21
)

type op struct {
	args   int
	writes bool
}

var ops = map[opCode]op{
	opAdd: {3, true},
	opOut: {1, false},
}

func (op opCode) String() string {
	switch op {
	case opHalt:
		return "halt"
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
