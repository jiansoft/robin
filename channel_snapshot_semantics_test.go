package robin

import (
	"sync"
	"sync/atomic"
	"testing"
)

// manualFiber is a deterministic test fiber that buffers tasks and executes them on demand.
// manualFiber 是可控的測試 fiber：先緩存任務，待測試主動觸發執行。
type manualFiber struct {
	mu       sync.Mutex
	disposed bool
	tasks    []task
}

func newManualFiber() *manualFiber {
	return &manualFiber{tasks: make([]task, 0)}
}

func (f *manualFiber) Dispose() {
	f.mu.Lock()
	f.disposed = true
	f.tasks = nil
	f.mu.Unlock()
}

func (f *manualFiber) Enqueue(taskFunc any, params ...any) {
	f.enqueueTask(newTask(taskFunc, params...))
}

func (f *manualFiber) enqueueTask(t task) {
	f.mu.Lock()
	if !f.disposed {
		f.tasks = append(f.tasks, t)
	}
	f.mu.Unlock()
}

func (f *manualFiber) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return noopDisposable{}
}

func (f *manualFiber) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return noopDisposable{}
}

func (f *manualFiber) runAll() {
	for {
		f.mu.Lock()
		if len(f.tasks) == 0 {
			f.mu.Unlock()
			return
		}
		batch := f.tasks
		f.tasks = make([]task, 0)
		f.mu.Unlock()

		for i := range batch {
			batch[i].execute()
			batch[i] = task{}
		}
	}
}

type noopDisposable struct{}

func (noopDisposable) Dispose() {}

// TestChannelPublishUsesCallTimeSnapshot verifies publish dispatch is based on
// subscribers present at Publish call time (not at task execution time).
// TestChannelPublishUsesCallTimeSnapshot 驗證 Publish 以呼叫當下的訂閱者快照為準，
// 而非以實際執行任務當下的訂閱者集合為準。
func TestChannelPublishUsesCallTimeSnapshot(t *testing.T) {
	f := newManualFiber()
	ch := NewChannel()
	ch.fiber = f

	var aCount atomic.Int32
	var bCount atomic.Int32

	subA := ch.Subscribe(func(_ string) {
		aCount.Add(1)
	})
	ch.Publish("m1")
	ch.Subscribe(func(_ string) {
		bCount.Add(1)
	})

	// m1 should see only A because B subscribed after Publish call.
	f.runAll()
	if aCount.Load() != 1 || bCount.Load() != 0 {
		t.Fatalf("after m1: a=%d b=%d, want a=1 b=0", aCount.Load(), bCount.Load())
	}

	ch.Publish("m2")
	subA.Unsubscribe()

	// m2 should still include A because A was present at Publish call time.
	f.runAll()
	if aCount.Load() != 2 || bCount.Load() != 1 {
		t.Fatalf("after m2: a=%d b=%d, want a=2 b=1", aCount.Load(), bCount.Load())
	}

	ch.Publish("m3")
	f.runAll()

	// m3 should no longer include A because A was unsubscribed before this Publish call.
	if aCount.Load() != 2 || bCount.Load() != 2 {
		t.Fatalf("after m3: a=%d b=%d, want a=2 b=2", aCount.Load(), bCount.Load())
	}
}

// TestTypedChannelPublishUsesCallTimeSnapshot verifies the same call-time snapshot
// behavior for TypedChannel.
// TestTypedChannelPublishUsesCallTimeSnapshot 驗證 TypedChannel 同樣採用呼叫當下快照語義。
func TestTypedChannelPublishUsesCallTimeSnapshot(t *testing.T) {
	f := newManualFiber()
	ch := NewTypedChannel[int]()
	ch.fiber = f

	var aCount atomic.Int32
	var bCount atomic.Int32

	subA := ch.Subscribe(func(_ int) {
		aCount.Add(1)
	})
	ch.Publish(1)
	ch.Subscribe(func(_ int) {
		bCount.Add(1)
	})

	f.runAll()
	if aCount.Load() != 1 || bCount.Load() != 0 {
		t.Fatalf("after m1: a=%d b=%d, want a=1 b=0", aCount.Load(), bCount.Load())
	}

	ch.Publish(2)
	subA.Unsubscribe()

	f.runAll()
	if aCount.Load() != 2 || bCount.Load() != 1 {
		t.Fatalf("after m2: a=%d b=%d, want a=2 b=1", aCount.Load(), bCount.Load())
	}

	ch.Publish(3)
	f.runAll()

	if aCount.Load() != 2 || bCount.Load() != 2 {
		t.Fatalf("after m3: a=%d b=%d, want a=2 b=2", aCount.Load(), bCount.Load())
	}
}
