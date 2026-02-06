package robin

import (
	"sync"
)

type taskQueue interface {
	count() int
	dequeueAll() ([]task, bool)
	dispose()
	enqueue(t task)
}

// defaultQueue struct
type defaultQueue struct {
	pending  []task
	upcoming []task
	mu       sync.Mutex
}

// newDefaultQueue return a new defaultQueue
func newDefaultQueue() *defaultQueue {
	return &defaultQueue{}
}

// dispose dispose defaultQueue
func (dq *defaultQueue) dispose() {
	dq.mu.Lock()
	dq.pending = nil
	dq.upcoming = nil
	dq.mu.Unlock()
}

// enqueue put a task into queue
func (dq *defaultQueue) enqueue(task task) {
	dq.mu.Lock()
	dq.pending = append(dq.pending, task)
	dq.mu.Unlock()
}

// dequeueAll return current tasks
func (dq *defaultQueue) dequeueAll() ([]task, bool) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	dq.upcoming, dq.pending = dq.pending, dq.upcoming
	dq.pending = dq.pending[:0]

	return dq.upcoming, len(dq.upcoming) > 0
}

// count return pending tasks length
func (dq *defaultQueue) count() int {
	dq.mu.Lock()
	count := len(dq.pending)
	dq.mu.Unlock()

	return count
}
