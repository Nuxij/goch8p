package mmu

import "fmt"

// 0x0000 - 0x0080 - font set
// 0x00A0 - 0x0E8F - unused
// 0x0E90 - 0x0FFF - 352 bytes for registers, stack, frame buffer

// 4096 bytes of memory gives us 16 pages of 256 bytes each
// 0x0000 - 0x00FF - page 0
// 0x0100 - 0x01FF - page 1
// 0x0200 - 0x02FF - page 2
// 0x0300 - 0x03FF - page 3
// 0x0400 - 0x04FF - page 4
// 0x0500 - 0x05FF - page 5
// 0x0600 - 0x06FF - page 6
// 0x0700 - 0x07FF - page 7
// 0x0800 - 0x08FF - page 8
// 0x0900 - 0x09FF - page 9
// 0x0A00 - 0x0AFF - page 10
// 0x0B00 - 0x0BFF - page 11
// 0x0C00 - 0x0CFF - page 12
// 0x0D00 - 0x0DFF - page 13
// 0x0E00 - 0x0EFF - page 14
// 0x0F00 - 0x0FFF - page 15

type PageFault struct {
	addr uint16
}

func (e PageFault) Error() string {
	return fmt.Sprintf("Page fault at address: %x", e.addr)
}

// SplitAddress splits an address into the address and options
func SplitAddress(addr uint16) (uint16, uint8) {
	return addr & 0x0FFF, uint8(addr >> 12) & 0x0F
}

// ConvertAddress converts an address to a page number and offset
// ignores the options (highest 4 bits)
func ConvertAddress(addr uint16) (uint8, uint8) {
	return uint8(addr >> 8) & 0x0F, uint8(addr & 0xFF)
}