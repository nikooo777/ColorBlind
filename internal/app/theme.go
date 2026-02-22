package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type darkTheme struct{}

func (d *darkTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 30, G: 30, B: 34, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 50, G: 50, B: 56, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 40, G: 40, B: 44, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 36, G: 36, B: 40, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 40, G: 40, B: 46, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 200, G: 60, B: 120, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 220, G: 220, B: 225, A: 255}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 60, G: 60, B: 66, A: 255}
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (d *darkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (d *darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (d *darkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 13
	}
	return theme.DefaultTheme().Size(name)
}
