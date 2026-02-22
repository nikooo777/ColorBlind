package app

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	fieldSize = 180
	barWidth  = 20
)

var confuserPresets = []struct {
	name  string
	color color.NRGBA
}{
	{"#61B9B9", color.NRGBA{R: 0x61, G: 0xB9, B: 0xB9, A: 255}},
	{"#569B9B", color.NRGBA{R: 0x56, G: 0x9B, B: 0x9B, A: 255}},
	{"#929292", color.NRGBA{R: 0x92, G: 0x92, B: 0x92, A: 255}},
}

type hsvPicker struct {
	h, s, v   float64
	field     *hsvField
	bar       *hsvBar
	swatch    *canvas.Rectangle
	hexLabel  *widget.Label
	onChanged func(color.NRGBA)
	Container fyne.CanvasObject
}

func newHSVPicker(onChanged func(color.NRGBA)) *hsvPicker {
	p := &hsvPicker{
		h: 300, s: 1, v: 1,
		onChanged: onChanged,
	}

	c := p.currentColor()
	p.swatch = canvas.NewRectangle(c)
	p.swatch.SetMinSize(fyne.NewSize(0, 24))
	p.swatch.CornerRadius = 4

	p.hexLabel = widget.NewLabel(colorToHex(c))
	p.hexLabel.Alignment = fyne.TextAlignCenter

	p.field = newHSVField(p)
	p.bar = newHSVBar(p)

	fieldBox := container.New(layout.NewGridWrapLayout(fyne.NewSize(fieldSize, fieldSize)), p.field)
	barBox := container.New(layout.NewGridWrapLayout(fyne.NewSize(barWidth, fieldSize)), p.bar)

	presetSwatches := make([]fyne.CanvasObject, 0, len(confuserPresets))
	for _, preset := range confuserPresets {
		c := preset.color
		sw := newTappableSwatch(c, func() {
			p.setColor(c)
		})
		presetSwatches = append(presetSwatches, sw)
	}
	presets := container.New(layout.NewGridLayout(len(confuserPresets)), presetSwatches...)

	p.Container = container.NewVBox(
		p.swatch,
		p.hexLabel,
		presets,
		container.NewHBox(fieldBox, barBox),
	)

	return p
}

func (p *hsvPicker) currentColor() color.NRGBA {
	return hsvToRGB(p.h, p.s, p.v)
}

func (p *hsvPicker) setColor(c color.NRGBA) {
	p.h, p.s, p.v = rgbToHSV(c)
	p.field.Refresh()
	p.bar.Refresh()
	p.notify()
}

func (p *hsvPicker) notify() {
	c := p.currentColor()
	p.swatch.FillColor = c
	p.swatch.Refresh()
	p.hexLabel.SetText(colorToHex(c))
	if p.onChanged != nil {
		p.onChanged(c)
	}
}

func colorToHex(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// --- tappable preset swatch ---

type tappableSwatch struct {
	widget.BaseWidget
	color   color.NRGBA
	onTap   func()
}

func newTappableSwatch(c color.NRGBA, onTap func()) *tappableSwatch {
	s := &tappableSwatch{color: c, onTap: onTap}
	s.ExtendBaseWidget(s)
	return s
}

func (s *tappableSwatch) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(s.color)
	rect.SetMinSize(fyne.NewSize(0, 24))
	rect.CornerRadius = 4
	return widget.NewSimpleRenderer(rect)
}

func (s *tappableSwatch) Tapped(_ *fyne.PointEvent) {
	if s.onTap != nil {
		s.onTap()
	}
}

// --- color field (hue x saturation square) ---

type hsvField struct {
	widget.BaseWidget
	picker *hsvPicker
}

func newHSVField(p *hsvPicker) *hsvField {
	f := &hsvField{picker: p}
	f.ExtendBaseWidget(f)
	return f
}

func (f *hsvField) CreateRenderer() fyne.WidgetRenderer {
	raster := canvas.NewRaster(f.draw)
	return widget.NewSimpleRenderer(raster)
}

func (f *hsvField) Tapped(ev *fyne.PointEvent) {
	f.pick(ev.Position)
}

func (f *hsvField) Dragged(ev *fyne.DragEvent) {
	f.pick(ev.Position)
}

func (f *hsvField) DragEnd() {}

func (f *hsvField) pick(pos fyne.Position) {
	size := f.Size()
	x := clampF(float64(pos.X/size.Width), 0, 1)
	y := clampF(float64(pos.Y/size.Height), 0, 1)
	f.picker.h = x * 360
	f.picker.s = 1.0 - y
	f.Refresh()
	f.picker.notify()
}

