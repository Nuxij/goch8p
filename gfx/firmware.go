package gfx

import (
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
	StyleMapData = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#000"))
)

type Firmware struct {
	tea.Model
	width, height int
	ram           machine.Memory
	Machine       machine.Ch8pInfo
}

func NewFirmware(width, height int) *Firmware {
	cuppa := &Firmware{
		width:  width,
		height: height,
		ram:    make(machine.Memory, width*height),
	}
	return cuppa
}

func (fw *Firmware) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (fw *Firmware) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case PixelMsg:
		fw.Machine = msg.Info
		for i, b := range msg.Pixels {
			fw.ram.WriteByte(uint16(i), b)
		}
		return fw, nil
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return fw, tea.Quit

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return fw, nil
}

func (fw *Firmware) View() string {
	// The header
	header := "JoeDraw\n"

	// Iterate over our choices
	mapData := header
	for i, bite := range fw.ram {
		// Render the row
		if bite != 0 {
			// num := fmt.Sprintf("%d", rand.Intn(i%fw.width+1))
			mapData += StyleMapData.Render(".")
			if i%fw.width == fw.width-1 {
				mapData += "\n"
			}
		}
	}
	mapData = StyleMapData.Render(mapData)

	// The footer
	footer := "(Press q to quit)\n"

	vmInfo := StyleMapData.Copy()
	info := "VM Info - " + lipgloss.PlaceHorizontal(fw.width, lipgloss.Center, footer)
	info += fmt.Sprintf("Opcode: %v\n", fw.Machine.Opcode)
	info += fmt.Sprintf("Tick: %d\n", fw.Machine.Tick)

	frames := lipgloss.JoinHorizontal(lipgloss.Top, StyleDefault.Render(mapData), vmInfo.Render(info))
	return lipgloss.JoinVertical(lipgloss.Top, frames, StyleDefault.Render(footer))
}
