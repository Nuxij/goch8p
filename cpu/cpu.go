package cpu

import (
	"fmt"

	"github.com/google/uuid"
)

type CPU struct {
	ram   Device
	pc    uint16
	sp    uint16
	stack *Stack
	screen Screen
}

func RandomStringUUID() string {
	id, _ := uuid.NewUUID()
	return fmt.Sprintf("%v", id)
}

type InstructionUnknown struct {
	opcode uint16
}
func (iu InstructionUnknown) Error() string {
	return fmt.Sprintf("Unknown instruction: %X", iu.opcode)
}

type Screen struct {
	width  uint16
	height uint16
	size  uint16
	address uint16
}

func NewCPU(ram Device) *CPU {
	addressableSize := uint16(0x1000)
	screenSize := uint16(64*32)/8
	screenAddress := addressableSize-(screenSize+1)
	rammer := NewRammer(256, []Device{ram})
	rammer.SetRegion(0x0, addressableSize, ram, 0x0)
	rammer.SetRegion(screenAddress, screenSize, ram, screenAddress)
	cpu := &CPU{
		ram: rammer,
		pc:  0x200,
		sp:  0x0,
		stack: NewStack(0x10),
		screen: Screen {
			width:  64,
			height: 32,
			size: screenSize,
			address: screenAddress,
		},
	}
	return cpu
}

func (c *CPU) ExecuteInstruction(instruction uint16) error {
	var err error
	switch instruction & 0x00FF {
	case 0x00E0:
		err = c.clearScreen()
	case 0x00EC:
		err = c.yieldCoroutine()
	case 0x00EE:
		err = c.returnFromSubroutine()
	default:
		return InstructionUnknown{instruction}
	}
	return err
}

func (c *CPU) clearScreen() error {
	for i := uint16(0); i < c.screen.size; i++ {
		c.ram.Write(c.screen.address+i, 0x0)
	}
	return nil
}
func (c *CPU) yieldCoroutine() error {
	return nil
}
func (c *CPU) returnFromSubroutine() error {
	return nil
}