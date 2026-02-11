package robin

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue only
// 基準測試：僅入列
// ---------------------------------------------------------------------------

// BenchmarkQueue_Enqueue measures pure enqueue throughput on ConcurrentQueue.
// BenchmarkQueue_Enqueue 量測 ConcurrentQueue 純入列操作的吞吐量。
func BenchmarkQueue_Enqueue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	b.ResetTimer()
	for i := range b.N {
		q.Enqueue(i)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Enqueue + Dequeue (steady state, no grow)
// 基準測試：入列+出列（穩態，不觸發擴容）
// ---------------------------------------------------------------------------

// BenchmarkQueue_EnqueueDequeue measures balanced enqueue/dequeue cost in steady state.
// BenchmarkQueue_EnqueueDequeue 量測穩態下入列與出列的平衡操作成本。
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
// 基準測試：Len 讀取
// ---------------------------------------------------------------------------

// BenchmarkQueue_Len measures lock/read overhead of Len() under pre-filled queue.
// BenchmarkQueue_Len 量測預填充佇列下 Len() 的鎖與讀取開銷。
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
// 基準測試：並發入列/出列
// ---------------------------------------------------------------------------

// BenchmarkQueue_Concurrent measures mixed enqueue/dequeue behavior under b.RunParallel.
// BenchmarkQueue_Concurrent 量測 b.RunParallel 下混合入列/出列的並發表現。
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
// 基準測試：填入 1000 筆後清空
// ---------------------------------------------------------------------------

// BenchmarkQueue_Clear measures Clear() cost after filling the queue with 1000 items.
// BenchmarkQueue_Clear 量測每輪先填入 1000 筆後執行 Clear() 的成本。
func BenchmarkQueue_Clear(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for range b.N {
		for j := range 1000 {
			q.Enqueue(j)
		}
		q.Clear()
	}
}
