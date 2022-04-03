package mmu

import "fmt"

type AddressOutOfRange struct {
	addr uint16
}

func (e AddressOutOfRange) Error() string {
	return "Address out of range: " + fmt.Sprint(e.addr)
}

type RAM struct {
	data []byte
	size uint16
}

func (r *RAM) Size() uint16 {
	return r.size
}
func (r *RAM) CheckBounds(addr uint16) (uint16, error) {
	if addr >= r.Size() {
		return 0, AddressOutOfRange{addr}
	}
	return addr, nil
}
func (r *RAM) Read(addr uint16) (byte, error) {
	baddr, error := r.CheckBounds(addr)
	if error != nil {
		return 0, error
	}
	return r.data[baddr], nil
}
func (r *RAM) Write(addr uint16, value byte) error {
	baddr, err := r.CheckBounds(addr)
	if err != nil {
		return err
	}
	r.data[baddr] = value
	return nil
}
func NewRAM(size uint16) *RAM {
	ram := &RAM{}
	ram.data[0x00] = 0xAA
	ram.data[0xFFF] = 0xFF
	return ram
}