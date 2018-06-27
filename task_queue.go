package robin

import (
	"sync"
)

type taskQueue interface {
	Enqueue(t task)
	DequeueAll() ([]task, bool)
	Count() int
	Dispose()
}

type defaultQueue struct {
	paddingTasks []task
	toDoTasks    []task
	lock         *sync.Mutex
}

func (d *defaultQueue) init() *defaultQueue {
	d.toDoTasks = []task{}
	d.paddingTasks = []task{}
	d.lock = new(sync.Mutex)
	return d
}

func NewDefaultQueue() *defaultQueue {
	return new(defaultQueue).init()
}

func (d *defaultQueue) Dispose() {
	d.paddingTasks = nil
	d.toDoTasks = nil
}

func (d *defaultQueue) Enqueue(task task) {
	d.lock.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.lock.Unlock()
}

func (d *defaultQueue) DequeueAll() ([]task, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.toDoTasks, d.paddingTasks = d.paddingTasks, d.toDoTasks
	d.paddingTasks = d.paddingTasks[:0]
	return d.toDoTasks, len(d.toDoTasks) > 0
}

func (d *defaultQueue) Count() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	return len(d.paddingTasks)
}
