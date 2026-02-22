package app

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/niko/colorblind/internal/generator"
)

const (
	defaultWidth  = 960
	defaultHeight = 720
	plateWidth    = 700
	plateHeight   = 650

	defaultDensity = 2000
	minDensity     = 10
	maxDensity     = 2000

	defaultBlue = 170
	minBlue     = 156
	maxBlue     = 190

	sidebarWidth = 220
)

type ColorBlindApp struct {
	fyneApp       fyne.App
	window        fyne.Window
	plateImage    *canvas.Image
	plate         *generator.Plate
	density       int
	blueValue     uint8
	summaryLabel  *widget.Label
	showOutline   bool
	confuserColor *color.NRGBA
}

func New() *ColorBlindApp {
	a := app.NewWithID("com.colorblind.app")
	a.Settings().SetTheme(&darkTheme{})
	return &ColorBlindApp{
		fyneApp:   a,
		density:   defaultDensity,
		blueValue: defaultBlue,
	}
}

func (a *ColorBlindApp) Run() {
	a.window = a.fyneApp.NewWindow("ColorBlind")
	a.setupUI()
	a.window.Resize(fyne.NewSize(defaultWidth, defaultHeight))
	a.regenerateAndRender()
	a.window.ShowAndRun()
}

func (a *ColorBlindApp) setupUI() {
	a.plateImage = canvas.NewImageFromImage(nil)
	a.plateImage.FillMode = canvas.ImageFillContain
	a.plateImage.SetMinSize(fyne.NewSize(plateWidth, plateHeight))

	plateContainer := container.NewCenter(a.plateImage)

	sidebar := a.buildSidebar()

	sep := widget.NewSeparator()
	sidebarWithSep := container.NewBorder(nil, nil, sep, nil, sidebar)

	content := container.NewBorder(nil, nil, nil, sidebarWithSep, plateContainer)
	a.window.SetContent(content)

	menu := fyne.NewMenu("File",
		fyne.NewMenuItem("Quit", func() { a.fyneApp.Quit() }),
	)
	a.window.SetMainMenu(fyne.NewMainMenu(menu))
}

func (a *ColorBlindApp) buildSidebar() fyne.CanvasObject {
	generateBtn := widget.NewButtonWithIcon("Generate", theme.MediaPlayIcon(), func() {
		a.regenerateAndRender()
	})
	generateBtn.Importance = widget.HighImportance

	densityLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("Density: %d", a.density),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	densitySlider := widget.NewSlider(minDensity, maxDensity)
	densitySlider.Value = float64(a.density)
	densitySlider.Step = 10
	densitySlider.OnChanged = func(val float64) {
		a.density = int(val)
		densityLabel.SetText(fmt.Sprintf("Density: %d", a.density))
	}

	colorLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("Shade: %d", a.blueValue),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	colorSlider := widget.NewSlider(minBlue, maxBlue)
	colorSlider.Value = float64(a.blueValue)
	colorSlider.Step = 1
	colorSlider.OnChanged = func(val float64) {
		a.blueValue = uint8(val)
		colorLabel.SetText(fmt.Sprintf("Shade: %d", a.blueValue))
		a.render()
	}

	a.summaryLabel = widget.NewLabel("")
	a.summaryLabel.Alignment = fyne.TextAlignCenter

	densityCard := widget.NewCard("", "Circle Density", container.NewVBox(
		densityLabel,
		densitySlider,
	))

	colorCard := widget.NewCard("", "Reveal Color", container.NewVBox(
		colorLabel,
		colorSlider,
	))

	confuserPicker := newHSVPicker(func(c color.NRGBA) {
		a.confuserColor = &c
		a.render()
	})
	confuserCard := widget.NewCard("", "Confuser (10%)", confuserPicker.Container)

	outlineCheck := widget.NewCheck("Show hidden shape", func(checked bool) {
		a.showOutline = checked
		a.render()
	})

	sidebar := container.NewVBox(
		container.NewPadded(generateBtn),
		container.NewPadded(a.summaryLabel),
		densityCard,
		colorCard,
		confuserCard,
		container.NewPadded(outlineCheck),
	)

	sized := container.New(layout.NewGridWrapLayout(fyne.NewSize(sidebarWidth, 0)), sidebar)
	return sized
}

func (a *ColorBlindApp) revealColor() color.NRGBA {
	return color.NRGBA{R: 255, G: 0, B: a.blueValue, A: 255}
}

func (a *ColorBlindApp) regenerateAndRender() {
	if a.plate == nil {
		a.plate = generator.NewPlate(plateWidth, plateHeight, a.density)
	} else {
		a.plate.Regenerate(a.density)
	}
	a.summaryLabel.SetText(fmt.Sprintf("Placed: %d / %d", a.plate.CircleCount(), a.density))
	a.render()
}

func (a *ColorBlindApp) render() {
	if a.plate == nil {
		return
	}
	img := a.plate.Render(a.revealColor(), a.confuserColor, a.showOutline)
	a.plateImage.Image = img
	a.plateImage.Refresh()
}
