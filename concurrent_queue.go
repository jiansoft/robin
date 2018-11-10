package robin

import (
	"container/list"
	"sync"
)

// ConcurrentQueue A "thread" safe string to anything container.
type ConcurrentQueue struct {
	lock      *sync.Mutex
	container *list.List
}

func (c *ConcurrentQueue) init() *ConcurrentQueue {
	c.container = list.New()
	c.lock = new(sync.Mutex)
	return c
}

// ConcurrentQueue Constructors
func NewConcurrentQueue() *ConcurrentQueue {
	return new(ConcurrentQueue).init()
}

// Enqueue Adds an object to the end of the ConcurrentQueue.
func (c *ConcurrentQueue) Enqueue(item interface{}) {
	c.lock.Lock()
	c.container.PushBack(item)
	c.lock.Unlock()
}

// TryPeek Tries to return an interface{} from the beginning of the ConcurrentQueue without removing it.
func (c *ConcurrentQueue) TryPeek() (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	lastItem := c.container.Back()
	if lastItem == nil {
		return nil, false
	}
	return lastItem.Value, true
}

// TryDequeue Tries to remove and return the interface{} at the beginning of the concurrent queue.
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

// Count Gets the number of elements contained in the ConcurrentQueue.
func (c ConcurrentQueue) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.container.Len()
}
