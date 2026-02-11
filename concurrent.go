package robin

import (
	"sync"
)

const defaultRingCapacity = 16

// ConcurrentQueue represents a thread-safe first in-first out (FIFO) collection.
// Uses a power-of-two ring buffer so index wrapping can use bitwise masking for O(1) math.
// Capacity grows when full and may shrink when utilization drops below 25%.
// ConcurrentQueue 是執行緒安全的先進先出（FIFO）集合。
// 透過 2 的冪次容量環形緩衝區，以位元遮罩完成 O(1) 索引計算。
// 佇列滿時會擴容，使用率低於 25% 時可能縮容。
type ConcurrentQueue[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
	mask  int // len(buf) - 1, bitwise AND replaces modulo
	mu    sync.Mutex
}

// NewConcurrentQueue creates an empty queue with default ring capacity.
// NewConcurrentQueue 會建立預設容量的空佇列。
func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		buf:  make([]T, defaultRingCapacity),
		mask: defaultRingCapacity - 1,
	}
}

// Enqueue appends one element to the tail of the queue.
// Enqueue 會將元素加入佇列尾端。
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

// grow doubles ring capacity while preserving FIFO order.
// grow 會把環形緩衝區容量加倍，並保留原有 FIFO 順序。
func (q *ConcurrentQueue[T]) grow() {
	q.resize((q.mask + 1) << 1)
}

// shrink halves ring capacity when utilization is low, but never below default capacity.
// shrink 在低使用率時把容量減半，但不會低於預設容量。
func (q *ConcurrentQueue[T]) shrink() {
	curCap := q.mask + 1
	if curCap <= defaultRingCapacity || q.count >= curCap>>2 {
		return
	}
	// curCap is always a power of two > defaultRingCapacity,
	// so curCap>>1 >= defaultRingCapacity is guaranteed.
	q.resize(curCap >> 1)
}

// resize rebuilds the ring buffer into a new capacity and re-linearizes current elements.
// resize 會以新容量重建環形緩衝區，並把現有元素重新線性排列。
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

// TryPeek returns the front element without removing it.
// The bool result is false when the queue is empty.
// TryPeek 會回傳佇列前端元素但不移除。
// 當佇列為空時，bool 會是 false。
func (q *ConcurrentQueue[T]) TryPeek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}

	return q.buf[q.head], true
}

// TryDequeue removes and returns the front element.
// The bool result is false when the queue is empty.
// TryDequeue 會移除並回傳佇列前端元素。
// 當佇列為空時，bool 會是 false。
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

// Len returns the number of currently stored elements.
// Len 回傳目前儲存的元素數量。
func (q *ConcurrentQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.count
}

// Clear removes all elements and resets queue storage to default capacity.
// Clear 會清空所有元素，並把佇列容量重設為預設值。
func (q *ConcurrentQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.buf = make([]T, defaultRingCapacity)
	q.head = 0
	q.tail = 0
	q.count = 0
	q.mask = defaultRingCapacity - 1
}

// ToArray copies queue contents into a new slice in FIFO order.
// ToArray 會以 FIFO 順序將佇列內容複製到新的 slice。
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
// ConcurrentStack 是執行緒安全的後進先出（LIFO）集合。
type ConcurrentStack[T any] struct {
	container []T
	mu        sync.Mutex
}

// NewConcurrentStack creates an empty stack.
// NewConcurrentStack 會建立空堆疊。
func NewConcurrentStack[T any]() *ConcurrentStack[T] {
	return &ConcurrentStack[T]{}
}

// Push adds an element to the top of the stack.
// Push 會把元素推入堆疊頂端。
func (s *ConcurrentStack[T]) Push(element T) {
	s.mu.Lock()
	s.container = append(s.container, element)
	s.mu.Unlock()
}

// TryPeek returns the top element without popping it.
// The bool result is false when the stack is empty.
// TryPeek 會回傳堆疊頂端元素但不彈出。
// 當堆疊為空時，bool 會是 false。
func (s *ConcurrentStack[T]) TryPeek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.container) == 0 {
		var zero T
		return zero, false
	}

	return s.container[len(s.container)-1], true
}

// TryPop pops and returns the current top element.
// The bool result is false when the stack is empty.
// TryPop 會彈出並回傳目前堆疊頂端元素。
// 當堆疊為空時，bool 會是 false。
func (s *ConcurrentStack[T]) TryPop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return popLast(&s.container)
}

// Len returns the current number of stack elements.
// Len 回傳目前堆疊元素數量。
func (s *ConcurrentStack[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.container)
}

// Clear removes all elements from the stack.
// Clear 會清空堆疊中的所有元素。
func (s *ConcurrentStack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.container)
	s.container = s.container[:0]
}

// ToArray copies the elements stored in the ConcurrentStack to a new slice.
// Elements are returned in LIFO order (top of stack first).
// ToArray 會把堆疊元素複製到新 slice，順序為 LIFO（頂端在前）。
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

// ConcurrentBag represents a thread-safe, unordered collection.
// ConcurrentBag 是執行緒安全且不保證順序的集合。
type ConcurrentBag[T any] struct {
	container []T
	mu        sync.Mutex
}

// NewConcurrentBag creates an empty bag with a small pre-allocated capacity.
// NewConcurrentBag 會建立含小型預配置容量的空 bag。
func NewConcurrentBag[T any]() *ConcurrentBag[T] {
	return &ConcurrentBag[T]{
		container: make([]T, 0, 64),
	}
}

// Add inserts one element into the bag.
// Add 會把一個元素加入 bag。
func (cb *ConcurrentBag[T]) Add(element T) {
	cb.mu.Lock()
	cb.container = append(cb.container, element)
	cb.mu.Unlock()
}

// Len returns the number of currently stored bag elements.
// Len 回傳目前 bag 內元素數量。
func (cb *ConcurrentBag[T]) Len() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return len(cb.container)
}

// TryTake removes and returns one element if available.
// The bool result is false when the bag is empty.
// TryTake 會移除並回傳一個元素；若 bag 為空則 bool 為 false。
func (cb *ConcurrentBag[T]) TryTake() (T, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return popLast(&cb.container)
}

// ToArray copies all bag elements into a new slice.
// Element order is unspecified.
// ToArray 會把 bag 所有元素複製到新 slice。
// 元素順序不保證固定。
func (cb *ConcurrentBag[T]) ToArray() []T {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	result := make([]T, len(cb.container))
	copy(result, cb.container)

	return result
}

// Clear removes all elements from the bag.
// Clear 會清空 bag 內所有元素。
func (cb *ConcurrentBag[T]) Clear() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	clear(cb.container)
	cb.container = cb.container[:0]
}

// popLast removes and returns the last element from the slice.
// Returns the zero value and false if the slice is empty.
// popLast 會移除並回傳切片最後一個元素。
// 若切片為空，則回傳 zero value 與 false。
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
