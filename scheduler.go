package robin

import "sync"

// Disposable an interface just has only one function
type Disposable interface {
	Dispose()
}

// IScheduler an interface that for GoroutineMulti and GoroutineSingle use.
type IScheduler interface {
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
	Remove(d Disposable)
	Dispose()
}

type scheduler struct {
	fiber     Fiber
	running   bool
	isDispose bool
	sync.Map
}

func newScheduler(executionState Fiber) *scheduler {
	s := new(scheduler)
	s.fiber = executionState
	s.running = true
	return s
}

// Schedule delay n millisecond then execute once the function
func (s *scheduler) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFun, params...)
}

// ScheduleOnInterval first time delay N millisecond then execute once the function,
// then interval N millisecond repeat execute the function.
func (s *scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	pending := newTimerTask(s, newTask(taskFun, params...), firstInMs, regularInMs)
	if s.isDispose {
		return pending
	}
	s.Store(pending, pending)
	pending.schedule()
	return pending
}

//Implement SchedulerRegistry.Enqueue
func (s *scheduler) Enqueue(taskFun interface{}, params ...interface{}) {
	s.EnqueueWithTask(newTask(taskFun, params...))
}

func (s *scheduler) EnqueueWithTask(task Task) {
	s.fiber.EnqueueWithTask(task)
}

//Implement SchedulerRegistry.Remove
func (s *scheduler) Remove(d Disposable) {
	s.Delete(d)
}

func (s *scheduler) Dispose() {
	s.isDispose = true
	s.Range(func(k, v interface{}) bool {
		s.Delete(k)
		v.(Disposable).Dispose()
		return true
	})
}
