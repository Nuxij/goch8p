package machine

import "fmt"

// Oper interface allows Ops to override execution behavior
type Oper interface {
	Name() string
	Execute(c *Ch8p, op OpCode)
}

type OpCode interface {
	Oper
	Op() uint16
	OpClass() byte
	X() byte
	Y() byte
	XY() (byte, byte)
	N() byte
	KK() uint16
	NNN() uint16
	Name() string
}

// Op will store a code and Execute via an Oper
type Op struct {
	Code uint16
	name string
	call func(c *Ch8p)
}
// Name prints the name to match Oper
func (o Op) Name() string {
	return o.name
}
// String implements the Stringer interface
func (o Op) String() string {
	return fmt.Sprintf("%X [%v]", o.Code, o.Name())
}
// Op returns the opcode only
func (o Op) Op() uint16 {
	return o.Code & 0xF000
}
// OpClass returns the opcode's highest byte only
func (o Op) OpClass() byte {
	return byte(o.Op() >> 12)
}
// X returns the byte in position X of the opcode
func (o Op) X() byte {
	return byte(o.Code & 0x0F00 >> 8)
}
// Y returns the byte in position Y of the opcode
func (o Op) Y() byte {
	return byte(o.Code & 0x00F0 >> 4)
}
// XY returns the X and Y bytes of the opcode
func (o Op) XY() (byte, byte) {
	return o.X(), o.Y()
}
// N returns the nibble in position N of the opcode (the lowest 4 bits)
func (o Op) N() byte {
	return byte(o.Code & 0x000F)
}
// KK returns the lowest byte of the opcode
func (o Op) KK() uint16 {
	return o.Code & 0x00FF
}
// NNN returns everything but the opcode itself (the lowest three bytes)
func (o Op) NNN() uint16 {
	return o.Code & 0x0FFF
}
