package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU_ExecuteInstruction(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	assert.Equal(t, len(cpu.opcodes), 4)
	cpu.ram.Writes(0x0, []byte{0x00, 0xE0})
	cpu.stack.Push(0x200) // Dummy stack value so 0x00EE doesn't fail
	for _, v := range []uint16{0x00E0, 0x00EC, 0x00EE, 0xA000} {
		err := cpu.ExecuteInstruction(v)
		if err != nil {
			t.Errorf("Expected no error for 0x%X, got %v", v, err)
		}
	}
	for _, v := range []uint16{0x00FF, 0x1000, 0x6000, 0x7000, 0xD000, 0xF000} {
		err := cpu.ExecuteInstruction(v)
		assert.IsType(t, InstructionUnknown{}, err)
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