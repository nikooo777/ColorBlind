package generator

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/niko/colorblind/palette"
)

const (
	maxAttemptsPer = 500
	tierMaxRadius  = 12.0
)

type Circle struct {
	X, Y   float64
	Radius float64
	Color  color.NRGBA
}

type radiusTier struct {
	minR, maxR float64
	fraction   float64
}

var tiers = []radiusTier{
	{9, 12, 0.20},
	{6, 9, 0.35},
	{3, 6, 0.45},
}

type spatialGrid struct {
	cellSize float64
	cells    map[[2]int][]int
}

func newSpatialGrid(cellSize float64) *spatialGrid {
	return &spatialGrid{
		cellSize: cellSize,
		cells:    make(map[[2]int][]int),
	}
}

func (g *spatialGrid) cellKey(x, y float64) [2]int {
	return [2]int{int(math.Floor(x / g.cellSize)), int(math.Floor(y / g.cellSize))}
}

func (g *spatialGrid) insert(idx int, c Circle) {
	key := g.cellKey(c.X, c.Y)
	g.cells[key] = append(g.cells[key], idx)
}

func (g *spatialGrid) hasOverlap(c Circle, circles []Circle) bool {
	key := g.cellKey(c.X, c.Y)
	spread := int(math.Ceil((tierMaxRadius*2)/g.cellSize)) + 1

	for dx := -spread; dx <= spread; dx++ {
		for dy := -spread; dy <= spread; dy++ {
			nk := [2]int{key[0] + dx, key[1] + dy}
			for _, idx := range g.cells[nk] {
				if overlaps(c, circles[idx]) {
					return true
				}
			}
		}
	}
	return false
}

func GenerateCircles(width, height float64, maxCount int) []Circle {
	circles := make([]Circle, 0, maxCount)
	grid := newSpatialGrid(tierMaxRadius * 2)

	remaining := maxCount
	for ti, tier := range tiers {
		var tierTarget int
		if ti == len(tiers)-1 {
			tierTarget = remaining
		} else {
			tierTarget = int(math.Round(float64(maxCount) * tier.fraction))
			if tierTarget > remaining {
				tierTarget = remaining
			}
		}
		remaining -= tierTarget

		failures := 0
		placed := 0
		for placed < tierTarget && failures < maxAttemptsPer {
			c := Circle{
				X:      width * rand.Float64(),
				Y:      height * rand.Float64(),
				Radius: tier.minR + rand.Float64()*(tier.maxR-tier.minR),
				Color:  palette.MagentaMain,
			}

			if grid.hasOverlap(c, circles) {
				failures++
				continue
			}

			grid.insert(len(circles), c)
			circles = append(circles, c)
			placed++
			failures = 0
		}
	}

	return circles
}

func HiddenShape(width, height float64) Circle {
	r := height / 3.0
	return Circle{
		X:      r + rand.Float64()*(width-2*r),
		Y:      height / 2.0,
		Radius: r,
	}
}

func RevealSecret(circles []Circle, shape Circle, revealColor color.NRGBA) {
	for i := range circles {
		if intersects(circles[i], shape) {
			circles[i].Color = revealColor
		}
	}
}
