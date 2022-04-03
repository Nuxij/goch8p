package machine

import "fmt"

// NewOp calls NewOper and returns it cast to an OpCode
func NewOp(opcode uint16) OpCode {
	return NewOper(opcode).(OpCode)
}
// NewOper returns a new Oper with the given code.
func NewOper(opcode uint16) Oper {
	op := Op{Code: opcode}
	var oper Oper
	switch op.OpClass() {
	case 0x0:
		op.name = "Sys"
		oper = OperSys{op}
	case 0x1:
		op.name = "Jump"
		oper = OperJump{op}
	case 0x2:
		op.name = "Call"
		oper = OperCall{op}
	case 0x3:
		op.name = "SkipEqual"
		oper = OperSE{op}
	case 0x4:
		op.name = "SkipNotEqual"
		oper = OperSNE{op}
	case 0x5:
		op.name = "SkipEqual:XY"
		oper = OperSEXY{op}
	case 0x6:
		op.name = "Load"
		oper = OperLD{op}
	case 0x7:
		op.name = "Add"
		oper = OperADD{op}
	case 0x8:
		op.name = "Bit"
		oper = OperBit{op}
	case 0x9:
		op.name = "SkipNotEqual:XY"
		oper = OperSNEXY{op}
	case 0xA:
		op.name = "Load:I"
		oper = OperLDI{op}
	case 0xB:
		op.name = "JumpV0"
		oper = OperJPV0{op}
	case 0xC:
		op.name = "RandomByte"
		oper = OperRND{op}
	case 0xD:
		op.name = "Draw"
		oper = OperDRW{op}
	case 0xE:
		switch op.KK() {
		case 0x9E:
			op.name = "SkipIfPressed"
			oper = OperSKP{op}
		case 0xA1:
			op.name = "SkipIfNotPressed"
			oper = OperSKNP{op}
		}
	case 0xF:
		op.name = "Special"
		oper = OperSpecial{op}
	default:
	fmt.Printf("Unknown opcode: %X\n", opcode)
	op.name = "Unknown"
	op.Code = 0x0
	oper = OperSys{op}
	}
	oper.(OpCode).call = func(c *Ch8p) {
		
	}
	return oper
}

// OperSys are the system instructions.
type OperSys struct{ Op }
// Execute the op
func (o OperSys) Execute(c *Ch8p, op OpCode) {
	switch op.KK() {
	case 0x00E0:
		c.ClearScreen()
	case 0x00EE:
		// c.ReturnFromSubroutine()
	}
}

// OperJump sets the program counter to the address in the operand.
type OperJump struct{ Op }
// Execute the op
func (o OperJump) Execute(c *Ch8p, op OpCode) {
	c.WriteCounter('P', op.NNN())
}

// OperCall calls a subroutine (likely by pushing to the stack and program counter).
type OperCall struct{ Op }
// Execute the op
func (o OperCall) Execute(c *Ch8p, op OpCode) {
	// c.CallSubroutine(op.NNN())
}

// OperSE is the skip if equal instruction.
type OperSE struct{ Op }
// Execute the op
func (o OperSE) Execute(c *Ch8p, op OpCode) {
	Vx := c.ReadRegister(op.X())
	if Vx == op.N() {
		c.IncrementProgramCounter()
	}
}

// OperSNE is the skip if not equal instruction.
type OperSNE struct{ Op }
// Execute the op
func (o OperSNE) Execute(c *Ch8p, op OpCode) {
	Vx := c.ReadRegister(op.X())
	if Vx != op.N() {
		c.IncrementProgramCounter()
	}
}

// OperSEXY will skip if X and Y registers are equal
type OperSEXY struct{ Op }
// Execute the op
func (o OperSEXY) Execute(c *Ch8p, op OpCode) {
	Vx := c.ReadRegister(op.X())
	if Vx != op.N() {
		c.IncrementProgramCounter()
	}
}

// OperSNEXY will skip if X and Y registers are not equal
type OperSNEXY struct{ Op }
// Execute the op
func (o OperSNEXY) Execute(c *Ch8p, op OpCode) {
	Vx := c.ReadRegister(op.X())
	if Vx != op.N() {
		c.IncrementProgramCounter()
	}
}

