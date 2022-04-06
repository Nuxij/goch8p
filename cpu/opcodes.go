package cpu

import "fmt"

var AllOpcodes = []Instruction{
	OxClearScreen{Opcode{0x00E0, "Clear Screen"}},
}

type InstructionHandler interface {
	HandleInstruction(uint16) error
}
type Instruction interface {
	Name() string
	Register(cpu *CPU) InstructionHandler
}
type InstructionHandlerFunc func(uint16) error

func (f InstructionHandlerFunc) HandleInstruction(instruction uint16) error {
	return f(instruction)
}

type InstructionUnknown struct {
	opcode uint16
}

func (iu InstructionUnknown) Error() string {
	return fmt.Sprintf("unknown instruction: %X", iu.opcode)
}

type InstructionOk struct {
	opcode uint16
}

func (io InstructionOk) Error() string {
	return fmt.Sprintf("success: %X", io.opcode)
}

type Opcode struct {
	opcode uint16
	name   string
}

func (i Opcode) Name() string {
	return i.name
}

type OxClearScreen struct {
	Opcode
}

func (o OxClearScreen) Register(cpu *CPU) InstructionHandler {
	return InstructionHandlerFunc(func(op uint16) error {
		if op == o.opcode {
			for i := uint16(0); i < cpu.screen.size; i++ {
				cpu.ram.Write(cpu.screen.address+i, 0x0)
			}
			return InstructionOk{op}
		}
		return nil
	})
}
