package core

import (
	"sync"
)

type IQueue interface {
	Enqueue(t Task)
	DequeueAll() ([]Task, bool)
	Count() int
	Dispose()
}

type defaultQueue struct {
	paddingTasks []Task
	toDoTasks    []Task
	lock         *sync.Mutex
}

func (d *defaultQueue) init() *defaultQueue {
	d.toDoTasks = []Task{}
	d.paddingTasks = []Task{}
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

func (d *defaultQueue) Enqueue(task Task) {
	d.lock.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.lock.Unlock()
}

func (d *defaultQueue) DequeueAll() ([]Task, bool) {
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
