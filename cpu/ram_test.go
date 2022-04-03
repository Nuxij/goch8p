package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var DataForTest = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
var Size = uint16(0x1000)
var Last = uint16(0xFFF)
var Half = uint16(0x800)
func TestRAM_size_is_separate_from_data_length(t *testing.T) {
	size := uint16(len(DataForTest) + 512)
	r := NewRAM(size)
	r.data = DataForTest
	assert.EqualValues(t, len(DataForTest), len(r.data)) // obvious, may change as we implement r.Write()
	assert.EqualValues(t, size, r.size)
	assert.NotEqual(t, len(r.data), r.size)
}
func TestRAM_stores_history_in_a_buffer_when_cleared(t *testing.T) {
	size := uint16(len(DataForTest) + 512)
	r := NewRAM(size)
	r.data = DataForTest
	assert.EqualValues(t, 0, len(r.backbuffer))
	r.Clear()
	assert.EqualValues(t, 1, len(r.backbuffer))
	assert.EqualValues(t, DataForTest[0], r.backbuffer[len(r.backbuffer)-1][0x0])
}
func TestRAM_Clear(t *testing.T) {
	size := uint16(len(DataForTest) + 512)
	r := NewRAM(size)
	r.data = DataForTest
	r.Clear()
	for _, v := range r.data {
		assert.EqualValues(t, 0, v)
	}
}
func TestRAM_Reset(t *testing.T) {
	r := NewRAM(0x1000)
	for i := range DataForTest {
		r.data[int(Half)+i] = DataForTest[i]
	}
	assert.EqualValues(t, DataForTest, r.data[Half:Half+uint16(len(DataForTest))])
	r.Clear()
	assert.Equal(t, 1, len(r.backbuffer))
	assert.EqualValues(t, DataForTest[0], r.backbuffer[len(r.backbuffer)-1][Half])
	assert.EqualValues(t, 0, r.data[Half])
	r.Reset()
	assert.EqualValues(t, 0, r.data[Half])
	assert.Equal(t, 0, len(r.backbuffer))
}

func TestRAM_Read(t *testing.T) {
	ram := NewRAM(Size)
	ram.data[Last] = 0xFF
	ram.data[Half] = 0xAA

	read, err := ram.Read(Last)
	assert.EqualValues(t, 0xFF, read)
	assert.NoError(t, err)

	read, err = ram.Read(Half)
	assert.EqualValues(t, 0xAA, read)
	assert.NoError(t, err)

	read, err = ram.Read(Size)
	assert.EqualValues(t, 0, read)
	assert.IsType(t, AddressOutOfRange{}, err)
}

func TestRAM_Write(t *testing.T) {
	ram := NewRAM(Size)
	err := ram.Write(Half, 0xFF)
	assert.NoError(t, err)
	err = ram.Write(Last, 0xFF)
	assert.NoError(t, err)
	err = ram.Write(Size, 0xFF)
	assert.IsType(t, AddressOutOfRange{}, err)
}