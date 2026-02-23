package generator

import (
	"image"
	"image/color"
)

type Shape interface {
	Intersects(c Circle) bool
	DrawOutline(img *image.NRGBA)
}

type CircleShape struct {
	circle Circle
}

func NewCircleShape(width, height float64) *CircleShape {
	return &CircleShape{circle: HiddenShape(width, height)}
}

func (cs *CircleShape) Intersects(c Circle) bool {
	return intersects(c, cs.circle)
}

func (cs *CircleShape) DrawOutline(img *image.NRGBA) {
	drawCircleOutline(img, cs.circle, color.NRGBA{A: 255}, 2)
}
