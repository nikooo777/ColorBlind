package generator

import (
	"image"
	"image/color"

	"github.com/niko/colorblind/internal/palette"
)

type Plate struct {
	Width, Height int
	circles       []Circle
	shape         Circle
	intersecting  []int
}

func NewPlate(width, height, density int) *Plate {
	p := &Plate{Width: width, Height: height}
	p.generate(density)
	return p
}

func (p *Plate) generate(density int) {
	p.circles = GenerateCircles(float64(p.Width), float64(p.Height), density)
	p.shape = HiddenShape(float64(p.Width), float64(p.Height))

	p.intersecting = nil
	for i, c := range p.circles {
		if intersects(c, p.shape) {
			p.intersecting = append(p.intersecting, i)
		}
	}
}

func (p *Plate) Regenerate(density int) {
	p.generate(density)
}

func (p *Plate) CircleCount() int {
	return len(p.circles)
}

func (p *Plate) Render(revealColor color.NRGBA) *image.NRGBA {
	for _, i := range p.intersecting {
		p.circles[i].Color = revealColor
	}
	for i, c := range p.circles {
		if c.Color != revealColor && c.Color != palette.MagentaMain {
			p.circles[i].Color = palette.MagentaMain
		}
	}

	img := image.NewNRGBA(image.Rect(0, 0, p.Width, p.Height))
	fillBackground(img)

	for _, c := range p.circles {
		drawFilledCircle(img, c)
	}

	return img
}

func fillBackground(img *image.NRGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetNRGBA(x, y, palette.Background)
		}
	}
}

func drawFilledCircle(img *image.NRGBA, c Circle) {
	cx := int(c.X)
	cy := int(c.Y)
	r := int(c.Radius) + 1

	rSq := c.Radius * c.Radius
	bounds := img.Bounds()

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if float64(dx*dx+dy*dy) <= rSq {
				px, py := cx+dx, cy+dy
				if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
					blendPixel(img, px, py, c.Color)
				}
			}
		}
	}
}

func blendPixel(img *image.NRGBA, x, y int, src color.NRGBA) {
	if src.A == 255 {
		img.SetNRGBA(x, y, src)
		return
	}

	dst := img.NRGBAAt(x, y)
	srcA := float64(src.A) / 255.0
	invA := 1.0 - srcA

	img.SetNRGBA(x, y, color.NRGBA{
		R: uint8(float64(src.R)*srcA + float64(dst.R)*invA),
		G: uint8(float64(src.G)*srcA + float64(dst.G)*invA),
		B: uint8(float64(src.B)*srcA + float64(dst.B)*invA),
		A: 255,
	})
}
