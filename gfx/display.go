package gfx

import (
	"github.com/Nuxij/goch8p/machine"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	WIDTH  = 64
	HEIGHT = 32
)

type PixelMsg struct {
	Pixels []byte
	Info   machine.Ch8pInfo
	Callback  func()
}

type Display struct {
	width    , height uint16
	mug      *Firmware
	firmware *tea.Program
}

func NewDisplay(width, height uint16) *Display {
	display := &Display{
		width:  width,
		height: height,
		mug:    NewFirmware(width, height),
	}
	display.firmware = tea.NewProgram(display.mug, tea.WithMouseCellMotion())
	return display
}

func (d Display) Start() error {
	return d.firmware.Start()
}

func (d Display) Update(info machine.Ch8pInfo, pixels []byte) {
	d.firmware.Send(PixelMsg{
		Pixels: pixels,
		Info:   info,
	})
}

//////////////////////////////////////////////////////////////////////////////////
