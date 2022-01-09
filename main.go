package main

import (
	"fmt"
	"time"

	"github.com/AllenDang/giu"
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

func NewRunner(width, height int) *Runner {
	vm := &machine.Ch8p{
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
	}
	vm.LoadFonts()

	display := &gfx.ImScreen{
		Window: giu.NewMasterWindow("Joe", int(1440), int(900), giu.MasterWindowFlagsNotResizable),
		Shortcuts: []giu.WindowShortcut{
			{
				Key:      giu.KeyR,
				// Modifier: giu.ModControl,
				Callback: func() { vm.Running = !vm.Running },
			},
		},
	}
	err := display.Init(int(1440), int(900))
	if err != nil {
		panic(err)
	}
	
	vm.WriteCounter('I', 0)
	p := vm.ReadCounter('P')

	vm.WriteRAMBytes(p, []byte{0x60, 0x04, 0x61, 0x04})
	vm.WriteRAMBytes(p+0x04, []byte{0xA0, 0x14})
	vm.WriteRAMBytes(p+0x06, []byte{0xD0, 0x15})

	vm.WriteRAMBytes(p+0x08, []byte{0x60, 0x04, 0x61, 0x0A, 0xA0, 0x0A, 0xD0, 0x15})

	return &Runner{
		VM: vm,
		Display: display,
	}
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
					runner.VM.Cycle()
					runner.Display.Update(machine.Ch8pInfo{
						Tick:   runner.VM.ReadCounter('T'),
						Opcode: runner.VM.LastOp,
						RAM:    runner.VM.RAM,
						V:      runner.VM.V,
						I:      runner.VM.ReadCounter('I'),
						PC: 	runner.VM.ReadCounter('P'),
						Stack:  runner.VM.Stack,
						DrawFlag: runner.VM.DrawFlag,
						Running: runner.VM.Running,
					}, runner.VM.GFX)
				}
			}
		}
	}()
	if runner.Display.Start() != nil {
		panic(fmt.Sprintf("runner failed to start %v", runner))
	}

}
