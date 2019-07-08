package robin

import (
	"container/list"
	"sync"
)

/*
// ConcurrentCollection
type ConcurrentCollection interface {
    Push(element interface{})
    TryPeek() (interface{}, bool)
    TryPop() (interface{}, bool)
    ToArray() (elements []interface{})
    Len() int
    Clear()
}
*/
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

// Push adds an object to the end of the ConcurrentQueue.
func (c *ConcurrentQueue) Enqueue(element interface{}) {
	c.Lock()
	c.container.PushFront(element)
	c.Unlock()
}

// TryPeek tries to return an interface{} from the beginning of the ConcurrentQueue without removing it.
func (c *ConcurrentQueue) TryPeek() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	element := c.container.Back()
	if element == nil {
		return nil, false
	}

	return element.Value, true
}

// TryPop tries to remove and return the interface{} at the beginning of the ConcurrentQueue.
func (c *ConcurrentQueue) TryDequeue() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	element := c.container.Back()
	if element == nil {
		return nil, false
	}

	c.container.Remove(element)

	return element.Value, true
}

// Len gets the number of elements contained in the ConcurrentQueue.
func (c ConcurrentQueue) Len() int {
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

// ToArray copies the elements stored in the ConcurrentQueue to a new array.
func (c ConcurrentQueue) ToArray() (elements []interface{}) {
	c.Lock()
	defer c.Unlock()

	count := c.container.Len()
	elements = make([]interface{}, count)
	for temp, i := c.container.Back(), 0; temp != nil; temp, i = temp.Prev(), i+1 {
		elements[i] = temp.Value
	}

	return
}

// ConcurrentStack represents a thread-safe last in-first out (LIFO) collection.
type ConcurrentStack struct {
	sync.Mutex
	container *list.List
}

// NewConcurrentStack new a ConcurrentStack instance
func NewConcurrentStack() *ConcurrentStack {
	c := new(ConcurrentStack)
	c.container = list.New()
	return c
}

// Push adds an object to the end of the ConcurrentStack.
func (c *ConcurrentStack) Push(element interface{}) {
	c.Lock()
	c.container.PushFront(element)
	c.Unlock()
}

// TryPeek tries to return an interface{} from the beginning of the ConcurrentStack without removing it.
func (c *ConcurrentStack) TryPeek() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	element := c.container.Front()
	if element == nil {
		return nil, false
	}

	return element.Value, true
}

// TryPop attempts to pop and return the object at the top of the
func (c *ConcurrentStack) TryPop() (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	element := c.container.Front()
	if element == nil {
		return nil, false
	}

	c.container.Remove(element)

	return element.Value, true
}

// Len gets the number of elements contained in the ConcurrentStack.
func (c ConcurrentStack) Len() int {
	c.Lock()
	defer c.Unlock()
	return c.container.Len()
}

// Clear remove all element in the ConcurrentStack.
func (c *ConcurrentStack) Clear() {
	c.Lock()
	var next *list.Element
	for e := c.container.Front(); e != nil; e = next {
		next = e.Next()
		c.container.Remove(e)
	}
	c.Unlock()
}

//ToArray copies the elements stored in the ConcurrentStack to a new array.
func (c ConcurrentStack) ToArray() (elements []interface{}) {
	c.Lock()
	defer c.Unlock()

	count := c.container.Len()
	elements = make([]interface{}, count)
	for temp, i := c.container.Front(), 0; temp != nil; temp, i = temp.Next(), i+1 {
		elements[i] = temp.Value
	}

	return
}

// ConcurrentBag represents a thread-safe, unordered collection of interface{}.
type ConcurrentBag struct {
	sync.Mutex
	container []interface{}
}

// NewConcurrentStack new a ConcurrentStack instance
func NewConcurrentBag() *ConcurrentBag {
	c := new(ConcurrentBag)
	c.container = make([]interface{}, 0, 64)
	return c
}

// Add  an interface{} to the ConcurrentBag.
func (cb *ConcurrentBag) Add(element interface{}) {
	cb.Lock()
	s, c := len(cb.container), cap(cb.container)
	if s >= c {
		nc := make([]interface{}, s, c*2)
		copy(nc, cb.container)
		cb.container = nc
	}

	cb.container = append(cb.container, element)
	cb.Unlock()
}

// Len gets the number of elements contained in the ConcurrentBag.
func (cb ConcurrentBag) Len() int {
	cb.Lock()
	defer cb.Unlock()
	return len(cb.container)
}

// TryTake attempts to remove and return an interface{} from the ConcurrentBag
func (cb *ConcurrentBag) TryTake() (interface{}, bool) {
	cb.Lock()
	defer cb.Unlock()

	var s int
	if s = len(cb.container); s == 0 {
		return nil, false
	}

	c := cap(cb.container)
	if s < (c/2) && c > 64 {
		nc := make([]interface{}, s, c/2)
		copy(nc, cb.container)
		cb.container = nc
	}

	takeOne := cb.container[0]
	cb.container = cb.container[1:]
	return takeOne, true
}

//ToArray copies the ConcurrentBag elements to a new array.
func (cb *ConcurrentBag) ToArray() (elements []interface{}) {
	cb.Lock()
	defer cb.Unlock()
	nc := make([]interface{}, len(cb.container))
	copy(nc, cb.container)
	return nc
}

// Clear remove all element in the ConcurrentBag.
func (cb *ConcurrentBag) Clear() {
	cb.Lock()
	cb.container = cb.container[:0]
	cb.Unlock()
}
