package textview

import (
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/muesli/reflow/wordwrap"
	"github.com/vwhitteron/go-pirateaudio/display"
)

type TextView struct {
	FontSize         int // font size
	LineSep          int // distance between lines in pixels
	vpos             uint8
	hpos             uint8
	dc               *gg.Context
	width            int
	margin           int
	dsp              *display.Display
	fontPath         string
	fgColor, bgColor Color
}

type Options struct {
	FontSize int
	FontPath string
	BGColor  Color
	FGColor  Color
}

var DefaultOpts = Options{
	FontSize: 16,
	BGColor:  Color{0, 0, 0},
	FGColor:  Color{255, 255, 255},
}

type Color [3]int

var YELLOW = Color{255, 255, 0}
var GREEN = Color{0, 255, 0}

func New() *TextView {
	return NewWithOptions(DefaultOpts)
}

func NewWithOptions(opts Options) *TextView {
	tv := &TextView{FontSize: opts.FontSize}
	tv.LineSep = 2
	tv.vpos = uint8(tv.FontSize)
	tv.margin = 2
	tv.hpos = uint8(tv.margin)
	tv.width = 240
	tv.dc = gg.NewContext(tv.width, tv.width)
	tv.loadFont(opts.FontPath)
	tv.clearContext()
	tv.bgColor = opts.BGColor
	tv.fgColor = opts.FGColor

	var err error
	tv.dsp, err = display.Init()
	if err != nil {
		panic(err)
	}

	return tv
}

func (t *TextView) loadFont(path string) {
	fonts := []string{
		path,
		"/usr/share/fonts/truetype/roboto/unhinted/RobotoTTF/Roboto-Medium.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	}

	floaded := false
	for _, f := range fonts {
		if f == "" {
			continue
		}
		if _, err := os.Stat(f); err == nil {
			if err := t.dc.LoadFontFace(f, float64(t.FontSize)); err != nil {
				panic(err)
			}
			floaded = true
			break
		}
	}
	if !floaded {
		panic("could not load a valid TTF font")
	}
}

func (t *TextView) drawText(text string) {
	text = wordwrap.String(text, t.width*2/t.FontSize)
	if strings.Contains(text, "\n") {
		s := strings.Split(text, "\n")
		for _, i := range s {
			t.drawText(i)
		}
	}

	if !strings.Contains(text, "\n") {
		if t.vpos >= uint8(t.width) {
			t.clearContext()
			t.vpos = uint8(t.FontSize)
		}
		t.dc.DrawString(text, float64(t.margin), float64(t.vpos))
		t.vpos += uint8(t.FontSize + t.LineSep)
	}
}

func (t *TextView) DrawChars(text string) {
	drawn := ""
	for _, s := range text {
		drawn += string(s)
		t.Draw(drawn)
		t.clearContext()
		time.Sleep(50 * time.Millisecond)
	}
}

func (t *TextView) clearContext() {
	t.dc.SetRGB255(t.bgColor[0], t.bgColor[1], t.bgColor[2])
	//t.dc.SetRGB(0, 0, 0)
	t.dc.Clear()
	t.dc.SetRGB255(t.fgColor[0], t.fgColor[1], t.fgColor[2])
}

func (t *TextView) Draw(text string) {
	t.drawText(text)
	t.drawToDisplay(text)
	t.resetPos()
}

func (t *TextView) resetPos() {
	t.vpos = uint8(t.FontSize)
	t.hpos = uint8(t.margin)
}

func (t *TextView) drawToDisplay(text string) {
	t.dsp.DrawRAW(t.dc.Image())
}

func (t *TextView) DrawFrames(textFrames []string) {
	for _, f := range textFrames {
		t.Draw(f)
		time.Sleep(50 * time.Millisecond)
	}
}
