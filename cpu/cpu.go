package cpu

import (
	"fmt"

	"github.com/google/uuid"
)

type CPU struct {
	ram     Device
	pc      uint16
	sp      uint16
	index   uint16
	stack   *Stack
	screen  Screen
	opcodes []InstructionHandler
}

func RandomStringUUID() string {
	id, _ := uuid.NewUUID()
	return fmt.Sprintf("%v", id)
}

func (c *CPU) CallInstruction(handler InstructionHandler, opcode uint16) error {
	return handler.HandleInstruction(opcode)
}

type Screen struct {
	width   uint16
	height  uint16
	size    uint16
	address uint16
}

func NewCPU(ram Device) *CPU {
	addressableSize := uint16(0x1000)
	screenSize := uint16(64*32) / 8
	screenAddress := addressableSize - (screenSize + 1)
	rammer := NewRammer(256, []Device{ram})
	rammer.SetRegion(0x0, addressableSize, ram, 0x0)
	rammer.SetRegion(screenAddress, screenSize, ram, screenAddress)
	cpu := &CPU{
		opcodes: []InstructionHandler{},
		ram:     rammer,
		pc:      0x200,
		sp:      0x0,
		index:   0x0,
		stack:   NewStack(0x10),
		screen: Screen{
			width:   64,
			height:  32,
			size:    screenSize,
			address: screenAddress,
		},
	}
	for _, i := range AllOpcodes {
		cpu.opcodes = append(cpu.opcodes, i.Register(cpu))
	}
	return cpu
}

func (c *CPU) ExecuteInstruction(instruction uint16) error {
	var err error
	for _, op := range c.opcodes {
		err = c.CallInstruction(op, instruction)
		if err != nil {
			return err
		}
	}
	// fmt.Printf("Executing instruction: %x\n", instruction)
	switch instruction >> 12 {
	case 0xA:
		err = c.loadIndex(instruction)
	case 0x2:
		err = c.callSubroutine(instruction)
	case 0x0:
		switch instruction & 0x00FF {
		case 0x00EC:
			err = c.yieldCoroutine()
		case 0x00EE:
			err = c.returnFromSubroutine()
		default:
			err = InstructionUnknown{instruction}
		}
	default:
		err = InstructionUnknown{instruction}
	}
	return err
}

func (c *CPU) IncrementPC(count uint16) {
	c.pc += count
}

func (c *CPU) SetPC(addr uint16) {
	c.pc = addr
}

func (c *CPU) loadIndex(instruction uint16) error {
	data, err := c.ram.Reads(instruction&0x0FFF, 2)
	if err != nil {
		return err
	}
	c.index = uint16(data[0])<<8 | uint16(data[1])
	c.IncrementPC(2)
	return nil
}

func (c *CPU) yieldCoroutine() error {
	return nil
}

// func (c *CPU) jump(instruction uint16) error {
// 	c.SetPC(instruction & 0x0FFF)
// 	return nil
// }
func (c *CPU) callSubroutine(instruction uint16) error {
	ok := c.stack.Push(c.pc)
	if !ok {
		return fmt.Errorf("Stack overflow")
	}
	c.SetPC(instruction & 0x0FFF)
	return nil
}
func (c *CPU) returnFromSubroutine() error {
	d, ok := c.stack.Pop()
	if !ok {
		return fmt.Errorf("Stack is empty")
	}
	c.SetPC(d)
	return nil
}
