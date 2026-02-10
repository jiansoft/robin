package robin

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue only
// ---------------------------------------------------------------------------

func BenchmarkQueue_Enqueue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue + Dequeue (steady state, no grow)
// ---------------------------------------------------------------------------

func BenchmarkQueue_EnqueueDequeue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
		q.TryDequeue()
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Len
// ---------------------------------------------------------------------------

func BenchmarkQueue_Len(b *testing.B) {
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

func BenchmarkQueue_Concurrent(b *testing.B) {
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

func BenchmarkQueue_Clear(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for range b.N {
		for j := range 1000 {
			q.Enqueue(j)
		}
		q.Clear()
	}
}
