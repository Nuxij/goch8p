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
	Display *gfx.Display
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
			Delay:    time.NewTicker(time.Second / 60),
			Sound:    time.NewTicker(time.Second / 60),
			LastOp:   "false",
		},
		Display: gfx.NewDisplay(width, height),
	}
	// initial program is a single draw call, 5 high sprite at 0,0
	for i := uint16(0); i < 5; i++ {
		r.VM.Write(r.VM.ReadCounter('P')*i, 0xD0)
		r.VM.Write(r.VM.ReadCounter('P')*i+1, 0x05)
	}
	r.VM.WriteCounter('I', 0)
	return r
}

func main() {
	runner := NewRunner(64, 64)
	hurts := time.NewTicker(time.Millisecond * 100)
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
					}, runner.VM.GFX)
				}
			}
		}
	}()
	if runner.Display.Start() != nil {
		panic(fmt.Sprintf("runner failed to start %v", runner))
	}

}
