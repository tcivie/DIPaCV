package topology

import (
	"DIPaCV/pathfinding"
	"image"
)

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
			if !Are8Adjacent(point, (*path)[i+1]) {
				return false
			}
		}
	}

	return true
}

// AreConnected determines if two points p and q are connected within a region of interest S.
// Two elements p and q of a set S are considered connected if there exists a path between them
// that consists entirely of elements belonging to the set S. The path represents the sequence
// of adjacent pixels that connects p to q while staying within the region S.
//
// Parameters:
//   - p: Starting point to check for connectivity
//   - q: Ending point to check for connectivity
//   - S: Region of interest (set of points) in which to find the connecting path
//
// Returns:
//   - foundPath: If a path exists, returns the sequence of points forming the path from p to q
//   - hasPath: true if p and q are connected (a valid path exists), false otherwise
//
// Notes:
//   - Points p and q must both belong to region S for a path to exist
//   - The path, if it exists, consists entirely of points from region S
//   - The complement of S (points not in S) is considered the background region
func AreConnected(p image.Point, q image.Point, S Region) (foundPath []image.Point, hasPath bool) {
	if p == q && S[p] {
		return []image.Point{p}, true
	}
	if !S[p] || !S[q] || len(S) == 0 {
		return nil, false
	}

	foundPath, hasPath = pathfinding.AStar(p, q, func(p image.Point) bool {
		return S[p]
	})
	return
}
