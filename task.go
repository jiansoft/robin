package robin

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// task a struct
type task struct {
	fn          func() // fast path: no reflection needed
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

// newTask returns a task instance.
func newTask(f any, p ...any) task {
	if fn, ok := f.(func()); ok && len(p) == 0 {
		return task{fn: fn}
	}
	t := task{funcCache: reflect.ValueOf(f)}
	t.params(p...)

	return t
}

func (t *task) params(p ...any) {
	t.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		t.paramsCache[k] = reflect.ValueOf(param)
	}
}

func (t *task) execute() {
	if t.fn != nil {
		t.fn()
		return
	}
	_ = t.funcCache.Call(t.paramsCache)
}

type timerTask struct {
	scheduler    Scheduler
	task         task
	done         chan struct{}
	firstInMs    int64
	intervalInMs int64
	mu           sync.Mutex
	disposed     atomic.Bool
}

func newTimerTask(scheduler Scheduler, task task, firstInMs int64, intervalInMs int64) *timerTask {
	return &timerTask{
		scheduler:    scheduler,
		task:         task,
		firstInMs:    firstInMs,
		intervalInMs: intervalInMs,
		done:         make(chan struct{}),
	}
}

// Dispose releases the resources associated with the timerTask.
func (tk *timerTask) Dispose() {
	if !tk.disposed.CompareAndSwap(false, true) {
		return
	}

	close(tk.done)
	tk.mu.Lock()
	s := tk.scheduler
	tk.scheduler = nil
	tk.mu.Unlock()
	if s != nil {
		s.Remove(tk)
	}
}

// schedule starts the execution of the timerTask.
// If firstInMs is 0 or less, it executes the task immediately,
// otherwise it schedules the task to run after firstInMs milliseconds.
func (tk *timerTask) schedule() {
	if tk.firstInMs <= 0 {
		tk.next()
		return
	}

	go tk.runFirst()
}

// runFirst runs the task for the first time after the delay specified in firstInMs.
func (tk *timerTask) runFirst() {
	timer := time.NewTimer(time.Duration(tk.firstInMs) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-timer.C:
		tk.next()
	case <-tk.done:
	}
}

// next schedules the next execution of the task.
// If the task is not recurring (intervalInMs is 0 or less), it disposes the task,
// otherwise it schedules the task to run again after intervalInMs milliseconds.
func (tk *timerTask) next() {
	tk.executeOnFiber()
	if tk.intervalInMs <= 0 {
		tk.Dispose()
		return
	}

	if tk.disposed.Load() {
		return
	}

	go tk.runInterval()
}

// runInterval runs the task at regular intervals specified by intervalInMs.
func (tk *timerTask) runInterval() {
	ticker := time.NewTicker(time.Duration(tk.intervalInMs) * time.Millisecond)
	defer ticker.Stop()
	for !tk.disposed.Load() {
		select {
		case <-ticker.C:
			tk.executeOnFiber()
		case <-tk.done:
			return
		}
	}
}

// executeOnFiber enqueues the task in the scheduler for execution.
func (tk *timerTask) executeOnFiber() {
	if tk.disposed.Load() {
		return
	}

	tk.mu.Lock()
	if tk.scheduler != nil {
		tk.scheduler.enqueueTask(tk.task)
	}
	tk.mu.Unlock()
}
