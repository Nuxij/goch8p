package machine

type Oper interface {
	Name() string
	Execute(c *Ch8p, op Op)
}

type OperSys struct{}

func (o OperSys) Name() string {
	return "Sys"
}
func (o OperSys) Execute(c *Ch8p, op Op) {
	switch op.KK() {
	case 0x00E0:
		c.ClearScreen()
	case 0x00EE:
		// c.ReturnFromSubroutine()
	}
}

type OperJump struct{}

func (o OperJump) Name() string {
	return "Jump"
}
func (o OperJump) Execute(c *Ch8p, op Op) {
	c.WriteCounter('P', op.NNN())
}

type OperCall struct{}

func (o OperCall) Name() string {
	return "Call"
}
func (o OperCall) Execute(c *Ch8p, op Op) {
	// c.CallSubroutine(op.NNN())
}

type OperSE struct{}

func (o OperSE) Name() string {
	return "SkipIfEqual"
}
func (o OperSE) Execute(c *Ch8p, op Op) {
	Vx := c.ReadRegister(op.X())
	if Vx == op.N() {
		c.IncrementProgramCounter()
	}
}

type OperSNE struct{}

func (o OperSNE) Name() string {
	return "SkipIfNotEqual"
}
func (o OperSNE) Execute(c *Ch8p, op Op) {
	Vx := c.ReadRegister(op.X())
	if Vx != op.N() {
		c.IncrementProgramCounter()
	}
}

type OperLD struct{}

func (o OperLD) Name() string {
	return "Load"
}
func (o OperLD) Execute(c *Ch8p, op Op) {
	c.WriteRegister(op.X(), op.N())
}

type OperADD struct{}

func (o OperADD) Name() string {
	return "Add"
}
func (o OperADD) Execute(c *Ch8p, op Op) {
	Vx := c.ReadRegister(op.X())
	Vy := c.ReadRegister(op.N())
	c.WriteRegister(op.X(), Vx+Vy)
}

type OperBit struct{}

func (o OperBit) Name() string {
	return "Bit"
}
func (o OperBit) Execute(c *Ch8p, op Op) {
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

type OperLDI struct{}

func (o OperLDI) Name() string {
	return "LoadI"
}
func (o OperLDI) Execute(c *Ch8p, op Op) {
	c.WriteCounter('I', op.NNN())
}

type OperJPV0 struct{}

func (o OperJPV0) Name() string {
	return "JumpV0"
}
func (o OperJPV0) Execute(c *Ch8p, op Op) {
	// c.JumpToAddress(op.NNN() + c.ReadRegister(0))
}

type OperRND struct{}

func (o OperRND) Name() string {
	return "Random"
}
func (o OperRND) Execute(c *Ch8p, op Op) {
	// c.WriteRegister(op.X(), c.Random(op.N()))
}

type OperDRW struct{}

func (o OperDRW) Name() string {
	return "Draw"
}
func (o OperDRW) Execute(c *Ch8p, op Op) {
	Vx := op.X()
	Vy := op.Y()
	X := uint16(c.ReadRegister(Vx) & 63)
	Y := uint16(c.ReadRegister(Vy) % 31)
	height := uint16(op.N())
	c.DrawSprite(X, Y, height)
}

type OperSKP struct{}

func (o OperSKP) Name() string {
	return "SkipIfKeyPressed"
}
func (o OperSKP) Execute(c *Ch8p, op Op) {
	// if c.KeyPressed(op.X()) {
	// 	c.IncrementProgramCounter()
	// }
}

type OperSpecial struct{}

func (o OperSpecial) Name() string {
	return "Special"
}
func (o OperSpecial) Execute(c *Ch8p, op Op) {
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