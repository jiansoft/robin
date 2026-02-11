package robin

import (
	"sync"
)

// taskQueue defines the internal queue contract used by fiber implementations.
// taskQueue 定義 fiber 實作所需的內部佇列介面契約。
type taskQueue interface {
	count() int
	dequeueAll() ([]task, bool)
	dispose()
	enqueue(t task)
}

// defaultQueue is a double-buffer queue:
// pending collects newly enqueued tasks, and upcoming is returned as the current batch on dequeue.
// defaultQueue 是雙緩衝佇列：
// pending 收集新進任務，dequeue 時會把 upcoming 作為目前批次回傳。
type defaultQueue struct {
	pending  []task
	upcoming []task
	mu       sync.Mutex
}

// newDefaultQueue creates a new defaultQueue with pre-allocated empty slices.
// Both slices are initialized to non-nil so that nil exclusively represents the disposed state.
// newDefaultQueue 會建立並預配置空切片。
// 兩個切片都保證非 nil，讓 nil 可專門表示「已 dispose」狀態。
func newDefaultQueue() *defaultQueue {
	return &defaultQueue{
		pending:  make([]task, 0),
		upcoming: make([]task, 0),
	}
}

// dispose marks the queue as disposed by setting buffers to nil.
// dispose 會把兩個緩衝切片設為 nil，表示佇列已釋放不可再用。
func (dq *defaultQueue) dispose() {
	dq.mu.Lock()
	dq.pending = nil
	dq.upcoming = nil
	dq.mu.Unlock()
}

// enqueue appends a task into pending unless the queue is already disposed.
// enqueue 會把任務加入 pending；若佇列已 dispose 則為 no-op。
func (dq *defaultQueue) enqueue(t task) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if dq.pending == nil { // already disposed
		return
	}
	dq.pending = append(dq.pending, t)
}

// dequeueAll swaps pending/upcoming buffers and returns the batch snapshot.
// It also zeroes old task slots to release references promptly for GC.
// dequeueAll 會交換 pending/upcoming 並回傳批次快照。
// 同時會把舊任務槽位清零，避免長時間持有引用而延後 GC 回收。
func (dq *defaultQueue) dequeueAll() ([]task, bool) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	// Swap: pending becomes upcoming (the batch to return),
	// old upcoming becomes the new pending buffer.
	dq.upcoming, dq.pending = dq.pending, dq.upcoming

	// Clear references in the old upcoming (now pending) to avoid holding
	// pointers (fn, invoke, invokeArg, funcCache, paramsCache) that prevent GC.
	for i := range dq.pending {
		dq.pending[i] = task{}
	}
	dq.pending = dq.pending[:0]

	return dq.upcoming, len(dq.upcoming) > 0
}

// count returns how many tasks are currently waiting in pending.
// count 回傳目前 pending 中等待處理的任務數量。
func (dq *defaultQueue) count() int {
	dq.mu.Lock()
	count := len(dq.pending)
	dq.mu.Unlock()

	return count
}
