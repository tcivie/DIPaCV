package topology

import "image"

// Are4Adjacent checks 4-connectivity (horizontal/vertical only)
func Are4Adjacent(p, q image.Point) bool {
	dx := abs(p.X - q.X)
	dy := abs(p.Y - q.Y)
	return (dx == 1 && dy == 0) || (dx == 0 && dy == 1)
}

// Are8Adjacent checks 8-connectivity (includes diagonals)
func Are8Adjacent(p, q image.Point) bool {
	dx := abs(p.X - q.X)
	dy := abs(p.Y - q.Y)
	return dx <= 1 && dy <= 1 && (dx != 0 || dy != 0)
}

// Are24Adjacent checks 24-connectivity (2nd level neighbors)
func Are24Adjacent(p, q image.Point) bool {
	return ChebyshevDistance(p, q) <= 2 && !p.Eq(q)
}

// ChebyshevDistance returns the Chebyshev distance (L∞ norm)
func ChebyshevDistance(p, q image.Point) int {
	return maxInt(abs(p.X-q.X), abs(p.Y-q.Y))
}

// ManhattanDistance returns the Manhattan distance (L1 norm)
func ManhattanDistance(p, q image.Point) int {
	return abs(p.X-q.X) + abs(p.Y-q.Y)
}

// GetNeighbors4 returns all 4-connected neighbors of p
func GetNeighbors4(p image.Point) []image.Point {
	return []image.Point{
		{p.X, p.Y - 1},
		{p.X - 1, p.Y}, {p.X + 1, p.Y},
		{p.X, p.Y + 1},
	}
}

// GetNeighbors8 returns all 8-connected neighbors of p
func GetNeighbors8(p image.Point) []image.Point {
	return []image.Point{
		{p.X - 1, p.Y - 1}, {p.X, p.Y - 1}, {p.X + 1, p.Y - 1},
		{p.X - 1, p.Y}, {p.X + 1, p.Y},
		{p.X - 1, p.Y + 1}, {p.X, p.Y + 1}, {p.X + 1, p.Y + 1},
	}
}

// GetNeighbors24 returns all 24-connected neighbors of p (2nd level)
func GetNeighbors24(p image.Point) []image.Point {
	neighbors := make([]image.Point, 0, 24)
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbors = append(neighbors, image.Point{X: p.X + dx, Y: p.Y + dy})
		}
	}
	return neighbors
}

// GetKthLevelNeighbors returns all neighbors at exactly the k-th level (Chebyshev distance = k)
// For k=1, this gives 8 neighbors. For k=2, this gives 16 neighbors (not including k=1).
func GetKthLevelNeighbors(p image.Point, k int) []image.Point {
	if k <= 0 {
		return nil
	}

	neighbors := make([]image.Point, 0, 8*k) // approximate capacity

	for dy := -k; dy <= k; dy++ {
		for dx := -k; dx <= k; dx++ {
			if maxInt(abs(dx), abs(dy)) == k {
				neighbors = append(neighbors, image.Point{X: p.X + dx, Y: p.Y + dy})
			}
		}
	}
	return neighbors
}

// GetNeighborsWithinK returns all neighbors within Chebyshev distance k (inclusive)
func GetNeighborsWithinK(p image.Point, k int) []image.Point {
	if k <= 0 {
		return nil
	}

	size := (2*k+1)*(2*k+1) - 1 // total cells minus center
	neighbors := make([]image.Point, 0, size)

	for dy := -k; dy <= k; dy++ {
		for dx := -k; dx <= k; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbors = append(neighbors, image.Point{X: p.X + dx, Y: p.Y + dy})
		}
	}
	return neighbors
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
