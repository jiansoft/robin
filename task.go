package robin

import (
	"reflect"
	"sync"
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
	task := Task{doFunc: f, funcCache: reflect.ValueOf(f)}
	var paramLen = len(p)
	task.paramsCache = make([]reflect.Value, paramLen)
	for k, param := range p {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t Task) run() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}

type timerTask struct {
	scheduler    IScheduler
	firstInMs    int64
	intervalInMs int64
	task         Task
	disposed     int32
	lock         sync.Mutex
}

func newTimerTask(scheduler IScheduler, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	var t = &timerTask{scheduler: scheduler, task: task, firstInMs: firstInMs, intervalInMs: intervalInMs}
	return t
}

// Dispose release resources
func (t *timerTask) Dispose() {
	if atomic.CompareAndSwapInt32(&t.disposed, 0, 1) {
		return
	}
	t.scheduler.Remove(t)
}

func (t *timerTask) schedule() {
	if t.firstInMs <= 0 {
		t.doFirstSchedule()
		return
	}

	first := time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
	go func() {
		select {
		case <-first.C:
			t.doFirstSchedule()
		}
	}()
}

func (t *timerTask) doFirstSchedule() {
	t.executeOnFiber()
	if t.intervalInMs <= 0 {
		t.Dispose()
		return
	}

	interval := time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	go func() {
		for atomic.LoadInt32(&t.disposed) == 0 {
			/*select {
			  case <-t.interval.C:
			  	t.executeOnFiber()
			  }*/
			<-interval.C
			t.executeOnFiber()
		}
		interval.Stop()
	}()
}

func (t *timerTask) executeOnFiber() {
	if atomic.LoadInt32(&t.disposed) == 1 {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}
