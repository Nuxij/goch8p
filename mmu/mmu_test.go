package mmu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestConvertAddress(t *testing.T) {
	page, offset := ConvertAddress(0x200)
	assert.Equal(t, uint8(0x2), page)
	assert.Equal(t, uint8(0x0), offset)

	page, offset = ConvertAddress(0x4EE)
	assert.Equal(t, uint8(0x4), page)
	assert.Equal(t, uint8(0xEE), offset)

	page, offset = ConvertAddress(0xDFEE)
	assert.Equal(t, uint8(0xF), page)
	assert.Equal(t, uint8(0xEE), offset)
}

// Running program sees:
//	0x0000: 0xF0 0x90 0x90 0x90 0xF0 0x00 0x00 0x00
//	0x0200: 0x20 0x60 0x20 0x20 0x70 0x00 0x00 0x00
//  ...
//	0x0E00: 0xF0 0x10 0xF0 0x80 0xF0 0x00 0x00 0x00

// CPU sees in memory:
//	0x0000: 0xF0 0x90 0x90 0x90 0xF0 0x00 0x00 0x00
//  0x01F0: 0x00 0x21 0x00 0x00 0x00 0x00 0x00 0x00
//  0x01F8 

func TestReadPageTable(t *testing.T) {
	ram := NewRAM(0x1000) // 4KB
	rammer := NewRammer(map[uint16]MemoryDevice{
		0x0000: ram,
	})
	rammer.Write(0x0000, 0xFF)
	rammer.Write(0x0100, 0x01F0) // 0x01F0 is the first page table (0x200 - (0x8*16)*1)

	mmu := NewMMU(rammer)
	mmu.SetProcessTable(map[uint8]uint16{
		0: 0x0100,
	})
	
	// Set up page table

	pt := &PageTable{0x01F0}
	frame, fault := pt.Read(0x200)
	assert.Equal(t, uint8(0x0), frame)
	assert.ErrorIs(t, fault, PageFault{0x200})
}