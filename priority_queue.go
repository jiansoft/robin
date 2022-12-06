package robin

import (
	"container/heap"
	"sync/atomic"
)

// Item store data in the PriorityQueue
type Item struct {
	Value    any
	Priority int64
	Index    int
}

// expired set the item has expired
func (item *Item) expired() {
	item.setPriority(0)
}

func (item *Item) setPriority(newVal int64) {
	atomic.StoreInt64(&item.Priority, newVal)
}

func (item *Item) getPriority() (priority int64) {
	priority = atomic.LoadInt64(&item.Priority)
	return
}

// A PriorityQueue implements heap.Interface and holds Items.
// ie. the 0th element is the lowest value
type PriorityQueue []*Item

func NewPriorityQueue(capacity int) *PriorityQueue {
	pg := make(PriorityQueue, 0, capacity)
	return &pg
}

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	var isLess = (*pq)[i].getPriority() < (*pq)[j].getPriority()
	return isLess
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
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

func (pq *PriorityQueue) PushItem(item *Item) {
	heap.Push(pq, item)
}

// Update modifies the priority of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item) {
	heap.Fix(pq, item.Index)
}

func (pq *PriorityQueue) Pop() any {
	n, c := len(*pq), cap(*pq)
	if n < (c/2) && c > 64 {
		npq := make(PriorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]

	return item
}

func (pq *PriorityQueue) PopItem(limit int64) (*Item, bool) {
	if pq.Len() == 0 {
		return nil, false
	}

	item := (*pq)[0]
	priority := item.getPriority()
	if priority > limit {
		return nil, false
	}

	heap.Remove(pq, 0)

	return item, true
}
