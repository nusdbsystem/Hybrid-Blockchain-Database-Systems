package veritas

import (
	"container/heap"
	"sync"
)

type PriorityQueue struct {
	mu    *sync.RWMutex
	nodes []*PqNode
}

type PqNode struct {
	Value    string
	Priority int
	index    int
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{mu: new(sync.RWMutex)}
	heap.Init(pq)
	return pq
}

func (pq *PriorityQueue) Put(v *PqNode) {
	defer pq.mu.Unlock()
	pq.mu.Lock()
	heap.Push(pq, v)
}

func (pq *PriorityQueue) Get() (interface{}, bool) {
	defer pq.mu.Unlock()
	pq.mu.Lock()
	if len(pq.nodes) > 0 {
		item := heap.Pop(pq)
		return item, true
	}
	return nil, false
}

func (pq PriorityQueue) Size() int {
	defer pq.mu.RUnlock()
	pq.mu.RLock()
	return len(pq.nodes)
}

func (pq *PriorityQueue) IsEmpty() bool {
	defer pq.mu.RUnlock()
	pq.mu.RLock()
	return !(len(pq.nodes) > 0)
}

func (pq PriorityQueue) Len() int {
	return len(pq.nodes)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq.nodes[i].Priority > pq.nodes[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.nodes[i], pq.nodes[j] = pq.nodes[j], pq.nodes[i]
	pq.nodes[i].index, pq.nodes[j].index = i, j
}

func (pq *PriorityQueue) Push(v interface{}) {
	item := v.(*PqNode)
	item.index = len(pq.nodes)
	pq.nodes = append(pq.nodes, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old.nodes)
	item := old.nodes[n-1]
	item.index = -1
	pq.nodes = old.nodes[0 : n-1]
	return item
}
