package cpu

import "fmt"

type AddressOutOfRange struct {
	size uint16
	addr uint16
}

func (a AddressOutOfRange) Error() string {
	return fmt.Sprintf("Address out of range: %X, max: %X", a.addr, a.size)
}

type RAM struct {
	uuid string
	backbuffer [][]byte
	data       []byte
	size       uint16
}

func NewRAM(size uint16) *RAM {
	r := &RAM{
		uuid: fmt.Sprintf("RAM::%s", RandomStringUUID()),
		size: size,
	}
	r.clearData(true)
	return r
}

func (r *RAM) UUID() string {
	return r.uuid
}

func (r *RAM) Read(addr uint16) (byte, error) {
	if err := r.checkBounds(addr); err != nil {
		return 0, err
	}
	return r.data[addr], nil
}

func (r *RAM) Reads(addr uint16, size uint16) ([]byte, error) {
	if err := r.checkBounds(addr); err != nil {
		return nil, err
	}
	if err := r.checkBounds(addr + size); err != nil {
		return nil, err
	}
	return r.data[addr : addr+size], nil
}

func (r *RAM) Write(addr uint16, value byte) error {
	if err := r.checkBounds(addr); err != nil {
		return err
	}
	r.data[addr] = value
	return nil
}

func (r *RAM) Writes(addr uint16, data []byte) error {
	if err := r.checkBounds(addr); err != nil {
		return err
	}
	if err := r.checkBounds(addr + uint16(len(data))); err != nil {
		return err
	}
	for i, b := range data {
		r.data[addr+uint16(i)] = b
	}
	return nil
}

func (r *RAM) checkBounds(addr uint16) error {
	if addr >= r.size {
		return AddressOutOfRange{r.size, addr}
	}
	return nil
}

func (r *RAM) flushData() {
	r.backbuffer = append(r.backbuffer, r.data)
	r.clearData(false)
}
func (r *RAM) clearData(buffer bool) {
	r.data = make([]byte, r.size)
	if buffer {
		r.clearBackbuffer()
	}
}
func (r *RAM) clearBackbuffer() {
	r.backbuffer = make([][]byte, 0)
}

func (r *RAM) Clear() {
	r.flushData()
	for i := range r.data {
		r.data[i] = 0
	}
}

func (r *RAM) Reset() {
	r.clearData(true)
}
