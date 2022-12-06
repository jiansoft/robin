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
func (t *timerTask) Dispose() {
	if !atomic.CompareAndSwapInt32(&t.disposed, 0, 1) {
		return
	}

	t.scheduler.Remove(t)
	close(t.exitC)
	t.scheduler = nil
}

func (t *timerTask) schedule() {
	firstInMs := atomic.LoadInt64(&t.firstInMs)
	if firstInMs <= 0 {
		t.next()
		return
	}

	go func(ft int64, tk *timerTask) {
		first := time.NewTimer(time.Duration(ft) * time.Millisecond)
		select {
		case <-first.C:
			tk.next()
		case <-tk.exitC:
		}
		first.Stop()
	}(firstInMs, t)
}

func (t *timerTask) next() {
	t.executeOnFiber()
	intervalInMs := atomic.LoadInt64(&t.intervalInMs)
	if intervalInMs <= 0 || atomic.LoadInt32(&t.disposed) == 1 {
		t.Dispose()
		return
	}

	go func(interval int64, tk *timerTask) {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		for atomic.LoadInt32(&t.disposed) == 0 {
			select {
			case <-ticker.C:
				t.executeOnFiber()
			case <-tk.exitC:
				break
			}
		}
		ticker.Stop()
	}(intervalInMs, t)
}

func (t *timerTask) executeOnFiber() {
	if atomic.LoadInt32(&t.disposed) == 1 {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}
