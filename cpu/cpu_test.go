package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU_ExecuteInstruction(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	assert.Equal(t, len(cpu.opcodes), 1)

	cpu.ram.Writes(0x0, []byte{0x00, 0xE0}) // clearscreen in memory for testing read + exec
	cpu.stack.Push(0x200)                   // Dummy stack value so 0x00EE doesn't fail

	type op struct {
		code uint16
		err  error
	}
	var ops = []op{
		{0x00E0, InstructionOk{}},
		{0x00EC, nil},
		{0x00EE, nil},
		{0xA000, nil},
		{0x00FF, InstructionUnknown{}},
		{0x1000, InstructionUnknown{}},
		{0xD000, InstructionUnknown{}},
		{0xF000, InstructionUnknown{}},
		{0x7000, InstructionUnknown{}},
		{0x6000, InstructionUnknown{}},
	}

	for _, v := range ops {
		err := cpu.ExecuteInstruction(v.code)
		assert.IsTypef(t, v.err, err, "expecting %X to error with %v", v.code, v.err)
	}

}

func TestCPU_ClearScreen(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	cpu.ram.Writes(cpu.screen.address, DataForTest)
	cpu.ExecuteInstruction(0x00E0)
	for i := range DataForTest {
		data, err := cpu.ram.Read(cpu.screen.address + uint16(i))
		assert.NoError(t, err)
		assert.EqualValues(t, 0, data)
	}
}

// func TestDrawToScreen(t *testing.T) {
// 	cpu := NewCPU()
// 	opcode := 0xD005
// }

func TestCPU_loadIndex(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	opcode := uint16(0xA000)
	target := uint16(0x0300)
	cpu.ram.Writes(target, []byte{0xFF, 0xFF})
	err := cpu.ExecuteInstruction(opcode | target)
	assert.NoError(t, err)
	assert.EqualValues(t, 0xFFFF, cpu.index)
}

func TestCPU_callSubroutine(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	opcode := uint16(0x2000)
	target := uint16(0x0A00)
	startPC := cpu.pc
	cpu.ExecuteInstruction(opcode | target)
	poppedPC, ok := cpu.stack.Pop()
	assert.True(t, ok)
	assert.EqualValues(t, startPC, poppedPC)
	assert.EqualValues(t, target, cpu.pc)
}
