package gfx

import (
	"fmt"

	"github.com/Nuxij/goch8p/machine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

var (
	BorderCute = lipgloss.Border{
		Top:         "._.:*:",
		Bottom:      "._.:*:",
		Left:        "|*",
		Right:       "|*",
		TopLeft:     "*",
		TopRight:    "*",
		BottomLeft:  "*",
		BottomRight: "*",
	}
	StyleDefault = lipgloss.NewStyle().BorderStyle(BorderCute).BorderForeground(lipgloss.Color("63"))
	StyleMapData = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FAFAFA")).
					Background(lipgloss.Color("#7D56F4")).
					PaddingTop(2).
					PaddingLeft(4)
)

type Firmware struct {
	tea.Model
	width, height int
	ready 	   bool
	ram           machine.Memory
	info     chan machine.Ch8pInfo
	events   chan tea.Msg
}

type TeaScreen struct {
	width    , height int
	firmware *Firmware
	mug *tea.Program
}

func (t *TeaScreen) Init(width, height int) error {
	t.width = width
	t.height = height
	t.firmware = NewFirmware(width, height)
	t.mug = tea.NewProgram(t.firmware, tea.WithMouseCellMotion())
	return nil
}

func (t *TeaScreen) Start() error {
	return t.mug.Start()
}

func (t *TeaScreen) Update(info machine.Ch8pInfo, pixels []byte) {
	t.mug.Send(PixelMsg{
		Pixels: pixels,
		Info:   info,
	})
}

func NewFirmware(width, height int) *Firmware {
	fw := &Firmware{
		ready:	false,
		width:  width,
		height: height,
		ram:    make(machine.Memory, width*height),
		info:   make(chan machine.Ch8pInfo, 1),
	}
	for i := 0; i < fw.width*fw.height; i++ {
		fw.ram.WriteByte(uint16(i), 9)
	}
	return fw
}

func (fw *Firmware) Init() tea.Cmd {
	return nil
}

func (fw *Firmware) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

		case tea.KeyMsg:
			if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
				return fw, tea.Quit
			}
		// Is it a key press?
		case PixelMsg:
			fw.info <- msg.Info
			for i, b := range msg.Pixels {
				fw.ram.WriteByte(uint16(i), b)
			}
		case tea.MouseMsg:
			// left click
			if msg.Type == tea.MouseLeft {
				fw.ram.WriteByte(uint16(msg.X+msg.Y*fw.width), 0xFF)
			}
	}

	return fw, nil
}

func (fw *Firmware) View() string {
	var s string

	s = mapView(fw)
	s = lipgloss.JoinHorizontal(lipgloss.Top, s, statsView(fw))
	return indent.String("\n"+s+"\n\n", 2)
}

func mapView(fw *Firmware) string {
	var s string
	for y := 0; y < fw.height; y++ {
		for x := 0; x < fw.width; x++ {
			s += fmt.Sprintf("%v", fw.ram.ReadByte(uint16(x+y*fw.width)))
		}
		s += "\n"
	}
	return s
}

func statsView(fw *Firmware) string {
	s := fmt.Sprintf("%vx%v\n", fw.width, fw.height)
	select {
	case info, ok := <- fw.info:
		if !ok {
			return "INFO CHANNEL CLOSED"
		}
		s += fmt.Sprintf("Tick %v\n", info.Tick)
		s += fmt.Sprintf("Opc %v\n", info.Opcode)
	default:
		s += "NO INFO"
	}
	return StyleDefault.Render(s)
}



// 	mapStyle := lipgloss.NewStyle().
// 					Width(int(fw.width)).Height(int(fw.height))

// 	// header := strings.Repeat("H", int(fw.width)) + "\n"
// 	content := ""
// 	for i, b := range fw.ram.ReadBytes(0, fw.width*fw.height) {
// 		if b == 0 {
// 			content += " "
// 		} else {
// 			c := "."
// 			// content is a lipgloss style with foreground based on position
// 			content += fmt.Sprintf("%v", lipgloss.NewStyle().Foreground(lipgloss.Color(hex.EncodeToString([]byte{b}))).Render(c))
// 		}
// 		if uint16(i)%fw.width == fw.width-1 {
// 			content += "\n"
// 		}
// 	}
// 	mapString := mapStyle.Render(content)



// // Join on the top edge
// // str := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB)
// return mapString
	
// }

