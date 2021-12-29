package gfx

import (
	"encoding/hex"
	"fmt"

	"github.com/Nuxij/goch8p/machine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	StyleMapData = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
)

type Firmware struct {
	tea.Model
	width, height uint16
	ready 	   bool
	ram           machine.Memory
	Machine       machine.Ch8pInfo
}

func NewFirmware(width, height uint16) *Firmware {
	cuppa := &Firmware{
		ready:	false,
		width:  width,
		height: height,
		ram:    make(machine.Memory, width*height),
	}
	return cuppa
}

func (fw *Firmware) Init() tea.Cmd {
	for i := uint16(0); i < fw.width*fw.height; i++ {
		fw.ram.WriteByte(i, 0x0E)
	}
	return nil
}

func (fw *Firmware) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

		case tea.KeyMsg:
			if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
				return fw, tea.Quit
			}
		// Is it a key press?
		case PixelMsg:
			fw.Machine = msg.Info
			
			for i, b := range msg.Pixels {
				fw.ram.WriteByte(uint16(i), b)
			}
		case tea.MouseMsg:
			// left click
			if msg.Type == tea.MouseLeft {
				fw.ram.WriteByte(uint16(msg.X)+uint16(msg.Y)*fw.width, 0xFF)
			}
		case tea.WindowSizeMsg:
			fw.ready = true
	}

	return fw, tea.Batch(cmds...)
}

func (fw *Firmware) View() string {
	if !fw.ready {
		return "\n  Initializing..."
	}

	// header := strings.Repeat("H", int(fw.width)) + "\n"
	content := ""
	for i, b := range fw.ram.ReadBytes(0, fw.width*fw.height) {
		if b == 0 {
			content += " "
		} else {
			c := "X"
			// content is a lipgloss style with foreground based on position
			content += fmt.Sprintf("%v", lipgloss.NewStyle().Foreground(lipgloss.Color(hex.EncodeToString([]byte{b}))).Render(c))
		}
		if uint16(i+1) % fw.width == 0 {
			content += "\n"
		}
	}
	mapString := lipgloss.NewStyle().
		Width(int(fw.width)).Height(int(fw.height)).
		Render(content)



// Join on the top edge
// str := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB)
return mapString
	
}

