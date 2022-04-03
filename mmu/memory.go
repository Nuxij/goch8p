package mmu

type MemoryDevice interface {
	Size() uint16
}

type ReadOnlyMemory interface {
	MemoryDevice
	Read(addr uint16) byte
	ReadWord(addr uint16) uint16
	ReadBytes(addr uint16, length uint16) []byte
}

type WriteOnlyMemory interface {
	MemoryDevice
	Write(addr uint16, value byte)
	WriteWord(addr uint16, value uint16)
	WriteBytes(addr uint16, bytes []byte)
}

type ReadWriteMemory interface {
	ReadOnlyMemory
	WriteOnlyMemory
}