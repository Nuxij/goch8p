package machine

import "fmt"

// Memory is a byte array uint16 class
type Memory []byte

func (m Memory) ReadByte(addr uint16) byte {
	return m[addr]
}
func (m Memory) ReadWord(addr uint16) uint16 {
	return (uint16(m[addr]) << 8) | uint16(m[addr+1])
}
func (m Memory) ReadBytes(addr uint16, length uint16) []byte {
	return m[addr : addr+length]
}
func (m Memory) WriteByte(addr uint16, value byte) {
	m[addr] = value
}
func (m Memory) WriteWord(addr uint16, value uint16) {
	m[addr] = byte(value >> 8)
	m[addr+1] = byte(value)
}
func (m Memory) WriteBytes(addr uint16, bytes []byte) {
	for i, b := range bytes {
		m[addr+uint16(i)] = b
	}
}
func (m Memory) String() string {
	return fmt.Sprintf("%v", m)
}

// type ByteMemory []byte
// type UintMemory []uint16

// func (m *ByteMemory) Clear() {
// 	*m = make([]byte, len(*m))
// }

// func (m *ByteMemory) WriteBytes(addr Int, bytes []byte) {
// 	for i, b := range bytes {
// 		m.WriteByte(addr+Int(i), b)
// 	}
// }

// func (m *ByteMemory) WriteByte(addr Int, value byte) {
// 	if addr < 0 && addr >= Int(len(*m)) {
// 		panic("Trying to write outside of memory")
// 	}
// 	// if addr < 0x200 {
// 	// 	panic("Trying to write to ROM")
// 	// }
// 	(*m)[addr] = value
// }

// func (m *ByteMemory) ReadByte(addr Int) byte {
// 	if addr < 0 && addr >= Int(len(*m)) {
// 		panic("Trying to read outside of memory")
// 	}
// 	return (*m)[addr]
// }

// func (m *ByteMemory) ReadBytes(addr Int, length int) []byte {
// 	if addr < 0 && addr >= Int(len(*m)) {
// 		panic("Trying to read outside of memory")
// 	}
// 	return (*m)[addr : addr+Int(length)]
// }
