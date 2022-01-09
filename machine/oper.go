package machine

import "fmt"

// Op will store a code and Execute via an Oper
type Op struct {
	Oper
	Code uint16
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
// N returns the byte in position N of the opcode (the lowest byte)
func (o Op) N() byte {
	return byte(o.Code & 0x000F)
}
// NN returns the lowest two bytes of the opcode
func (o Op) NN() uint16 {
	return o.Code & 0x00FF
}
// NNN returns everything but the opcode itself (the lowest three bytes)
func (o Op) NNN() uint16 {
	return o.Code & 0x0FFF
}

// NewOp returns a new Op with the given code.
func NewOp(opcode uint16) Op {
	switch opcode & 0xF000 >> 12 {
	case 0x0:
		return Op{Oper: OperSys{}, Code: opcode}
	case 0x1:
		return Op{Oper: OperJump{}, Code: opcode}
	case 0x2:
		return Op{Oper: OperCall{}, Code: opcode}
	case 0x3:
		return Op{Oper: OperSE{}, Code: opcode}
	case 0x4:
		return Op{Oper: OperSNE{}, Code: opcode}
	case 0x5:
		return Op{Oper: OperSE{}, Code: opcode}
	case 0x6:
		return Op{Oper: OperLD{}, Code: opcode}
	case 0x7:
		return Op{Oper: OperADD{}, Code: opcode}
	case 0x8:
		return Op{Oper: OperBit{}, Code: opcode}
	case 0x9:
		return Op{Oper: OperSNE{}, Code: opcode}
	case 0xA:
		return Op{Oper: OperLDI{}, Code: opcode}
	case 0xB:
		return Op{Oper: OperJPV0{}, Code: opcode}
	case 0xC:
		return Op{Oper: OperRND{}, Code: opcode}
	case 0xD:
		return Op{Oper: OperDRW{}, Code: opcode}
	case 0xE:
		return Op{Oper: OperSKP{}, Code: opcode}
	case 0xF:
		return Op{Oper: OperSpecial{}, Code: opcode}
	default:
		fmt.Printf("Unknown opcode: %X\n", opcode)
		return Op{Oper: OperSys{}, Code: 0x0000}
	}
}