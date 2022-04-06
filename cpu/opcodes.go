package cpu

import (
	"encoding/binary"
	"fmt"
)

var AllOpcodes = []Instruction{
	OxClearScreen{Opcode{0x00E0, "Clear Screen"}},
	OxDrawSprite{Opcode{0xD000, "Draw Sprite"}},
	OxReturn{Opcode{0x00EE, "Return from subroutine"}},
	OxCall{Opcode{0x2000, "Call Subroutine"}},
	OxLoadIndex{Opcode{0xA000, "Load Index"}},
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

type InstructionNOP struct {
	opcode uint16
}

func (io InstructionNOP) Error() string {
	return fmt.Sprintf("NOP: %X", io.opcode)
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
			return nil
		}
		return InstructionNOP{op}
	})
}

type OxDrawSprite struct {
	Opcode
}
func (o OxDrawSprite) Register(cpu *CPU) InstructionHandler {
	return InstructionHandlerFunc(func(op uint16) error {
		if op & 0xF000 == o.opcode {
			x := op & 0x0F00 >> 8
			y := op & 0x00F0 >> 4
			height := op & 0x000F
			spriteAddress := cpu.index
			xCoord := uint16(cpu.v[x] % 64)
			yCoord := uint16(cpu.v[y] % 32)
			cpu.v[0xF] = 0


			sprite, err := cpu.ram.Reads(spriteAddress, height)
			if err != nil {
				return err
			}
			for yPos, b := range sprite {
				for xPos := uint16(0); xPos < 8; xPos++ {
					if xCoord + xPos >= 64 || yCoord + uint16(yPos) >= 32 {
						continue
					}
					pos := (y + uint16(yPos)) * 64 + x + xPos
					pos += cpu.screen.address
					onScreen, err := cpu.ram.Read(pos)
					if err != nil {
						return err
					}
					toBe := b & (0x80 >> xPos)
					if toBe != 0 && onScreen != 0 {
						cpu.ram.Write(pos, 0)
						cpu.v[0xF] = 1
					} else if toBe != 0 && onScreen == 0 {
						cpu.ram.Write(pos, 1)
					}
				}
			}
			return nil
		}
		return InstructionNOP{op}
	})
}

type OxCall struct {
	Opcode
}

func (o OxCall) Register(cpu *CPU) InstructionHandler {
	return InstructionHandlerFunc(func(op uint16) error {
		if op & 0xF000 == o.opcode {
			ok := cpu.stack.Push(cpu.pc)
			if !ok {
				return fmt.Errorf("Stack overflow")
			}
			cpu.SetPC(op & 0x0FFF)
			return nil
		}
		return InstructionNOP{op}
	})
}

type OxReturn struct {
	Opcode
}	

func (o OxReturn) Register(cpu *CPU) InstructionHandler {
	return InstructionHandlerFunc(func(op uint16) error {
		if op & 0xF000 == o.opcode {
			d, ok := cpu.stack.Pop()
			if !ok {
				return fmt.Errorf("Stack is empty")
			}
			cpu.SetPC(d)
			return nil
		}
		return InstructionNOP{op}
	})
}

type OxLoadIndex struct {
	Opcode
}

func (o OxLoadIndex) Register(cpu *CPU) InstructionHandler {
	return InstructionHandlerFunc(func(op uint16) error {
		if op & 0xF000 == o.opcode {
			data, err := cpu.ram.Reads(op&0x0FFF, 2)
			if err != nil {
				return err
			}
			cpu.index = binary.BigEndian.Uint16(data)
			cpu.IncrementPC(2)
			return nil
		}
		return InstructionNOP{op}
	})
}