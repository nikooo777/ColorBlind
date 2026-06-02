package generator

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/nikooo777/ColorBlind/palette"
)

const (
	maxAttemptsPer            = 500
	standardScaleStartDensity = 1500
	standardScaleEndDensity   = 7000
	standardMinScale          = 0.43
)

type CircleProfile int

const (
	CircleProfileStandard CircleProfile = iota
	CircleProfileFine
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

type circleProfileConfig struct {
	tiers     []radiusTier
	maxRadius float64
}

var standardCircleProfile = circleProfileConfig{
	tiers: []radiusTier{
		{9, 12, 0.20},
		{6, 9, 0.35},
		{3, 6, 0.45},
	},
	maxRadius: 12,
}

var fineCircleProfile = circleProfileConfig{
	tiers: []radiusTier{
		{4.5, 6.5, 0.15},
		{3, 4.5, 0.35},
		{1.8, 3, 0.50},
	},
	maxRadius: 6.5,
}

type spatialGrid struct {
	cellSize float64
	maxR     float64
	cells    map[[2]int][]int
}

func newSpatialGrid(maxR float64) *spatialGrid {
	return &spatialGrid{
		cellSize: maxR * 2,
		maxR:     maxR,
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
	spread := int(math.Ceil((c.Radius+g.maxR)/g.cellSize)) + 1

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

func generateCircles(width, height float64, maxCount int, profile CircleProfile) []Circle {
	config := circleProfileFor(profile, maxCount)
	circles := make([]Circle, 0, maxCount)
	grid := newSpatialGrid(config.maxRadius)

	remaining := maxCount
	for ti, tier := range config.tiers {
		var tierTarget int
		if ti == len(config.tiers)-1 {
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

func circleProfileFor(profile CircleProfile, density int) circleProfileConfig {
	if profile == CircleProfileFine {
		return fineCircleProfile
	}
	return scaleCircleProfile(standardCircleProfile, standardScaleForDensity(density))
}

func standardScaleForDensity(density int) float64 {
	if density <= standardScaleStartDensity {
		return 1
	}
	if density >= standardScaleEndDensity {
		return standardMinScale
	}

	progress := float64(density-standardScaleStartDensity) / float64(standardScaleEndDensity-standardScaleStartDensity)
	return 1 - progress*(1-standardMinScale)
}

func scaleCircleProfile(profile circleProfileConfig, scale float64) circleProfileConfig {
	scaledTiers := make([]radiusTier, len(profile.tiers))
	for i, tier := range profile.tiers {
		scaledTiers[i] = radiusTier{
			minR:     tier.minR * scale,
			maxR:     tier.maxR * scale,
			fraction: tier.fraction,
		}
	}
	return circleProfileConfig{
		tiers:     scaledTiers,
		maxRadius: profile.maxRadius * scale,
	}
}

func HiddenShape(width, height float64) Circle {
	r := height / 3.0
	return Circle{
		X:      r + rand.Float64()*(width-2*r),
		Y:      height / 2.0,
		Radius: r,
	}
}