func (f *hsvField) draw(w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	if w < 2 || h < 2 {
		return img
	}

	for py := 0; py < h; py++ {
		s := 1.0 - float64(py)/float64(h-1)
		for px := 0; px < w; px++ {
			hue := float64(px) / float64(w-1) * 360
			img.SetNRGBA(px, py, hsvToRGB(hue, s, f.picker.v))
		}
	}

	cx := int(f.picker.h / 360 * float64(w-1))
	cy := int((1.0 - f.picker.s) * float64(h-1))
	drawCrosshair(img, cx, cy)

	return img
}

// --- value bar (vertical brightness strip) ---

type hsvBar struct {
	widget.BaseWidget
	picker *hsvPicker
}

func newHSVBar(p *hsvPicker) *hsvBar {
	b := &hsvBar{picker: p}
	b.ExtendBaseWidget(b)
	return b
}

func (b *hsvBar) CreateRenderer() fyne.WidgetRenderer {
	raster := canvas.NewRaster(b.draw)
	return widget.NewSimpleRenderer(raster)
}

func (b *hsvBar) Tapped(ev *fyne.PointEvent) {
	b.pick(ev.Position)
}

func (b *hsvBar) Dragged(ev *fyne.DragEvent) {
	b.pick(ev.Position)
}

func (b *hsvBar) DragEnd() {}

func (b *hsvBar) pick(pos fyne.Position) {
	y := clampF(float64(pos.Y/b.Size().Height), 0, 1)
	b.picker.v = 1.0 - y
	b.Refresh()
	b.picker.field.Refresh()
	b.picker.notify()
}

func (b *hsvBar) draw(w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	if w < 2 || h < 2 {
		return img
	}

	for py := 0; py < h; py++ {
		v := 1.0 - float64(py)/float64(h-1)
		gray := uint8(v * 255)
		c := color.NRGBA{R: gray, G: gray, B: gray, A: 255}
		for px := 0; px < w; px++ {
			img.SetNRGBA(px, py, c)
		}
	}

	cy := int((1.0 - b.picker.v) * float64(h-1))
	black := color.NRGBA{A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	for px := 0; px < w; px++ {
		if cy > 0 {
			img.SetNRGBA(px, cy-1, white)
		}
		img.SetNRGBA(px, cy, black)
		if cy < h-1 {
			img.SetNRGBA(px, cy+1, white)
		}
	}

	return img
}

// --- helpers ---

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func drawCrosshair(img *image.NRGBA, cx, cy int) {
	bounds := img.Bounds()
	black := color.NRGBA{A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	for d := -6; d <= 6; d++ {
		if d >= -2 && d <= 2 {
			continue
		}
		px := cx + d
		if px >= bounds.Min.X && px < bounds.Max.X {
			img.SetNRGBA(px, cy, black)
			if cy > 0 {
				img.SetNRGBA(px, cy-1, white)
			}
			if cy < bounds.Max.Y-1 {
				img.SetNRGBA(px, cy+1, white)
			}
		}
		py := cy + d
		if py >= bounds.Min.Y && py < bounds.Max.Y {
			img.SetNRGBA(cx, py, black)
			if cx > 0 {
				img.SetNRGBA(cx-1, py, white)
			}
			if cx < bounds.Max.X-1 {
				img.SetNRGBA(cx+1, py, white)
			}
		}
	}
}

func hsvToRGB(h, s, v float64) color.NRGBA {
	c := v * s
	hh := math.Mod(h/60.0, 6)
	x := c * (1 - math.Abs(math.Mod(hh, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case hh < 1:
		r, g, b = c, x, 0
	case hh < 2:
		r, g, b = x, c, 0
	case hh < 3:
		r, g, b = 0, c, x
	case hh < 4:
		r, g, b = 0, x, c
	case hh < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.NRGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

func rgbToHSV(c color.NRGBA) (h, s, v float64) {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))
	delta := maxC - minC

	v = maxC

	if maxC == 0 {
		return 0, 0, v
	}
	s = delta / maxC

	if delta == 0 {
		return 0, s, v
	}

	switch maxC {
	case r:
		h = 60 * math.Mod((g-b)/delta, 6)
	case g:
		h = 60 * ((b-r)/delta + 2)
	case b:
		h = 60 * ((r-g)/delta + 4)
	}
	if h < 0 {
		h += 360
	}

	return h, s, v
}
