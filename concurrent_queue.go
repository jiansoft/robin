package robin

import (
	"container/list"
	"sync"
)

// ConcurrentQueue a "thread" safe string to anything items.
type ConcurrentQueue struct {
	sync.Mutex
	container *list.List
}

func (c *ConcurrentQueue) init() *ConcurrentQueue {
	c.container = list.New()
	return c
}

// NewConcurrentQueue ConcurrentQueue Constructors
func NewConcurrentQueue() *ConcurrentQueue {
	return new(ConcurrentQueue).init()
}

// Enqueue Adds an object to the end of the ConcurrentQueue.
func (c *ConcurrentQueue) Enqueue(item interface{}) {
	c.Lock()
	c.container.PushBack(item)
	c.Unlock()
}

// TryPeek Tries to return an interface{} from the beginning of the ConcurrentQueue without removing it.
func (c *ConcurrentQueue) TryPeek() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	lastItem := c.container.Back()
	if lastItem == nil {
		return nil, false
	}
	return lastItem.Value, true
}

// TryDequeue Tries to remove and return the interface{} at the beginning of the concurrent queue.
func (c *ConcurrentQueue) TryDequeue() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	lastItem := c.container.Back()
	if lastItem == nil {
		return nil, false
	}
	item := c.container.Remove(lastItem)
	return item, true
}

// Count Gets the number of elements contained in the ConcurrentQueue.
func (c ConcurrentQueue) Count() int {
	c.Lock()
	defer c.Unlock()
	return c.container.Len()
}

// Clean remove all element in the ConcurrentQueue.
func (c ConcurrentQueue) Clean() {
	c.Lock()
	var next *list.Element
	for e := c.container.Front(); e != nil; e = next {
		next = e.Next()
		c.container.Remove(e)
	}
	c.Unlock()
}
