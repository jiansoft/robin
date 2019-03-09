package robin

import (
	"reflect"
	"sync/atomic"
	"time"
)

//Task a struct
type Task struct {
	doFunc      interface{}
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func newTask(f interface{}, p ...interface{}) Task {
	var paramLen = len(p)
	task := Task{doFunc: f, funcCache: reflect.ValueOf(f), paramsCache: make([]reflect.Value, paramLen)}
	for k, param := range p {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t Task) execute() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}

type timerTask struct {
	scheduler    IScheduler
	firstInMs    int64
	intervalInMs int64
	task         Task
	disposed     int32
	exitC        chan bool
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
}

func (t *timerTask) schedule() {
	firstInMs := atomic.LoadInt64(&t.firstInMs)
	if firstInMs <= 0 {
		t.next()
		return
	}

	go func(first time.Duration, exitC chan bool) {
		select {
		case <-time.After(first):
			t.next()
		case <-exitC:
			return
		}
	}(time.Duration(firstInMs)*time.Millisecond, t.exitC)
}

func (t *timerTask) next() {
	t.executeOnFiber()
	intervalInMs := atomic.LoadInt64(&t.intervalInMs)
	if intervalInMs <= 0 {
		t.Dispose()
		return
	}

	ticker := time.NewTicker(time.Duration(intervalInMs) * time.Millisecond)
	go func(ticker *time.Ticker, exitC chan bool) {
		for atomic.LoadInt32(&t.disposed) == 0 {
			select {
			case <-ticker.C:
				t.executeOnFiber()
			case <-exitC:
				break
			}
		}
		ticker.Stop()
	}(ticker, t.exitC)
}

func (t *timerTask) executeOnFiber() {
	if atomic.LoadInt32(&t.disposed) == 1 {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}
