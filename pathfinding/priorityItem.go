package pathfinding

import "container/heap"

// PriorityItem interface that items in the priority queue must implement
type PriorityItem interface {
	Priority() float64
	SetIndex(int)
	GetIndex() int
}

// PriorityQueue is a generic priority queue for pathfinding algorithms
type PriorityQueue[T PriorityItem] []T

func (pq *PriorityQueue[T]) Len() int { return len(*pq) }

func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return (*pq)[i].Priority() < (*pq)[j].Priority()
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].SetIndex(i)
	(*pq)[j].SetIndex(j)
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	item := x.(T)
	item.SetIndex(n)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.SetIndex(-1)
	*pq = old[0 : n-1]
	return item
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue[T PriorityItem]() *PriorityQueue[T] {
	pq := &PriorityQueue[T]{}
	heap.Init(pq)
	return pq
}
