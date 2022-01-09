package gfx

import (
	"github.com/Nuxij/goch8p/machine"
)

type PixelMsg struct {
	Pixels []byte
	Info   machine.Ch8pInfo
	Callback  func()
}

type Display interface {
	Init(width, height int) error
	Start() error
	Update(info machine.Ch8pInfo, pixels []byte)
}