package main

import (
	"fmt"
	"time"

	"github.com/Nuxij/goch8p/gfx"
	"github.com/Nuxij/goch8p/machine"
)

const (
	RAM_SIZE      = 0xFFF // 4096
	STACK_SIZE    = 0xF   // 16
	GFX_SIZE      = 0x800 // 2048
	KEYBOARD_SIZE = 0xF   // 16
)

type Runner struct {
	VM      *machine.Ch8p
	Display gfx.Display
}

func NewRunner(width, height uint16) *Runner {
	r := &Runner{
		VM: &machine.Ch8p{
			RAM: make(machine.Memory, RAM_SIZE),
			V:   make(machine.Registers, 8),
			Counters: machine.Counters{
				'T': 0,
				'O': 0,
				'I': 0,
				'S': 0,
				'P': 0x200,
			},
			Stack:    machine.Stack{0},
			GFX:      make(machine.Memory, width*height),
			Keyboard: make(machine.Memory, KEYBOARD_SIZE),
			Delay:    time.NewTicker(time.Second/60),
			Sound:    time.NewTicker(time.Second/60),
			LastOp:   "false",
		},
		Display: &gfx.ImScreen{},
	}
	r.VM.LoadFonts()
	err := r.Display.Init(int(1440), int(900))
	if err != nil {
		panic(err)
	}
	go r.VM.Pipeline()
	p := r.VM.ReadCounter('P')
	// Write 4,4 to V0 and V1
	r.VM.WriteRAMBytes(p, []byte{0x60, 0x04, 0x61, 0x04})
	// Write 4 to I
	r.VM.WriteRAMBytes(p+0x04, []byte{0xA0, 0x14})
	// Draw V0 and V1 5 height
	r.VM.WriteRAMBytes(p+0x06, []byte{0xD0, 0x15})
	
	r.VM.WriteCounter('I', 0)
	return r
}

func main() {
	runner := NewRunner(64, 32)
	hurts := time.NewTicker(time.Second/1)
	defer hurts.Stop()
	go func() {
		for {
			select {
			case _, ok := <-hurts.C:
				if runner.VM != nil {
					if !ok {
						fmt.Println("Hurts channel closed")
						panic("shiot")
					}
					runner.VM.Tick()
					runner.Display.Update(machine.Ch8pInfo{
						Tick:   runner.VM.ReadCounter('T'),
						Opcode: runner.VM.LastOp,
						RAM:   runner.VM.RAM,
						V:      runner.VM.V,
						I:      runner.VM.ReadCounter('I'),
						Stack:  runner.VM.Stack,
						
					}, runner.VM.GFX)
				}
			}
		}
	}()
	if runner.Display.Start() != nil {
		panic(fmt.Sprintf("runner failed to start %v", runner))
	}

}
