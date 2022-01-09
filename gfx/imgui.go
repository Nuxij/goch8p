package gfx

import (
	"image"
	"image/color"

	"github.com/AllenDang/giu"
	"github.com/Nuxij/goch8p/machine"
)
type ImScreen struct {
	Window *giu.MasterWindow
	Width  int
	Height int
	Title  string
	info  machine.Ch8pInfo
	buffer *image.RGBA
	pixels []byte
	texture *giu.Texture
	memoryWidget *giu.MemoryEditorWidget
}

 func (s *ImScreen) Init(width, height int) error {
	s.Title = "Goch8p::IMGUI"
	window := giu.NewMasterWindow(s.Title, width, height, giu.MasterWindowFlagsNotResizable)
	s.Window = window
	s.Width = width
	s.Height = height
	s.memoryWidget = giu.MemoryEditor()
	buffer, err := giu.LoadImage("gfx/gopher.png")
	if err != nil {
		return err
	}
	s.buffer = buffer
	// image.NewRGBA(image.Rect(0, 0, 64, 32))
	
	return nil
}

// ImScreen.Start should draw to the window with imgui
func (s *ImScreen) Start() error {
	s.Window.Run(func() {
		s.Draw()
	})
	
	return nil
}

func (s *ImScreen) Draw() {
	giu.SingleWindow().Layout(
		giu.SplitLayout(giu.DirectionHorizontal, float32(s.Width)/8,
			giu.Child().Layout(
				giu.Labelf("Goch8p::IMGUI %d", s.info.Tick),
				giu.Labelf("OP: %v", s.info.Opcode),
			),
			giu.SplitLayout(giu.DirectionVertical, float32(s.Height)/4,
				giu.Custom(func() {
					if len(s.info.RAM) > 0 {
						s.memoryWidget.Build()
					}
				}),
				giu.Image(s.texture).Size(64*16, 32*16),
			),
		),
	)
}

func (s *ImScreen) Close() {
	s.Window.Close()
}

func (s *ImScreen) Update(info machine.Ch8pInfo, pixels []byte) {
	s.info = info
	s.pixels = pixels
	scale := 4
	m := image.NewRGBA(image.Rect(0, 0, 64*scale, 32*scale))
	// fill m with pixels
	if len(pixels) > 0 {
		for i := 0; i < 64; i++ {
			for j := 0; j < 32; j++ {
				index := i+j*64
				p := s.pixels[index]
				// px := byte(i+j)
				if p == 1 {
					p = 0xFF
				}
				c := color.RGBA{byte(int(p)+i*j), p, p, 0xFF}
				x := i * scale
				y := j * scale
				for xx := 0; xx < scale; xx++ {
					for yy := 0; yy < scale; yy++ {
						m.Set(x+xx, y+yy, c)
					}
				}
			}
		}
	}

	
	if len(s.info.RAM) > 0 {
		s.memoryWidget.Contents(s.info.RAM[0x200:0x300])
	}
	giu.NewTextureFromRgba(m, func(texture *giu.Texture) {
		s.texture = texture
		giu.Update()
	})
	
}

