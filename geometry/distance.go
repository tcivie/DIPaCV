package geometry

import (
	"image"
	"math"
)

// ManhattanLength calculates the Manhattan distance between two points.
// The Manhattan distance is the sum of the absolute differences of coordinates
// and represents the length of the path between two points when movement is
// restricted to horizontal and vertical directions only.
func ManhattanLength(p image.Point, q image.Point) int {
	dx := p.X - q.X
	dy := p.Y - q.Y

	// int ABS
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}

// EuclideanLength calculates the Euclidean distance between two points.
// The Euclidean distance represents the length of the straight line connecting
// two points in a plane and is calculated using the Pythagorean theorem.
func EuclideanLength(p image.Point, q image.Point) float64 {
	dx := float64(p.X - q.X)
	dy := float64(p.Y - q.Y)

	return math.Sqrt(dx*dx + dy*dy)
}
