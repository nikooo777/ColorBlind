package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	toggleWidth  float32 = 60
	toggleHeight float32 = 28
	knobPadding  float32 = 3
)

var (
	toggleOnColor  = color.NRGBA{R: 200, G: 60, B: 120, A: 255}
	toggleOffColor = color.NRGBA{R: 80, G: 80, B: 86, A: 255}
	knobColor      = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

type toggleSwitch struct {
	widget.BaseWidget
	on        bool
	onChanged func(bool)
	bg        *canvas.Rectangle
	knob      *canvas.Circle
}

func newToggleSwitch(onChanged func(bool)) *toggleSwitch {
	t := &toggleSwitch{onChanged: onChanged}

	t.bg = canvas.NewRectangle(toggleOffColor)
	t.bg.CornerRadius = toggleHeight / 2

	t.knob = canvas.NewCircle(knobColor)

	t.ExtendBaseWidget(t)
	return t
}

func (t *toggleSwitch) CreateRenderer() fyne.WidgetRenderer {
	return &toggleRenderer{toggle: t}
}

func (t *toggleSwitch) Tapped(_ *fyne.PointEvent) {
	t.on = !t.on
	t.Refresh()
	if t.onChanged != nil {
		t.onChanged(t.on)
	}
}

func (t *toggleSwitch) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (t *toggleSwitch) MinSize() fyne.Size {
	return fyne.NewSize(toggleWidth, toggleHeight)
}

type toggleRenderer struct {
	toggle *toggleSwitch
}

func (r *toggleRenderer) Layout(size fyne.Size) {
	r.toggle.bg.Resize(size)
	r.toggle.bg.Move(fyne.NewPos(0, 0))

	knobSize := size.Height - knobPadding*2
	r.toggle.knob.Resize(fyne.NewSize(knobSize, knobSize))

	var knobX float32
	if r.toggle.on {
		knobX = size.Width - knobSize - knobPadding
	} else {
		knobX = knobPadding
	}
	r.toggle.knob.Move(fyne.NewPos(knobX, knobPadding))
}

func (r *toggleRenderer) MinSize() fyne.Size {
	return fyne.NewSize(toggleWidth, toggleHeight)
}

func (r *toggleRenderer) Refresh() {
	if r.toggle.on {
		r.toggle.bg.FillColor = toggleOnColor
	} else {
		r.toggle.bg.FillColor = toggleOffColor
	}
	r.toggle.bg.Refresh()
	r.toggle.knob.Refresh()
	r.Layout(r.toggle.Size())
}

func (r *toggleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.toggle.bg, r.toggle.knob}
}

func (r *toggleRenderer) Destroy() {}
