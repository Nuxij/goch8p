package main

import "github.com/Nuxij/goch8p/mmu"

func main() {
	realmemoryblock := mmu.NewByteMemory(0x1000)
	mapped := &mmu.MappedMemory{}
	
	mapped.AddBlock(0x0, realmemoryblock)
	mapped.Read(0x1)
}