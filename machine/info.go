package machine

type Ch8pInfo struct {
	Name    string `json:"name",default:"joeks ch8p"`
	Version string `json:"version",default:"0.0.1"`
	Tick    uint16 `json:"tick",default:"0"`
	Opcode  string `json:"opcode",default:"none"`
}
