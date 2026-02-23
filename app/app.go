package app

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/niko/colorblind/generator"
)

const (
	defaultWidth  = 1100
	defaultHeight = 1150
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
	shapeMode     string
	shapeText     string
}

func New() *ColorBlindApp {
	a := app.NewWithID("com.colorblind.app")
	a.Settings().SetTheme(&darkTheme{})
	return &ColorBlindApp{
		fyneApp:   a,
		density:   defaultDensity,
		blueValue: defaultBlue,
		shapeMode: "Circle",
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

	sidebar := a.buildSidebar()

	sep := widget.NewSeparator()
	sidebarWithSep := container.NewBorder(nil, nil, sep, nil, sidebar)

	content := container.NewBorder(nil, nil, nil, sidebarWithSep, a.plateImage)
	a.window.SetContent(content)

	menu := fyne.NewMenu("File",
		fyne.NewMenuItem("Save PNG...", func() { a.savePNG() }),
		fyne.NewMenuItemSeparator(),
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

	textEntry := widget.NewEntry()
	textEntry.SetPlaceHolder("Enter text...")
	textEntry.Disable()
	textEntry.OnChanged = func(s string) {
		a.shapeText = s
	}

	shapeRadio := widget.NewRadioGroup([]string{"Circle", "Text"}, func(selected string) {
		a.shapeMode = selected
		if selected == "Text" {
			textEntry.Enable()
		} else {
			textEntry.Disable()
		}
	})
	shapeRadio.Horizontal = true
	shapeRadio.SetSelected("Circle")

	shapeCard := widget.NewCard("", "Hidden Shape", container.NewVBox(
		shapeRadio,
		textEntry,
	))

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

	outlineToggle := newToggleSwitch(func(on bool) {
		a.showOutline = on
		a.render()
	})
	outlineRow := container.NewHBox(widget.NewLabel("Show Shape"), outlineToggle)
	outlineCard := widget.NewCard("", "", outlineRow)

	sidebarContent := container.NewVBox(
		container.NewPadded(generateBtn),
		container.NewPadded(a.summaryLabel),
		shapeCard,
		densityCard,
		colorCard,
		confuserCard,
		outlineCard,
	)

	scrollable := container.NewVScroll(sidebarContent)
	scrollable.SetMinSize(fyne.NewSize(sidebarWidth, 0))
	return scrollable
}

func (a *ColorBlindApp) revealColor() color.NRGBA {
	return color.NRGBA{R: 255, G: 0, B: a.blueValue, A: 255}
}

func (a *ColorBlindApp) shapeFactory() func(w, h int) generator.Shape {
	if a.shapeMode == "Text" && a.shapeText != "" {
		text := a.shapeText
		fontData := theme.DefaultTextBoldFont().Content()
		return func(w, h int) generator.Shape {
			return generator.NewTextShape(text, fontData, w, h)
		}
	}
	return func(w, h int) generator.Shape {
		return generator.NewCircleShape(float64(w), float64(h))
	}
}

func (a *ColorBlindApp) regenerateAndRender() {
	factory := a.shapeFactory()
	if a.plate == nil {
		a.plate = generator.NewPlate(plateWidth, plateHeight, a.density, factory)
	} else {
		a.plate.SetShapeFactory(factory)
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

func (a *ColorBlindApp) savePNG() {
	if a.plateImage.Image == nil {
		return
	}
	d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()
		err = png.Encode(writer, a.plateImage.Image)
		if err != nil {
			dialog.ShowError(err, a.window)
		}
	}, a.window)
	d.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
	d.SetFileName("reverse_colorblind.png")
	cwd, err := os.Getwd()
	if err == nil {
		uri := storage.NewFileURI(cwd)
		listable, err := storage.ListerForURI(uri)
		if err == nil {
			d.SetLocation(listable)
		}
	}
	d.Show()
}
