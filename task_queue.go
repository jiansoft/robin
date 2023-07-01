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
	paddingTasks []task
	toDoTasks    []task
	sync.Mutex
}

// newDefaultQueue return a new defaultQueue
func newDefaultQueue() *defaultQueue {
	q := &defaultQueue{toDoTasks: []task{}, paddingTasks: []task{}}
	return q
}

// Dispose dispose defaultQueue
func (d *defaultQueue) dispose() {
	d.Lock()
	d.paddingTasks = nil
	d.toDoTasks = nil
	d.Unlock()
}

// Enqueue put a task into queue
func (d *defaultQueue) enqueue(task task) {
	d.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.Unlock()
}

// DequeueAll return current tasks
func (d *defaultQueue) dequeueAll() ([]task, bool) {
	d.Lock()
	defer d.Unlock()
	d.toDoTasks, d.paddingTasks = d.paddingTasks, d.toDoTasks
	d.paddingTasks = d.paddingTasks[:0]

	return d.toDoTasks, len(d.toDoTasks) > 0
}

// Count return padding tasks length
func (d *defaultQueue) count() int {
	d.Lock()
	count := len(d.paddingTasks)
	d.Unlock()

	return count
}
