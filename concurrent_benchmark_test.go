package robin

import (
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// oldConcurrentQueue: pre-optimization ring buffer (modulo, mutex Len, no shrink,
// Clear loops to zero). Generic to isolate ring buffer structural differences only.
// ---------------------------------------------------------------------------

type oldConcurrentQueue[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
	mu    sync.Mutex
}

func newOldConcurrentQueue[T any]() *oldConcurrentQueue[T] {
	return &oldConcurrentQueue[T]{
		buf: make([]T, 16),
	}
}

func (q *oldConcurrentQueue[T]) Enqueue(element T) {
	q.mu.Lock()
	if q.count == len(q.buf) {
		newBuf := make([]T, len(q.buf)*2)
		if q.head < q.tail {
			copy(newBuf, q.buf[q.head:q.tail])
		} else {
			n := copy(newBuf, q.buf[q.head:])
			copy(newBuf[n:], q.buf[:q.tail])
		}
		q.head = 0
		q.tail = q.count
		q.buf = newBuf
	}
	q.buf[q.tail] = element
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
	q.mu.Unlock()
}

func (q *oldConcurrentQueue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		var zero T
		return zero, false
	}
	element := q.buf[q.head]
	var zero T
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	return element, true
}

func (q *oldConcurrentQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count
}

func (q *oldConcurrentQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	var zero T
	for i := range q.buf {
		q.buf[i] = zero
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue only
// ---------------------------------------------------------------------------

func BenchmarkOldQueue_Enqueue(b *testing.B) {
	q := newOldConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
	}
}

func BenchmarkNewQueue_Enqueue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue + Dequeue (steady state, no grow)
// ---------------------------------------------------------------------------

func BenchmarkOldQueue_EnqueueDequeue(b *testing.B) {
	q := newOldConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
		q.TryDequeue()
	}
}

func BenchmarkNewQueue_EnqueueDequeue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
		q.TryDequeue()
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Len (mutex vs atomic)
// ---------------------------------------------------------------------------

func BenchmarkOldQueue_Len(b *testing.B) {
	q := newOldConcurrentQueue[int]()
	for i := range 1000 {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for range b.N {
		_ = q.Len()
	}
}

func BenchmarkNewQueue_Len(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for i := range 1000 {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for range b.N {
		_ = q.Len()
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Concurrent (parallel enqueue/dequeue)
// ---------------------------------------------------------------------------

func BenchmarkOldQueue_Concurrent(b *testing.B) {
	q := newOldConcurrentQueue[int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i&1 == 0 {
				q.Enqueue(i)
			} else {
				q.TryDequeue()
			}
			i++
		}
	})
}

func BenchmarkNewQueue_Concurrent(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i&1 == 0 {
				q.Enqueue(i)
			} else {
				q.TryDequeue()
			}
			i++
		}
	})
}

// ---------------------------------------------------------------------------
// Benchmarks: Clear after filling 1000 elements
// ---------------------------------------------------------------------------

func BenchmarkOldQueue_Clear(b *testing.B) {
	q := newOldConcurrentQueue[int]()
	for range b.N {
		for j := range 1000 {
			q.Enqueue(j)
		}
		q.Clear()
	}
}

func BenchmarkNewQueue_Clear(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for range b.N {
		for j := range 1000 {
			q.Enqueue(j)
		}
		q.Clear()
	}
}
