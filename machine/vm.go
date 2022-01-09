package machine

import (
	"encoding/binary"
	"fmt"
	"time"
)

type Counters map[rune]uint16
type Registers map[uint8]byte

// Ch8p is the CHIP-8 machine itself
type Ch8p struct {
	V        Registers
	Counters Counters
	Stack    Stack
	GFX      Memory
	RAM      Memory
	Keyboard Memory
	Delay    *time.Ticker
	Sound    *time.Ticker
	DrawFlag bool
	LastOp   string `default:"none"`
	pipeline chan Op
}

func (c *Ch8p) Tick() {
	opcode := c.ReadInstruction()
	c.DrawFlag = false
	op := NewOp(opcode)
	if op.Code == 0x0000 {
		fmt.Printf("Opcode: %X\n", opcode)
	} else {
		c.LastOp = fmt.Sprintf("%v\n", op) + c.LastOp
		op.Execute(c, op)
	}
	c.IncrementCounter('T')
}

// LoadFonts will put each of the fonts in Fonts into memory
func (c *Ch8p) LoadFonts() {
	for i, font := range Fonts {
		c.RAM.WriteBytes(uint16(i*5), font[:])
	}
	c.LastOp = ""
}

func (c *Ch8p) DrawSprite(x, y uint16, height uint16) {
	sprite := c.ReadRAMBytes(c.ReadCounter('I'), height)
	for yPos, b := range sprite {
		for xPos := uint16(0); xPos < 8; xPos++ {
			if x + xPos >= 64 || y + uint16(yPos) >= 32 {
				continue
			}
			pos := (y + uint16(yPos)) * 64 + x + xPos
			onScreen := c.GFX.ReadByte(pos)
			toBe := b & (0x80 >> xPos)
			if toBe != 0 && onScreen != 0 {
				c.GFX.WriteByte(pos, 0)
				c.WriteRegister('F', 1)
			} else if toBe != 0 && onScreen == 0 {
				c.GFX.WriteByte(pos, 1)
			}
		}
	}
	c.WriteCounter('I', 0)
	c.DrawFlag = true
}

func (c *Ch8p) ClearScreen() {
	for i := 0; i < len(c.GFX); i++ {
		c.GFX.WriteByte(uint16(i), 0)
	}
	c.DrawFlag = true
}

func (c *Ch8p) IncrementProgramCounter() uint16 {
	pc := c.ReadCounter('P')
	c.WriteCounter('P', pc+2)
	return pc
}

func (c *Ch8p) ReadInstruction() uint16 {
	pc := c.IncrementProgramCounter()
	opbytes := c.ReadRAMBytes(pc, 2)
	opcode := binary.BigEndian.Uint16(opbytes)
	return opcode
}

// ReadRAM does what it says on the tin
func (c *Ch8p) ReadRAM(addr uint16) byte {
	return c.RAM.ReadByte(addr)
}
func (c *Ch8p) ReadRAMBytes(addr uint16, length uint16) []byte {
	return c.RAM.ReadBytes(addr, length)
}
// WriteRAM does what it says on the tin
func (c *Ch8p) WriteRAM(addr uint16, value byte) {
	c.RAM.WriteByte(addr, value)
}
// WriteRAMBytes does WriteRAM but for a slice of bytes
func (c *Ch8p) WriteRAMBytes(addr uint16, bytes []byte) {
	c.RAM.WriteBytes(addr, bytes)
}

// ReadRegister does what it says on the tin
func (c *Ch8p) ReadRegister(reg uint8) byte {
	return c.V[reg]
}
// WriteRegister does what it says on the tin
func (c *Ch8p) WriteRegister(reg uint8, value byte) {
	c.V[reg] = value
}
// ReadCounter does what it says on the tin
func (c *Ch8p) ReadCounter(reg rune) uint16 {
	return c.Counters[reg]
}
// WriteCounter does what it says on the tin
func (c *Ch8p) WriteCounter(reg rune, value uint16) {
	c.Counters[reg] = value
}
// IncrementCounter does what it says on the tin
func (c *Ch8p) IncrementCounter(reg rune) {
	c.Counters[reg]++
}

// func (c *Ch8p) Debug() {
// 	tm.Clear()

// 	tm.MoveCursor(1, 1)

// 	tm.Printf(`opcode: %x
// pc: %d
// sp: %d
// i: %d
// ---Registers---
// V0: %d
// V1: %d
// V2: %d
// V3: %d
// V4: %d
// V5: %d
// V6: %d
// V7: %d
// V8: %d
// V9: %d
// VA: %d
// VB: %d
// VC: %d
// VD: %d
// VE: %d
// VF: %d\n`,
// 		0x00, c.Counters['P'], c.Counters['S'], c.Counters['I'],
// 		c.V.Get('0'), c.V.Get('1'), c.V.Get('2'), c.V.Get('3'),
// 		c.V.Get('4'), c.V.Get('5'), c.V.Get('6'), c.V.Get('7'),
// 		c.V.Get('8'), c.V.Get('9'), c.V.Get('A'), c.V.Get('B'),
// 		c.V.Get('C'), c.V.Get('D'), c.V.Get('E'), c.V.Get('F'),
// 	)
// 	tm.Flush()
// }
