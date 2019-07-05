package robin

import (
	"container/list"
	"sync"
)

// ConcurrentQueue represents a thread-safe first in-first out (FIFO) collection.
type ConcurrentQueue struct {
	sync.Mutex
	container *list.List
}

// NewConcurrentQueue new a ConcurrentQueue instance
func NewConcurrentQueue() *ConcurrentQueue {
	c := new(ConcurrentQueue)
	c.container = list.New()
	return c
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
	lastItem := c.container.Front()
	if lastItem == nil {
		return nil, false
	}
	return lastItem.Value, true
}

// TryDequeue Tries to remove and return the interface{} at the beginning of the concurrent queue.
func (c *ConcurrentQueue) TryDequeue() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	lastItem := c.container.Front()
	if lastItem == nil {
		return nil, false
	}

	c.container.Remove(lastItem)
	return lastItem.Value, true
}

// Count Gets the number of elements contained in the ConcurrentQueue.
func (c ConcurrentQueue) Count() int {
	c.Lock()
	defer c.Unlock()
	return c.container.Len()
}

// Clear remove all element in the ConcurrentQueue.
func (c *ConcurrentQueue) Clear() {
	c.Lock()
	var next *list.Element
	for e := c.container.Front(); e != nil; e = next {
		next = e.Next()
		c.container.Remove(e)
	}
	c.Unlock()
}

//ToArray copies the elements stored in the ConcurrentQueue to a new array.
func (c *ConcurrentQueue) ToArray() (elements []interface{}) {
	c.Lock()
	defer c.Unlock()
	count := c.container.Len()
	elements = make([]interface{}, count)
	for temp, i := c.container.Front(), 0; temp != nil; temp, i = temp.Next(), i+1 {
		elements[i] = temp.Value
	}
	return
}
