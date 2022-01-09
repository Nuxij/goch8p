package machine

import (
	"encoding/binary"
	"fmt"
	"time"
)

type Counters map[rune]uint16
type Register byte
type Registers map[uint8]Register

func NewRegisterBank() Registers {
	return make(Registers, 8)
}

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
	drawFlag bool
	LastOp   string `default:"none"`
}

func (c *Ch8p) Tick() {
	opcode := c.ReadInstruction()
	c.drawFlag = false
	c.ParseInstruction(opcode)
}

// LoadFonts will put each of the fonts in Fonts into memory
func (c *Ch8p) LoadFonts() {
	for i, font := range Fonts {
		c.RAM.WriteBytes(uint16(i*5), font[:])
	}
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
	c.drawFlag = true
}

func (c *Ch8p) ClearScreen() {
	for i := 0; i < len(c.GFX); i++ {
		c.GFX.WriteByte(uint16(i), 0)
	}
	c.drawFlag = true
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

func (c *Ch8p) ParseInstruction(opcode uint16) {
	instruction := (opcode & 0xF000)
	lastOp := fmt.Sprintf("[%X, %X]", opcode, instruction)
	switch instruction {
	case 0xD000: // DXYN: Draw sprite at (VX, VY) with N bytes of sprite data
		Vx := uint8(opcode & 0x0F00 >> 8)
		Vy := uint8(opcode & 0x00F0 >> 4)
		X := uint16(c.ReadRegister(Vx) & 63)
		Y := uint16(c.ReadRegister(Vy) % 31)
		height := uint16(opcode & 0x000F)

		c.DrawSprite(X, Y, height)
		lastOp += " 0xDXYN"
	case 0xA000: // ANNN: Sets I to the address NNN.
		c.WriteCounter('I', opcode&0x0FFF)
		lastOp += " 0xANNN"
	case 0x6000: // 6XNN: Sets VX to NN.
		Vx := uint8(opcode & 0x0F00 >> 8)
		NN := opcode & 0x00FF
		c.WriteRegister(Vx, Register(NN))
		lastOp += " 0x6XNN"
	}
	lastOp += "\n"
	if c.ReadCounter('T') % 1000 == 0 {
		c.LastOp = lastOp
	} else {
		c.LastOp = lastOp + c.LastOp 
	}
	c.IncrementCounter('T')
}

var opcodes = map[uint16]string {
	0x00E0: "CLS",
	0x00EE: "RET",
	0x1000: "1NNN",
	0x2000: "2NNN",
	0x3000: "3XNN",
	0x4000: "4XNN",
	0x5000: "5XY0",
	0x6000: "6XNN",
	0x7000: "7XNN",
	0x8000: "8XY0",
	0x8001: "8XY1",
	0x8002: "8XY2",
	0x8003: "8XY3",
	0x8004: "8XY4",
	0x8005: "8XY5",
	0x8006: "8XY6",
	0x8007: "8XY7",
	0x800E: "8XYE",
	0x9000: "9XY0",
	0xA000: "ANNN",
	0xB000: "BNNN",
	0xC000: "CXNN",
	0xD000: "DXYN",
	0xE09E: "EX9E",
	0xE0A1: "EXA1",
	0xF007: "FX07",
	0xF00A: "FX0A",
	0xF015: "FX15",
	0xF018: "FX18",
	0xF01E: "FX1E",
	0xF029: "FX29",
	0xF033: "FX33",
	0xF055: "FX55",
	0xF065: "FX65",
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
func (c *Ch8p) ReadRegister(reg uint8) Register {
	return c.V[reg]
}
// WriteRegister does what it says on the tin
func (c *Ch8p) WriteRegister(reg uint8, value Register) {
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
