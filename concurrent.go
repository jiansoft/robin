package robin

import (
	"sync"
)

const defaultRingCapacity = 16

// ConcurrentQueue represents a thread-safe first in-first out (FIFO) collection.
// Uses a ring buffer with power-of-two capacity and bitwise masking for O(1) index calculations.
// Automatically grows when full and shrinks when usage drops below 25% capacity.
type ConcurrentQueue[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
	mask  int // len(buf) - 1, bitwise AND replaces modulo
	mu    sync.Mutex
}

// NewConcurrentQueue creates a new ConcurrentQueue instance.
func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		buf:  make([]T, defaultRingCapacity),
		mask: defaultRingCapacity - 1,
	}
}

// Enqueue adds an element to the end of the ConcurrentQueue.
func (q *ConcurrentQueue[T]) Enqueue(element T) {
	q.mu.Lock()
	if q.count == q.mask+1 {
		q.grow()
	}
	q.buf[q.tail] = element
	q.tail = (q.tail + 1) & q.mask
	q.count++
	q.mu.Unlock()
}

func (q *ConcurrentQueue[T]) grow() {
	q.resize((q.mask + 1) << 1)
}

func (q *ConcurrentQueue[T]) shrink() {
	curCap := q.mask + 1
	if curCap <= defaultRingCapacity || q.count >= curCap>>2 {
		return
	}
	// curCap is always a power of two > defaultRingCapacity,
	// so curCap>>1 >= defaultRingCapacity is guaranteed.
	q.resize(curCap >> 1)
}

func (q *ConcurrentQueue[T]) resize(newCap int) {
	newBuf := make([]T, newCap)
	if q.head < q.tail {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		copied := copy(newBuf, q.buf[q.head:])
		copy(newBuf[copied:], q.buf[:q.tail])
	}
	q.head = 0
	q.tail = q.count
	q.buf = newBuf
	q.mask = newCap - 1
}

// TryPeek tries to return an element from the beginning of the ConcurrentQueue without removing it.
func (q *ConcurrentQueue[T]) TryPeek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}

	return q.buf[q.head], true
}

// TryDequeue tries to remove and return the element at the beginning of the ConcurrentQueue.
func (q *ConcurrentQueue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}

	element := q.buf[q.head]
	var zero T
	q.buf[q.head] = zero // clear reference for GC
	q.head = (q.head + 1) & q.mask
	q.count--
	q.shrink()
	return element, true
}

// Len gets the number of elements contained in the ConcurrentQueue.
func (q *ConcurrentQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.count
}

// Clear removes all elements from the ConcurrentQueue and resets the buffer to default capacity.
func (q *ConcurrentQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.buf = make([]T, defaultRingCapacity)
	q.head = 0
	q.tail = 0
	q.count = 0
	q.mask = defaultRingCapacity - 1
}

// ToArray copies the elements stored in the ConcurrentQueue to a new slice.
func (q *ConcurrentQueue[T]) ToArray() []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	result := make([]T, q.count)
	for i := range q.count {
		result[i] = q.buf[(q.head+i)&q.mask]
	}

	return result
}

// ConcurrentStack represents a thread-safe last in-first out (LIFO) collection.
type ConcurrentStack[T any] struct {
	container []T
	mu        sync.Mutex
}

// NewConcurrentStack creates a new ConcurrentStack instance.
func NewConcurrentStack[T any]() *ConcurrentStack[T] {
	return &ConcurrentStack[T]{}
}

// Push adds an element to the top of the ConcurrentStack.
func (s *ConcurrentStack[T]) Push(element T) {
	s.mu.Lock()
	s.container = append(s.container, element)
	s.mu.Unlock()
}

// TryPeek tries to return the element at the top of the ConcurrentStack without removing it.
func (s *ConcurrentStack[T]) TryPeek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.container) == 0 {
		var zero T
		return zero, false
	}

	return s.container[len(s.container)-1], true
}

// TryPop attempts to pop and return the element at the top of the ConcurrentStack.
func (s *ConcurrentStack[T]) TryPop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return popLast(&s.container)
}

// Len gets the number of elements contained in the ConcurrentStack.
func (s *ConcurrentStack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.container)
}

// Clear removes all elements from the ConcurrentStack.
func (s *ConcurrentStack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.container)
	s.container = s.container[:0]
}

// ToArray copies the elements stored in the ConcurrentStack to a new slice.
// Elements are returned in LIFO order (top of stack first).
func (s *ConcurrentStack[T]) ToArray() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(s.container)
	result := make([]T, n)
	for i, j := n-1, 0; i >= 0; i, j = i-1, j+1 {
		result[j] = s.container[i]
	}

	return result
}

// ConcurrentBag represents a thread-safe, unordered collection of elements.
type ConcurrentBag[T any] struct {
	container []T
	mu        sync.Mutex
}

// NewConcurrentBag creates a new ConcurrentBag instance.
func NewConcurrentBag[T any]() *ConcurrentBag[T] {
	return &ConcurrentBag[T]{
		container: make([]T, 0, 64),
	}
}

// Add adds an element to the ConcurrentBag.
func (cb *ConcurrentBag[T]) Add(element T) {
	cb.mu.Lock()
	cb.container = append(cb.container, element)
	cb.mu.Unlock()
}

// Len gets the number of elements contained in the ConcurrentBag.
func (cb *ConcurrentBag[T]) Len() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return len(cb.container)
}

// TryTake attempts to remove and return an element from the ConcurrentBag.
func (cb *ConcurrentBag[T]) TryTake() (T, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return popLast(&cb.container)
}

// ToArray copies the ConcurrentBag elements to a new slice.
func (cb *ConcurrentBag[T]) ToArray() []T {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	result := make([]T, len(cb.container))
	copy(result, cb.container)

	return result
}

// Clear removes all elements from the ConcurrentBag.
func (cb *ConcurrentBag[T]) Clear() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	clear(cb.container)
	cb.container = cb.container[:0]
}

// popLast removes and returns the last element from the slice.
// Returns the zero value and false if the slice is empty.
func popLast[T any](container *[]T) (T, bool) {
	n := len(*container)
	if n == 0 {
		var zero T
		return zero, false
	}

	element := (*container)[n-1]
	var zero T
	(*container)[n-1] = zero // clear reference for GC
	*container = (*container)[:n-1]
	return element, true
}
