package generator

import "math"

func distance(a, b Circle) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func overlaps(a, b Circle) bool {
	return (a.Radius + b.Radius) > distance(a, b)
}

func intersects(a, b Circle) bool {
	minR := math.Min(a.Radius, b.Radius)
	maxR := math.Max(a.Radius, b.Radius)
	d := distance(a, b)

	if minR+maxR < d {
		return false
	}

	return d < (maxR+minR) && d > math.Abs(maxR-minR)
}
