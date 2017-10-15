package collections

import (
	"container/list"
	"sync"
)

type ConcurrentQueue struct {
	lock      *sync.Mutex
	container *list.List
}

func (c *ConcurrentQueue) init() *ConcurrentQueue {
	c.container = list.New()
	c.lock = new(sync.Mutex)
	return c
}

func NewConcurrentQueue() *ConcurrentQueue {
	return new(ConcurrentQueue).init()
}

func (c *ConcurrentQueue) Enqueue(item interface{}) {
	c.lock.Lock()
	c.container.PushBack(item)
	c.lock.Unlock()
}

func (c *ConcurrentQueue) TryPeek() (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	lastItem := c.container.Back()
	if lastItem == nil {
		return nil, false
	}
	return lastItem.Value, true
}

func (c *ConcurrentQueue) TryDequeue() (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	lastItem := c.container.Back()
	if lastItem == nil {
		return nil, false
	}
	item := c.container.Remove(lastItem)
	return item, true
}

func (c ConcurrentQueue) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.container.Len()
}