// OperLD will load the value in the operand into the register.
type OperLD struct{ Op }
// Execute the op
func (o OperLD) Execute(c *Ch8p, op OpCode) {
	c.WriteRegister(op.X(), op.N())
}

// OperADD adds the value in the operand to the register.
type OperADD struct{ Op }
// Execute the op
func (o OperADD) Execute(c *Ch8p, op OpCode) {
	Vx := c.ReadRegister(op.X())
	Vy := c.ReadRegister(op.N())
	c.WriteRegister(op.X(), Vx+Vy)
}

// OperBit will handle BitWise operations and other related operations.
type OperBit struct{ Op }
// Execute the op
func (o OperBit) Execute(c *Ch8p, op OpCode) {
	switch op.N() {
	case 0x0:
		c.WriteRegister(op.X(), c.ReadRegister(op.Y()))
		// case 0x1:
		// 	c.Or(o.X(), o.Y())
		// case 0x2:
		// 	c.And(o.X(), o.Y())
		// case 0x3:
		// 	c.Xor(o.X(), o.Y())
		// case 0x4:
		// 	c.AddToRegister(o.X(), o.Y())
		// case 0x5:
		// 	c.Subtract(o.X(), o.Y())
		// case 0x6:
		// 	c.ShiftRight(o.X())
		// case 0x7:
		// 	c.Subtract(o.Y(), o.X())
		// case 0xE:
		// 	c.ShiftLeft(o.X())
	}
}

// OperLDI will load the value in the operand into the 'I' counter.
type OperLDI struct{ Op }
// Execute the op
func (o OperLDI) Execute(c *Ch8p, op OpCode) {
	c.WriteCounter('I', op.NNN())
}

// OperJPV0 will jump to the address in the operand plus the value in the 'V0' register.
type OperJPV0 struct{ Op }
// Execute the op
func (o OperJPV0) Execute(c *Ch8p, op OpCode) {
	// c.JumpToAddress(op.NNN() + c.ReadRegister(0))
}

// OperRND will set the register to a random number between 0 and the operand.
type OperRND struct{ Op }
// Execute the op
func (o OperRND) Execute(c *Ch8p, op OpCode) {
	// c.WriteRegister(op.X(), c.Random(op.N()))
}

// OperDRW will draw a sprite at the x,y position from X,Y registers in the operand.
type OperDRW struct{ Op }
// Execute the op
func (o OperDRW) Execute(c *Ch8p, op OpCode) {
	Vx := op.X()
	Vy := op.Y()
	X := uint16(c.ReadRegister(Vx) & 63)
	Y := uint16(c.ReadRegister(Vy) % 31)
	height := uint16(op.N())
	c.DrawSprite(X, Y, height)
}

// OperSKP will skip the next instruction if the key in the operand is pressed.
type OperSKP struct{ Op }
// Execute the op
func (o OperSKP) Execute(c *Ch8p, op OpCode) {
	// if c.KeyPressed(op.X()) {
	// 	c.IncrementProgramCounter()
	// }
}

// OperSKNP will skip the next instruction if the key in the operand is not pressed.
type OperSKNP struct{ Op }
// Execute the op
func (o OperSKNP) Execute(c *Ch8p, op OpCode) {
	// if c.KeyPressed(op.X()) {
	// 	c.IncrementProgramCounter()
	// }
}

// OperSpecial will handle special operations like delay, sound + input.
type OperSpecial struct{ Op }
// Execute the op
func (o OperSpecial) Execute(c *Ch8p, op OpCode) {
	// if o.KK() == 0x07 {
	// 	c.ReadDelayTimer(o.X())
	// }
	// if o.KK() == 0x0A {
	// 	c.WaitForKey(o.X())
	// }
	// if o.KK() == 0x15 {
	// 	c.LoadDelayTimer(o.X())
	// }
	// if o.KK() == 0x18 {
	// 	c.LoadSoundTimer(o.X())
	// }
	// if o.KK() == 0x1E {
	// 	c.AddToI(o.X())
	// }
	// if o.KK() == 0x29 {
	// 	c.LoadSprite(o.X())
	// }
	// if o.KK() == 0x33 {
	// 	c.Bcd(o.X())
	// }
	// if o.KK() == 0x55 {
	// 	c.StoreRegisters(o.X())
	// }
	// if o.KK() == 0x65 {
	// 	c.LoadRegisters(o.X())
	// }

}