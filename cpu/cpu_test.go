package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU_ExecuteInstruction(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	for _, v := range []uint16{0x00E0, 0x00EC, 0x00EE} {
		err := cpu.ExecuteInstruction(v)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}	
	err := cpu.ExecuteInstruction(0x00FF)
	assert.ErrorIs(t, err, InstructionUnknown{0xFF})
}

func TestCPU_ClearScreen(t *testing.T) {
	cpu := NewCPU(NewRAM(0x1000))
	cpu.ram.Writes(cpu.screen.address, DataForTest)
	_, err := cpu.ram.Reads(cpu.screen.address, cpu.screen.size)
	assert.NoError(t, err)
	cpu.ExecuteInstruction(0x00E0)
	for i := range DataForTest {
		data, err := cpu.ram.Read(cpu.screen.address+uint16(i))
		assert.NoError(t, err)
		assert.EqualValues(t, 0, data)
	}
}

// func TestDrawToScreen(t *testing.T) {
// 	cpu := NewCPU()
// 	opcode := 0xD005
// }