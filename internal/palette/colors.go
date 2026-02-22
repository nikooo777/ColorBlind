package palette

import (
	"image/color"
	"math/rand/v2"
)

const (
	minOpacity = 128
	maxOpacity = 255
)

var (
	MagentaMain      = color.NRGBA{R: 255, G: 0, B: 155, A: 255}
	MagentaSecondary = color.NRGBA{R: 255, G: 0, B: 170, A: 255}
	Background       = color.NRGBA{R: 218, G: 218, B: 218, A: 255}
)

func randomOpacity() uint8 {
	return uint8(rand.IntN(maxOpacity-minOpacity+1) + minOpacity)
}

func darkGreen() color.NRGBA {
	return color.NRGBA{R: 0, G: 65, B: 1, A: randomOpacity()}
}

func darkBrown() color.NRGBA {
	return color.NRGBA{R: 75, G: 56, B: 0, A: randomOpacity()}
}

func bluish() color.NRGBA {
	return color.NRGBA{R: 141, G: 161, B: 222, A: randomOpacity()}
}

func pinkish() color.NRGBA {
	return color.NRGBA{R: 169, G: 148, B: 224, A: randomOpacity()}
}

func green() color.NRGBA {
	return color.NRGBA{R: 0, G: 174, B: 0, A: randomOpacity()}
}

func darkRed() color.NRGBA {
	return color.NRGBA{R: 110, G: 0, B: 0, A: randomOpacity()}
}

func blue() color.NRGBA {
	return color.NRGBA{R: 31, G: 31, B: 248, A: randomOpacity()}
}

func red() color.NRGBA {
	return color.NRGBA{R: 255, G: 34, B: 51, A: randomOpacity()}
}

func grey() color.NRGBA {
	return color.NRGBA{R: 119, G: 123, B: 140, A: randomOpacity()}
}

func orange() color.NRGBA {
	return color.NRGBA{R: 242, G: 142, B: 244, A: randomOpacity()}
}

func RandomColor() color.NRGBA {
	r := rand.Float64()
	switch {
	case r < 0.1:
		return pinkish()
	case r < 0.2:
		return grey()
	case r < 0.3:
		return green()
	case r < 0.4:
		return darkRed()
	case r < 0.5:
		return darkGreen()
	case r < 0.6:
		return blue()
	case r < 0.7:
		return darkBrown()
	case r < 0.8:
		return red()
	case r < 0.9:
		return orange()
	default:
		return bluish()
	}
}
