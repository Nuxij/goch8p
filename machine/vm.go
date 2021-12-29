package machine

import (
	"fmt"
	"strconv"
	"time"
)

type Counters map[rune]uint16
type Register byte
type Registers map[uint8]Register

func NewRegisterBank() Registers {
	return make(Registers, 8)
}

type HexSprite [5]byte
type Stack [16]uint16

var HexSprites = [16]HexSprite{
	HexSprite{0xF0, 0x90, 0x90, 0x90, 0xF0}, // 0
	HexSprite{0x20, 0x60, 0x20, 0x20, 0x70}, // 1
	HexSprite{0xF0, 0x10, 0xF0, 0x80, 0xF0}, // 2
	HexSprite{0xF0, 0x10, 0xF0, 0x10, 0xF0}, // 3
	HexSprite{0x90, 0x90, 0xF0, 0x10, 0x10}, // 4
	HexSprite{0xF0, 0x80, 0xF0, 0x10, 0xF0}, // 5
	HexSprite{0xF0, 0x80, 0xF0, 0x90, 0xF0}, // 6
	HexSprite{0xF0, 0x10, 0x20, 0x40, 0x40}, // 7
	HexSprite{0xF0, 0x90, 0xF0, 0x90, 0xF0}, // 8
	HexSprite{0xF0, 0x90, 0xF0, 0x10, 0xF0}, // 9
	HexSprite{0xF0, 0x90, 0xF0, 0x90, 0x90}, // A
	HexSprite{0xE0, 0x90, 0xE0, 0x90, 0xE0}, // B
	HexSprite{0xF0, 0x80, 0x80, 0x80, 0xF0}, // C
	HexSprite{0xE0, 0x90, 0x90, 0x90, 0xE0}, // D
	HexSprite{0xF0, 0x80, 0xF0, 0x80, 0xF0}, // E
	HexSprite{0xF0, 0x80, 0xF0, 0x80, 0x80}, // F
}

type Ch8p struct {
	RAM      Memory
	V        Registers
	Counters Counters
	Stack    Stack
	GFX      Memory
	Keyboard Memory
	Delay    *time.Ticker
	Sound    *time.Ticker
	drawFlag bool
	LastOp   string `default:"none"`
}

func (c *Ch8p) Read(addr uint16) byte {
	return c.RAM.ReadByte(addr)
}

func (c *Ch8p) Write(addr uint16, value byte) {
	c.RAM.WriteByte(addr, value)
}
func (c *Ch8p) WriteBytes(addr uint16, bytes []byte) {
	c.RAM.WriteBytes(addr, bytes)
}

func (c *Ch8p) ReadRegister(reg uint8) Register {
	return c.V[reg]
}
func (c *Ch8p) WriteRegister(reg uint8, value Register) {
	c.V[reg] = value
}

func (c *Ch8p) ReadCounter(reg rune) uint16 {
	return c.Counters[reg]
}
func (c *Ch8p) WriteCounter(reg rune, value uint16) {
	c.Counters[reg] = value
}

func (c *Ch8p) DrawSprite(x, y uint16, height uint16) {
	// Draws a sprite at (x, y) with height N
	// The sprite is stored in memory starting at I
	// I is incremented after each byte is read
	// The sprite is drawn starting at the current position
	// If a pixel is set, the screen is set to 1
	// If a pixel is not set, the screen is set to 0
	// The screen is updated after each byte is read
	for yPos := uint16(0); yPos < height; yPos++ {
		spriteByte := c.RAM.ReadByte(c.ReadCounter('I') + yPos)
		for xPos := uint16(0); xPos < 8; xPos++ {
			if (spriteByte & (0x80 >> xPos)) != 0 {
				c.GFX.WriteByte(x+xPos+((y+yPos)*64), 1)
			}
		}
	}
	c.drawFlag = true
}

func (c *Ch8p) Tick() {
	pc := c.ReadCounter('P')
	smallestByte := c.RAM.ReadByte(pc)
	largestByte := c.RAM.ReadByte(pc + 1)
	opcode := uint16(largestByte)<<8 | uint16(smallestByte)
	instruction := (opcode & 0xF000)
	c.drawFlag = false

	switch instruction {
	case 0xD000: // DXYN: Draw sprite at (VX, VY) with N bytes of sprite data
		Vx := uint8(opcode & 0x0F00 >> 8)
		Vy := uint8(opcode & 0x00F0 >> 4)
		N := opcode & 0x000F // height of sprite
		c.DrawSprite(uint16(c.ReadRegister(Vx)), uint16(c.ReadRegister(Vy)), N)
		c.WriteCounter('P', pc+2)
		c.LastOp = "0xDXYN"
	case 0xA000: // ANNN: Sets I to the address NNN.
		c.WriteCounter('I', opcode&0x0FFF)
		c.WriteCounter('P', pc+2)
		c.LastOp = "0xANNN"
	default:
		c.LastOp = fmt.Sprintf("0x%v (%d)", strconv.FormatInt(int64(opcode), 16), strconv.FormatInt(int64(instruction), 10))
	}
	c.WriteCounter('T', c.ReadCounter('T')+1)

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
