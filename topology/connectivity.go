package topology

import "image"

// Region represents a set of pixels (region of interest)
type Region map[image.Point]bool

// IsDiscretePath determines if the given sequence of points forms a valid discrete path.
// A discrete path (or curve) from pixel p(x,y) to pixel q(s,t) is defined as a non-repeating
// sequence of adjacent pixels with coordinates (xi,yi), for 0 ≤ i ≤ n, where:
//   - (x0, y0) = (x, y) is the starting point
//   - (xn, yn) = (s, t) is the ending point
//   - Each consecutive pair of pixels must be adjacent
//   - No pixel appears more than once in the sequence
//
// Parameters:
//   - path: A pointer to a slice of image.Point representing the sequence of pixels
//
// Returns:
//   - true if the sequence forms a valid discrete path
//   - false if the path is nil, empty, contains repeated points, or has non-adjacent pixels
func IsDiscretePath(path *[]image.Point) bool {
	if path == nil || len(*path) == 0 {
		return false
	}

	if len(*path) == 1 {
		return true
	}

	visited := make(map[image.Point]bool)

	for i, point := range *path {
		if visited[point] {
			return false
		}
		visited[point] = true

		if i < len(*path)-1 {
			if !AreAdjacent(point, (*path)[i+1]) {
				return false
			}
		}
	}

	return true
}
