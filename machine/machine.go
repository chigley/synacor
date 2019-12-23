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
	memory    []uint16
	registers [8]uint16
	pc        uint16

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

	prg = []uint16{9, 32768, 32769, 4, 19, 32768, 0}

	return &Machine{
		memory: prg,
		out:    cfg.outWriter,
		logger: cfg.logger,
	}, nil
}

func (m *Machine) Run() error {
	m.registers[1] = 61
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

func (m *Machine) step() error {
	opCode := opCode(m.peek(m.pc))
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
		zap.String("opcode", opCode.String()),
		zap.Uint16s("args", args),
	)

	switch opCode {
	case opHalt:
		return errHalt
	case opAdd:
		m.writeArgument(args[0], (args[1]+args[2])%modulus)
		return nil
	case opOut:
		_, err := m.out.Write([]byte{byte(args[0])})
		return err
	case opNoop:
		return nil
	default:
		return fmt.Errorf("machine: unsupported opcode %d", opCode)
	}
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