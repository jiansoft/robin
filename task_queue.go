package robin

import (
	"sync"
)

type taskQueue interface {
	Enqueue(t Task)
	DequeueAll() ([]Task, bool)
	Count() int
	Dispose()
}

// DefaultQueue struct
type DefaultQueue struct {
	paddingTasks []Task
	toDoTasks    []Task
	lock         *sync.Mutex
}

func (d *DefaultQueue) init() *DefaultQueue {
	d.toDoTasks = []Task{}
	d.paddingTasks = []Task{}
	d.lock = new(sync.Mutex)
	return d
}

//NewDefaultQueue return a new DefaultQueue
func NewDefaultQueue() *DefaultQueue {
	return new(DefaultQueue).init()
}

// Dispose dispose DefaultQueue
func (d *DefaultQueue) Dispose() {
	d.paddingTasks = nil
	d.toDoTasks = nil
}

// Enqueue put a task into queue
func (d *DefaultQueue) Enqueue(task Task) {
	d.lock.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.lock.Unlock()
}

// DequeueAll return currrent tasks
func (d *DefaultQueue) DequeueAll() ([]Task, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.toDoTasks, d.paddingTasks = d.paddingTasks, d.toDoTasks
	d.paddingTasks = d.paddingTasks[:0]
	return d.toDoTasks, len(d.toDoTasks) > 0
}

// Count return padding tasks length
func (d *DefaultQueue) Count() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	return len(d.paddingTasks)
}
