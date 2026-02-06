package robin

import (
	"sync"
	"sync/atomic"
)

// Disposable an interface just has only one function
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

type scheduler struct {
	fiber    Fiber
	tasks    sync.Map
	disposed atomic.Bool
}

func newScheduler(fiber Fiber) *scheduler {
	return &scheduler{fiber: fiber}
}

// Schedule delay n millisecond then execute once the function
func (s *scheduler) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFunc, params...)
}

// ScheduleOnInterval first time delay N millisecond then execute once the function,
// then interval N millisecond repeat execute the function.
func (s *scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	pending := newTimerTask(s, newTask(taskFunc, params...), firstInMs, regularInMs)
	if s.disposed.Load() {
		return pending
	}
	s.tasks.Store(pending, pending)
	pending.schedule()
	return pending
}

// Enqueue enqueues a task for execution.
func (s *scheduler) Enqueue(taskFunc any, params ...any) {
	s.enqueueTask(newTask(taskFunc, params...))
}

func (s *scheduler) enqueueTask(task task) {
	s.fiber.enqueueTask(task)
}

// Remove removes a disposable from the scheduler.
func (s *scheduler) Remove(d Disposable) {
	s.tasks.Delete(d)
}

// Dispose disposes of all scheduled tasks.
func (s *scheduler) Dispose() {
	s.disposed.Store(true)
	s.tasks.Range(func(k, v any) bool {
		s.tasks.Delete(k)
		v.(Disposable).Dispose()
		return true
	})
}
