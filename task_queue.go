package robin

import (
	"sync"
)

type taskQueue interface {
	count() int
	dequeueAll() ([]Task, bool)
	dispose()
	enqueue(t Task)
}

// defaultQueue struct
type defaultQueue struct {
	paddingTasks []Task
	toDoTasks    []Task
	sync.Mutex
}

// newDefaultQueue return a new defaultQueue
func newDefaultQueue() *defaultQueue {
	q := &defaultQueue{toDoTasks: []Task{}, paddingTasks: []Task{}}
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
func (d *defaultQueue) enqueue(task Task) {
	d.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.Unlock()
}

// DequeueAll return current tasks
func (d *defaultQueue) dequeueAll() ([]Task, bool) {
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
