package robin

import (
	"reflect"
	"sync/atomic"
	"time"
)

// Task a struct
type Task struct {
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

// newTask returns a Task instance.
func newTask(f any, p ...any) Task {
	task := Task{funcCache: reflect.ValueOf(f)}
	task.params(p...)
	return task
}

func (t *Task) params(p ...any) {
	var paramLen = len(p)
	t.paramsCache = make([]reflect.Value, paramLen)
	for k, param := range p {
		t.paramsCache[k] = reflect.ValueOf(param)
	}
}

func (t *Task) execute() {
	_ = t.funcCache.Call(t.paramsCache)
	//func(f reflect.Value, in []reflect.Value) { _ = f.Call(in) }(t.funcCache, t.paramsCache)
}

type timerTask struct {
	scheduler    IScheduler
	exitC        chan bool
	task         Task
	firstInMs    int64
	intervalInMs int64
	disposed     int32
}

func newTimerTask(scheduler IScheduler, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	var t = &timerTask{scheduler: scheduler, task: task, firstInMs: firstInMs, intervalInMs: intervalInMs, exitC: make(chan bool)}
	return t
}

// Dispose releases the resources associated with the timerTask.
// It closes the exitC channel, removes the task from the scheduler, and sets the scheduler to nil.
func (tk *timerTask) Dispose() {
	if !atomic.CompareAndSwapInt32(&tk.disposed, 0, 1) {
		return
	}

	tk.scheduler.Remove(tk)
	close(tk.exitC)
	tk.scheduler = nil
}

// schedule starts the execution of the timerTask.
// If firstInMs is 0 or less, it executes the task immediately,
// otherwise it schedules the task to run after firstInMs milliseconds.
func (tk *timerTask) schedule() {
	if atomic.LoadInt64(&tk.firstInMs) <= 0 {
		tk.next()
		return
	}

	go tk.runFirst()
}

// runFirst runs the task for the first time after the delay specified in firstInMs.
// If the task is signaled to stop through exitC before the delay expires, it stops the timer and returns.
func (tk *timerTask) runFirst() {
	firstDuration := time.Duration(tk.firstInMs) * time.Millisecond
	firstRun := time.NewTimer(firstDuration)
	select {
	case <-firstRun.C:
		tk.next()
	case <-tk.exitC:
		if !firstRun.Stop() {
			<-firstRun.C
		}
	}
}

// next schedules the next execution of the task.
// If the task is not recurring (intervalInMs is 0 or less), it disposes the task,
// otherwise it schedules the task to run again after intervalInMs milliseconds.
func (tk *timerTask) next() {
	tk.executeOnFiber()
	if atomic.LoadInt64(&tk.intervalInMs) <= 0 {
		tk.Dispose()
		return
	}

	if atomic.LoadInt32(&tk.disposed) == 1 {
		return
	}

	go tk.runInterval()
}

// runInterval runs the task at regular intervals specified by intervalInMs.
// If the task is signaled to stop through exitC, it stops the timer and returns.
func (tk *timerTask) runInterval() {
	intervalDuration := time.Duration(tk.intervalInMs) * time.Millisecond
	intervalRun := time.NewTicker(intervalDuration)
	for atomic.LoadInt32(&tk.disposed) == 0 {
		select {
		case <-intervalRun.C:
			tk.executeOnFiber()
		case <-tk.exitC:
			intervalRun.Stop()
			break
		}
	}
}

// executeOnFiber enqueues the task in the scheduler for execution.
// If the task is disposed, it does nothing.
func (tk *timerTask) executeOnFiber() {
	if atomic.LoadInt32(&tk.disposed) == 1 {
		return
	}
	tk.scheduler.EnqueueWithTask(tk.task)
}
