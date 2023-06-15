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

// Dispose release resources
func (tk *timerTask) Dispose() {
	if !atomic.CompareAndSwapInt32(&tk.disposed, 0, 1) {
		return
	}

	tk.scheduler.Remove(tk)
	close(tk.exitC)
	tk.scheduler = nil
}

func (tk *timerTask) schedule() {
	firstInMs := atomic.LoadInt64(&tk.firstInMs)
	if firstInMs <= 0 {
		tk.next()
		return
	}

	tk.scheduler.Enqueue(func(tk *timerTask) {
		first := time.NewTimer(time.Duration(tk.firstInMs) * time.Millisecond)
		select {
		case <-first.C:
			tk.next()
		case <-tk.exitC:
		}

		if !first.Stop() {
			<-first.C
		}
	}, tk)
}

func (tk *timerTask) next() {
	tk.executeOnFiber()
	if atomic.LoadInt64(&tk.intervalInMs) <= 0 {
		tk.Dispose()
		return
	}

	if atomic.LoadInt32(&tk.disposed) == 1 {
		return
	}

	tk.scheduler.Enqueue(func(tk *timerTask) {
		ticker := time.NewTicker(time.Duration(tk.intervalInMs) * time.Millisecond)
		for atomic.LoadInt32(&tk.disposed) == 0 {
			select {
			case <-ticker.C:
				tk.executeOnFiber()
			case <-tk.exitC:
				ticker.Stop()
				break
			}
		}
	}, tk)
}

func (tk *timerTask) executeOnFiber() {
	if atomic.LoadInt32(&tk.disposed) == 1 {
		return
	}
	tk.scheduler.EnqueueWithTask(tk.task)
}
