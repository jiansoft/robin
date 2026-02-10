package robin

import (
	"sync"
	"sync/atomic"
)

// Disposable is an interface for releasing resources.
type Disposable interface {
	Dispose()
}

// Scheduler is an interface for scheduling and managing tasks on fibers.
type Scheduler interface {
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
	Enqueue(taskFunc any, params ...any)
	enqueueTask(task task)
	Remove(d Disposable)
	Dispose()
}

// scheduler implements the Scheduler interface using sync.Map for task tracking.
type scheduler struct {
	fiber    Fiber
	tasks    sync.Map
	disposed atomic.Bool
}

// newScheduler creates a new scheduler bound to the given fiber.
func newScheduler(fiber Fiber) *scheduler {
	return &scheduler{fiber: fiber}
}

// Schedule executes the task once after firstInMs milliseconds.
func (s *scheduler) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFunc, params...)
}

// ScheduleOnInterval executes the task after firstInMs milliseconds,
// then repeats every regularInMs milliseconds.
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

// Enqueue enqueues a task for execution.
func (s *scheduler) Enqueue(taskFunc any, params ...any) {
	s.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask forwards a task to the underlying fiber for execution.
func (s *scheduler) enqueueTask(task task) {
	s.fiber.enqueueTask(task)
}

// Remove removes a disposable from the scheduler.
func (s *scheduler) Remove(d Disposable) {
	s.tasks.Delete(d)
}

// Dispose disposes of all scheduled tasks.
// Safe to call multiple times; only the first call takes effect.
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
