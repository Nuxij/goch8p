package machine

type Ch8pInfo struct {
	Name    string `json:"name",default:"joeks ch8p"`
	Version string `json:"version",default:"0.0.1"`
	Tick    uint16 `json:"tick",default:"0"`
	Opcode  string `json:"opcode",default:"none"`
	RAM     []byte `json:"ram",default:"[]"`
}

// 5-high sprite for fonts
type Font [5]byte

// Fonts is a list of the default fonts 0-F
var Fonts = [16]Font{
	Font{0xF0, 0x90, 0x90, 0x90, 0xF0}, // 0
	Font{0x20, 0x60, 0x20, 0x20, 0x70}, // 1
	Font{0xF0, 0x10, 0xF0, 0x80, 0xF0}, // 2
	Font{0xF0, 0x10, 0xF0, 0x10, 0xF0}, // 3
	Font{0x90, 0x90, 0xF0, 0x10, 0x10}, // 4
	Font{0xF0, 0x80, 0xF0, 0x10, 0xF0}, // 5
	Font{0xF0, 0x80, 0xF0, 0x90, 0xF0}, // 6
	Font{0xF0, 0x10, 0x20, 0x40, 0x40}, // 7
	Font{0xF0, 0x90, 0xF0, 0x90, 0xF0}, // 8
	Font{0xF0, 0x90, 0xF0, 0x10, 0xF0}, // 9
	Font{0xF0, 0x90, 0xF0, 0x90, 0x90}, // A
	Font{0xE0, 0x90, 0xE0, 0x90, 0xE0}, // B
	Font{0xF0, 0x80, 0x80, 0x80, 0xF0}, // C
	Font{0xE0, 0x90, 0x90, 0x90, 0xE0}, // D
	Font{0xF0, 0x80, 0xF0, 0x80, 0xF0}, // E
	Font{0xF0, 0x80, 0xF0, 0x80, 0x80}, // F
}