package mmu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


var ram = NewRAM(0x1000)
func TestNewRAM(t *testing.T) {
	assert.Equal(t, 0x1000, ram.Size())
}

func TestReadRAM(t *testing.T) {
	addr := uint16(0x00)
	value, err := ram.Read(addr)
	assert.NotErrorIs(t, err, AddressOutOfRange{addr})
	assert.Equal(t, byte(0xAA), value)

	value, err = ram.Read(0xFFF)
	assert.NotErrorIs(t, err, AddressOutOfRange{0xFFF})
	assert.Equal(t, byte(0xBB), value)

	value, err = ram.Read(0x1100)
	assert.ErrorIs(t, err, AddressOutOfRange{0x1100})
}

func TestWriteRAM(t *testing.T) {
	addr := uint16(0x00)
	value := byte(0xCC)
	err := ram.Write(addr, value)
	assert.NotErrorIs(t, err, AddressOutOfRange{addr})
	assert.Equal(t, value, ram.data[addr])

	addr = 0xFFF
	value = byte(0xDD)
	err = ram.Write(addr, value)
	assert.NotErrorIs(t, err, AddressOutOfRange{addr})
	assert.Equal(t, value, ram.data[addr])

	addr = 0x1100
	err = ram.Write(addr, value)
	assert.ErrorIs(t, err, AddressOutOfRange{addr})
}