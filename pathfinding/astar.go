package pathfinding

import (
	"DIPaCV/topology"
	"container/heap"
	"image"
	"math"
)

// ValidPointChecker function type for checking if a point is valid/walkable
type ValidPointChecker func(image.Point) bool

// astarNode implements PriorityItem interface
type astarNode struct {
	point  image.Point
	gCost  float64
	fCost  float64
	parent *astarNode
	index  int
}

func (n *astarNode) Priority() float64 { return n.fCost }
func (n *astarNode) SetIndex(i int)    { n.index = i }
func (n *astarNode) GetIndex() int     { return n.index }

// AStar finds the shortest path using A* algorithm
func AStar(start, goal image.Point, isValidPoint ValidPointChecker) ([]image.Point, bool) {
	if !isValidPoint(start) || !isValidPoint(goal) {
		return nil, false
	}

	if start == goal {
		return []image.Point{start}, true
	}

	openSet := NewPriorityQueue[*astarNode]()
	closedSet := make(map[image.Point]*astarNode)

	startNode := &astarNode{
		point:  start,
		gCost:  0,
		fCost:  manhattanLength(start, goal),
		parent: nil,
	}

	heap.Push(openSet, startNode)

	// 8-connected neighbors
	directions := []image.Point{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0} /*************/, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*astarNode)

		if current.point == goal {
			return reconstructPath(current), true
		}

		closedSet[current.point] = current

		for _, dir := range directions {
			neighbor := image.Point{
				X: current.point.X + dir.X,
				Y: current.point.Y + dir.Y,
			}

			if !isValidPoint(neighbor) || closedSet[neighbor] != nil {
				continue
			}

			// Movement cost
			moveCost := 1.0
			if dir.X != 0 && dir.Y != 0 {
				moveCost = math.Sqrt2
			}

			tentativeGCost := current.gCost + moveCost

			existingNode := findNodeInQueue(openSet, neighbor)
			if existingNode == nil {
				neighborNode := &astarNode{
					point:  neighbor,
					gCost:  tentativeGCost,
					fCost:  tentativeGCost + manhattanLength(neighbor, goal),
					parent: current,
				}
				heap.Push(openSet, neighborNode)
			} else if tentativeGCost < existingNode.gCost {
				existingNode.gCost = tentativeGCost
				existingNode.fCost = tentativeGCost + manhattanLength(neighbor, goal)
				existingNode.parent = current
				heap.Fix(openSet, existingNode.GetIndex())
			}
		}
	}

	return nil, false
}

func manhattanLength(neighbor image.Point, goal image.Point) float64 {
	return float64(topology.ManhattanLength(neighbor, goal))
}

func reconstructPath(node *astarNode) []image.Point {
	var path []image.Point
	current := node

	for current != nil {
		path = append([]image.Point{current.point}, path...)
		current = current.parent
	}

	return path
}

func findNodeInQueue(pq *PriorityQueue[*astarNode], point image.Point) *astarNode {
	for _, node := range *pq {
		if node.point == point {
			return node
		}
	}
	return nil
}
