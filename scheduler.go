package robin

import "sync"

// Disposable an interface just has only one function
type Disposable interface {
	Dispose()
}

// IScheduler an interface that for GoroutineMulti and GoroutineSingle use.
type IScheduler interface {
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
	Enqueue(taskFunc any, params ...any)
	enqueueTask(task task)
	Remove(d Disposable)
	Dispose()
}

type scheduler struct {
	fiber IFiber
	sync.Map
	running   bool
	isDispose bool
}

func newScheduler(fiber IFiber) *scheduler {
	s := new(scheduler)
	s.fiber = fiber
	s.running = true
	return s
}

// Schedule delay n millisecond then execute once the function
func (s *scheduler) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFunc, params...)
}

// ScheduleOnInterval first time delay N millisecond then execute once the function,
// then interval N millisecond repeat execute the function.
func (s *scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	pending := newTimerTask(s, newTask(taskFunc, params...), firstInMs, regularInMs)
	if s.isDispose {
		return pending
	}
	s.Store(pending, pending)
	pending.schedule()
	return pending
}

// Enqueue Implement SchedulerRegistry.Enqueue
func (s *scheduler) Enqueue(taskFunc any, params ...any) {
	s.enqueueTask(newTask(taskFunc, params...))
}

func (s *scheduler) enqueueTask(task task) {
	s.fiber.enqueueTask(task)
}

// Remove Implement SchedulerRegistry.Forget
func (s *scheduler) Remove(d Disposable) {
	s.Delete(d)
}

func (s *scheduler) Dispose() {
	s.isDispose = true
	s.Range(func(k, v any) bool {
		s.Delete(k)
		v.(Disposable).Dispose()
		return true
	})
}
