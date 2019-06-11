package robin

import (
	"container/heap"
	"sync"
)

var (
	lock sync.Mutex
)

// Item store data in the PriorityQueue
type Item struct {
	Value    interface{}
	Priority int64
	Index    int
	sync.Mutex
}

// expired set the item has expired
func (item *Item) expired() {
	item.Lock()
	item.Priority = 0
	item.Unlock()
}

// A PriorityQueue implements heap.Interface and holds Items.
// ie. the 0th element is the lowest value
type PriorityQueue []*Item

func NewPriorityQueue(capacity int) *PriorityQueue {
	pg := make(PriorityQueue, 0, capacity)
	return &pg
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	pq[i].Lock()
	pq[j].Lock()
	defer func() {
		pq[j].Unlock()
		pq[i].Unlock()
	}()
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n, c := len(*pq), cap(*pq)
	if n >= c {
		npq := make(PriorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}

	*pq = (*pq)[0 : n+1]
	item := x.(*Item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) PushItem(item *Item) {
	lock.Lock()
	heap.Push(pq, item)
	lock.Unlock()
}

// update modifies the priority of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item) {
	lock.Lock()
	heap.Fix(pq, item.Index)
	lock.Unlock()
}

// TryDequeue
func (pq *PriorityQueue) TryDequeue(limit int64) (*Item, bool) {
	lock.Lock()
	defer lock.Unlock()
	if pq.Len() == 0 {
		return nil, false
	}

	item := (*pq)[0]
	item.Lock()
	defer item.Unlock()
	if item.Priority > limit {
		return nil, false
	}

	heap.Remove(pq, 0)

	return item, true
}
