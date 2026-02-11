package robin

import (
	"sync"
	"sync/atomic"
)

// Disposable represents a cancellable or releasable runtime resource.
// Disposable 表示可取消或可釋放的執行期資源。
type Disposable interface {
	Dispose()
}

// Scheduler defines task scheduling and lifecycle management operations on a fiber.
// Scheduler 定義在 fiber 上進行任務排程與生命週期管理的操作集合。
type Scheduler interface {
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
	Enqueue(taskFunc any, params ...any)
	enqueueTask(task task)
	Remove(d Disposable)
	Dispose()
}

// scheduler is the default Scheduler implementation backed by sync.Map.
// It tracks timer tasks so Dispose can cancel all outstanding schedules safely.
// scheduler 是預設的 Scheduler 實作，底層以 sync.Map 追蹤計時任務。
// 透過此追蹤可在 Dispose 時安全取消所有仍在排程中的任務。
type scheduler struct {
	fiber    Fiber
	tasks    sync.Map
	disposed atomic.Bool
}

// newScheduler creates a scheduler bound to a specific fiber instance.
// newScheduler 會建立並綁定到指定 fiber 的排程器。
func newScheduler(fiber Fiber) *scheduler {
	return &scheduler{fiber: fiber}
}

// Schedule runs a one-shot task after firstInMs milliseconds.
// Schedule 會在 firstInMs 毫秒後執行一次性任務。
func (s *scheduler) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFunc, params...)
}

// ScheduleOnInterval executes the task after firstInMs milliseconds,
// then repeats every regularInMs milliseconds.
// ScheduleOnInterval 會在 firstInMs 毫秒後首次執行，之後每 regularInMs 毫秒重複執行。
func (s *scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	pending := newTimerTask(s, newTask(taskFunc, params...), firstInMs, regularInMs)
	if s.disposed.Load() {
		return pending
	}
	// Store first, then re-check: if Dispose() raced in between,
	// its Range will see this entry and clean it up. If Dispose()
	// already completed, we clean up ourselves.
	s.tasks.Store(pending, pending)
	if s.disposed.Load() {
		s.tasks.Delete(pending)
		pending.Dispose()
		return pending
	}
	pending.schedule()
	return pending
}

// Enqueue forwards an immediate task to the underlying fiber queue.
// Enqueue 會把立即執行的任務轉送到底層 fiber 佇列。
func (s *scheduler) Enqueue(taskFunc any, params ...any) {
	s.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask pushes an already-built task object to the target fiber.
// enqueueTask 會把已建立好的 task 物件推送至目標 fiber。
func (s *scheduler) enqueueTask(task task) {
	s.fiber.enqueueTask(task)
}

// Remove forgets a disposable from scheduler tracking.
// This method only removes map bookkeeping and does not call Dispose on d.
// Remove 只會把 disposable 從排程器追蹤表移除，
// 不會主動呼叫 d.Dispose()。
func (s *scheduler) Remove(d Disposable) {
	s.tasks.Delete(d)
}

// Dispose disposes of all scheduled tasks.
// Safe to call multiple times; only the first call takes effect.
// Dispose 會釋放所有已排程任務。
// 可重複呼叫；僅第一次會生效。
func (s *scheduler) Dispose() {
	if !s.disposed.CompareAndSwap(false, true) {
		return
	}
	s.tasks.Range(func(k, v any) bool {
		s.tasks.Delete(k)
		if d, ok := v.(Disposable); ok {
			d.Dispose()
		}
		return true
	})
}
