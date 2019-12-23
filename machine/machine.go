package machine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"go.uber.org/zap"
)

type Machine struct {
	memory []uint16
	pc     uint16

	registers [8]uint16
	stack     []uint16

	out    io.Writer
	logger *zap.Logger
}

const modulus = math.MaxInt16 + 1

var errHalt = errors.New("machine: halted")

func New(r io.Reader, opts ...Option) (*Machine, error) {
	prg, err := readProgram(r)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		logger:    zap.NewNop(),
		outWriter: os.Stdout,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Machine{
		memory: prg,
		out:    cfg.outWriter,
		logger: cfg.logger,
	}, nil
}

func (m *Machine) Run() error {
	for {
		if err := m.step(); errors.Is(err, errHalt) {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func (m *Machine) peek(addr uint16) uint16 {
	return m.memory[addr]
}

func (m *Machine) readArgument() (uint16, error) {
	val := m.peek(m.pc)
	m.pc++

	if val <= math.MaxInt16 {
		return val, nil
	}
	if val <= math.MaxInt16+8 {
		return m.registers[val-modulus], nil
	}
	return 0, fmt.Errorf("machine: invalid read source %d", val)
}

func (m *Machine) writeArgument(arg, val uint16) {
	m.registers[arg-modulus] = val
}

func (m *Machine) writeBool(arg uint16, p bool) {
	var val uint16
	if p {
		val = 1
	}
	m.writeArgument(arg, val)
}

func (m *Machine) step() error {
	pc := m.pc
	opCode := opCode(m.peek(pc))
	m.pc++

	var args []uint16
	if op := ops[opCode]; op.args > 0 {
		args = make([]uint16, op.args)
		for i := 0; i < op.args; i++ {
			if i == 0 && op.writes {
				args[i] = m.peek(m.pc)
				m.pc++
			} else {
				arg, err := m.readArgument()
				if err != nil {
					return err
				}
				args[i] = arg
			}
		}
	}

	m.logger.Debug("step",
		zap.Uint16("pc", pc),
		zap.String("opcode", opCode.String()),
		zap.Uint16s("args", args),
	)

	switch opCode {
	case opHalt:
		return errHalt
	case opSet:
		m.writeArgument(args[0], args[1])
	case opPush:
		m.stack = append(m.stack, args[0])
	case opPop:
		var val uint16
		val, m.stack = m.stack[len(m.stack)-1], m.stack[:len(m.stack)-1]
		m.writeArgument(args[0], val)
	case opEq:
		m.writeBool(args[0], args[1] == args[2])
	case opGt:
		m.writeBool(args[0], args[1] > args[2])
	case opJmp:
		m.pc = args[0]
	case opJt:
		if args[0] > 0 {
			m.pc = args[1]
		}
	case opJf:
		if args[0] == 0 {
			m.pc = args[1]
		}
	case opAnd:
		m.writeArgument(args[0], args[1]&args[2])
	case opOr:
		m.writeArgument(args[0], args[1]|args[2])
	case opNot:
		m.writeArgument(args[0], (^args[1])&math.MaxInt16)
	case opAdd:
		m.writeArgument(args[0], (args[1]+args[2])%modulus)
	case opOut:
		_, err := m.out.Write([]byte{byte(args[0])})
		return err
	case opNoop:
	default:
		return fmt.Errorf("machine: unsupported opcode %d", opCode)
	}
	return nil
}

func readProgram(r io.Reader) ([]uint16, error) {
	var program []uint16
	for {
		bs := make([]byte, 2)
		if _, err := io.ReadFull(r, bs); err == io.EOF {
			return program, nil
		} else if err != nil {
			return nil, err
		}
		program = append(program, binary.LittleEndian.Uint16(bs))
	}
}
