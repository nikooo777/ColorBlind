package generator

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/niko/colorblind/palette"
)

const confuserFraction = 0.10

type Plate struct {
	Width, Height int
	circles       []Circle
	shape         Shape
	shapeFactory  func(w, h int) Shape
	intersecting  []int
	confusers     []int
}

func NewPlate(width, height, density int, shapeFactory func(w, h int) Shape) *Plate {
	p := &Plate{
		Width:        width,
		Height:       height,
		shapeFactory: shapeFactory,
	}
	p.generate(density)
	return p
}

func (p *Plate) SetShapeFactory(factory func(w, h int) Shape) {
	p.shapeFactory = factory
}

func (p *Plate) generate(density int) {
	p.circles = GenerateCircles(float64(p.Width), float64(p.Height), density)
	p.shape = p.shapeFactory(p.Width, p.Height)

	p.intersecting = nil
	p.confusers = nil
	intersectingSet := make(map[int]bool)
	for i, c := range p.circles {
		if p.shape.Intersects(c) {
			p.intersecting = append(p.intersecting, i)
			intersectingSet[i] = true
		}
	}

	var nonIntersecting []int
	for i := range p.circles {
		if !intersectingSet[i] {
			nonIntersecting = append(nonIntersecting, i)
		}
	}
	confuserCount := int(float64(len(nonIntersecting)) * confuserFraction)
	rand.Shuffle(len(nonIntersecting), func(i, j int) {
		nonIntersecting[i], nonIntersecting[j] = nonIntersecting[j], nonIntersecting[i]
	})
	p.confusers = nonIntersecting[:confuserCount]
}

func (p *Plate) Regenerate(density int) {
	p.generate(density)
}

func (p *Plate) CircleCount() int {
	return len(p.circles)
}

func (p *Plate) Render(revealColor color.NRGBA, confuserColor *color.NRGBA, showOutline bool) *image.NRGBA {
	for i := range p.circles {
		p.circles[i].Color = palette.MagentaMain
	}
	for _, i := range p.intersecting {
		p.circles[i].Color = revealColor
	}
	if confuserColor != nil {
		for _, i := range p.confusers {
			p.circles[i].Color = *confuserColor
		}
	}

	img := image.NewNRGBA(image.Rect(0, 0, p.Width, p.Height))
	fillBackground(img)

	for _, c := range p.circles {
		drawFilledCircle(img, c)
	}

	if showOutline {
		p.shape.DrawOutline(img)
	}

	return img
}

func drawCircleOutline(img *image.NRGBA, c Circle, col color.NRGBA, thickness float64) {
	cx := int(c.X)
	cy := int(c.Y)
	r := int(c.Radius) + int(thickness) + 1
	bounds := img.Bounds()

	outerSq := (c.Radius + thickness/2) * (c.Radius + thickness/2)
	innerSq := (c.Radius - thickness/2) * (c.Radius - thickness/2)

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			dSq := float64(dx*dx + dy*dy)
			if dSq <= outerSq && dSq >= innerSq {
				px, py := cx+dx, cy+dy
				if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
					img.SetNRGBA(px, py, col)
				}
			}
		}
	}
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
