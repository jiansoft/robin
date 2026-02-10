package robin

import (
	"sync"
)

// taskQueue defines the internal interface for fiber task queues.
type taskQueue interface {
	count() int
	dequeueAll() ([]task, bool)
	dispose()
	enqueue(t task)
}

// defaultQueue is a double-buffered task queue that swaps pending and upcoming slices on dequeue.
type defaultQueue struct {
	pending  []task
	upcoming []task
	mu       sync.Mutex
}

// newDefaultQueue creates a new defaultQueue with pre-allocated empty slices.
// Both slices are initialized to non-nil so that nil exclusively represents the disposed state.
func newDefaultQueue() *defaultQueue {
	return &defaultQueue{
		pending:  make([]task, 0),
		upcoming: make([]task, 0),
	}
}

// dispose releases the queue by setting both slices to nil.
func (dq *defaultQueue) dispose() {
	dq.mu.Lock()
	dq.pending = nil
	dq.upcoming = nil
	dq.mu.Unlock()
}

// enqueue appends a task to the pending queue. No-op if disposed.
func (dq *defaultQueue) enqueue(t task) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if dq.pending == nil { // already disposed
		return
	}
	dq.pending = append(dq.pending, t)
}

// dequeueAll swaps the pending and upcoming buffers and returns all pending tasks.
func (dq *defaultQueue) dequeueAll() ([]task, bool) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	// Swap: pending becomes upcoming (the batch to return),
	// old upcoming becomes the new pending buffer.
	dq.upcoming, dq.pending = dq.pending, dq.upcoming

	// Clear references in the old upcoming (now pending) to avoid holding
	// pointers (fn, invoke, invokeArg, funcCache, paramsCache) that prevent GC.
	for i := range dq.pending {
		dq.pending[i] = task{}
	}
	dq.pending = dq.pending[:0]

	return dq.upcoming, len(dq.upcoming) > 0
}

// count returns the number of pending tasks in the queue.
func (dq *defaultQueue) count() int {
	dq.mu.Lock()
	count := len(dq.pending)
	dq.mu.Unlock()

	return count
}
