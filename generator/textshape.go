package generator

import (
	"image"
	"image/color"
	"math"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	targetWidthFraction  = 0.80
	targetHeightFraction = 0.25
	minFontSize          = 12.0
	maxFontSize          = 500.0
	binarySearchIters    = 20
	fontDPI              = 72
	letterSpacingFactor  = 0.15
)

type TextShape struct {
	mask    []bool
	maskW   int
	maskH   int
	offsetX int
	offsetY int
}

var (
	parsedFont     *opentype.Font
	parsedFontOnce sync.Once
	parsedFontErr  error
)

func parseFont(fontData []byte) (*opentype.Font, error) {
	parsedFontOnce.Do(func() {
		parsedFont, parsedFontErr = opentype.Parse(fontData)
	})
	return parsedFont, parsedFontErr
}

func NewTextShape(text string, fontData []byte, plateWidth, plateHeight int) *TextShape {
	f, err := parseFont(fontData)
	if err != nil {
		return newFallbackCircle(plateWidth, plateHeight)
	}

	targetW := float64(plateWidth) * targetWidthFraction
	targetH := float64(plateHeight) * targetHeightFraction

	fontSize := findFontSize(f, text, targetW, targetH)
	spacing := fixed.I(int(fontSize * letterSpacingFactor))

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  fontDPI,
	})
	if err != nil {
		return newFallbackCircle(plateWidth, plateHeight)
	}
	defer face.Close()

	textW, textH, baseline := measureWithSpacing(face, text, spacing)
	if textW <= 0 || textH <= 0 {
		return newFallbackCircle(plateWidth, plateHeight)
	}

	img := image.NewAlpha(image.Rect(0, 0, textW, textH))
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: face,
		Dot:  fixed.P(0, baseline),
	}
	drawStringWithSpacing(&drawer, text, spacing)

	mask := make([]bool, textW*textH)
	for y := range textH {
		for x := range textW {
			if img.AlphaAt(x, y).A > 0 {
				mask[y*textW+x] = true
			}
		}
	}

	return &TextShape{
		mask:    mask,
		maskW:   textW,
		maskH:   textH,
		offsetX: (plateWidth - textW) / 2,
		offsetY: (plateHeight - textH) / 2,
	}
}

func measureWithSpacing(face font.Face, text string, spacing fixed.Int26_6) (w, h, baseline int) {
	runes := []rune(text)
	var totalAdvance fixed.Int26_6
	var minY, maxY fixed.Int26_6

	for i, r := range runes {
		bounds, advance, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}
		totalAdvance += advance
		if i < len(runes)-1 {
			totalAdvance += spacing
		}
		if bounds.Min.Y < minY {
			minY = bounds.Min.Y
		}
		if bounds.Max.Y > maxY {
			maxY = bounds.Max.Y
		}
	}

	return totalAdvance.Ceil(), (maxY - minY).Ceil(), (-minY).Ceil()
}

func drawStringWithSpacing(d *font.Drawer, text string, spacing fixed.Int26_6) {
	runes := []rune(text)
	for i, r := range runes {
		d.DrawString(string(r))
		if i < len(runes)-1 {
			d.Dot.X += spacing
		}
	}
}

func findFontSize(f *opentype.Font, text string, targetW, targetH float64) float64 {
	lo := minFontSize
	hi := maxFontSize
	best := minFontSize

	for range binarySearchIters {
		mid := (lo + hi) / 2
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size: mid,
			DPI:  fontDPI,
		})
		if err != nil {
			hi = mid
			continue
		}

		spacing := fixed.I(int(mid * letterSpacingFactor))
		w, h, _ := measureWithSpacing(face, text, spacing)
		face.Close()

		if float64(w) <= targetW && float64(h) <= targetH {
			best = mid
			lo = mid
		} else {
			hi = mid
		}
	}

	return best
}

func newFallbackCircle(plateWidth, plateHeight int) *TextShape {
	cs := NewCircleShape(float64(plateWidth), float64(plateHeight))
	r := int(cs.circle.Radius)
	diameter := r * 2
	mask := make([]bool, diameter*diameter)
	rSq := float64(r * r)

	for dy := range diameter {
		for dx := range diameter {
			fx := float64(dx-r) + 0.5
			fy := float64(dy-r) + 0.5
			if fx*fx+fy*fy <= rSq {
				mask[dy*diameter+dx] = true
			}
		}
	}

	return &TextShape{
		mask:    mask,
		maskW:   diameter,
		maskH:   diameter,
		offsetX: int(cs.circle.X) - r,
		offsetY: int(cs.circle.Y) - r,
	}
}

func (ts *TextShape) Intersects(c Circle) bool {
	minX := int(math.Floor(c.X - c.Radius))
	maxX := int(math.Ceil(c.X + c.Radius))
	minY := int(math.Floor(c.Y - c.Radius))
	maxY := int(math.Ceil(c.Y + c.Radius))
	rSq := c.Radius * c.Radius

	for py := minY; py <= maxY; py++ {
		my := py - ts.offsetY
		if my < 0 || my >= ts.maskH {
			continue
		}
		for px := minX; px <= maxX; px++ {
			mx := px - ts.offsetX
			if mx < 0 || mx >= ts.maskW {
				continue
			}
			if !ts.mask[my*ts.maskW+mx] {
				continue
			}
			dx := float64(px) - c.X
			dy := float64(py) - c.Y
			if dx*dx+dy*dy <= rSq {
				return true
			}
		}
	}

	return false
}

func (ts *TextShape) DrawOutline(img *image.NRGBA) {
	bounds := img.Bounds()
	black := color.NRGBA{A: 255}

	for my := range ts.maskH {
		for mx := range ts.maskW {
			if !ts.mask[my*ts.maskW+mx] {
				continue
			}
			if !ts.isEdge(mx, my) {
				continue
			}
			px := mx + ts.offsetX
			py := my + ts.offsetY
			if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
				img.SetNRGBA(px, py, black)
			}
		}
	}
}

func (ts *TextShape) isEdge(mx, my int) bool {
	neighbors := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, n := range neighbors {
		nx, ny := mx+n[0], my+n[1]
		if nx < 0 || nx >= ts.maskW || ny < 0 || ny >= ts.maskH {
			return true
		}
		if !ts.mask[ny*ts.maskW+nx] {
			return true
		}
	}
	return false
}
