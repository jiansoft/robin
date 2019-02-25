package robin

import (
	"reflect"
	"sync"
	"time"
)

//Task a struct
type Task struct {
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func newTask(f interface{}, p ...interface{}) Task {
	var paramLen = len(p)
	task := Task{funcCache: reflect.ValueOf(f), paramsCache: make([]reflect.Value, paramLen)}
	if paramLen > 0 {
		for k, param := range p {
			task.paramsCache[k] = reflect.ValueOf(param)
		}
	}
	return task
}

func (t Task) run() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}

func (t *Task) Params(p ...interface{}) {
	t.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		t.paramsCache[k] = reflect.ValueOf(param)
	}
}

type timerTask struct {
	scheduler    SchedulerRegistry
	firstInMs    int64
	intervalInMs int64
	task         Task
	disposed     bool
	lock         sync.Mutex
}

func newTimerTask(scheduler SchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	t := &timerTask{}
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	return t
}

// Dispose release resources
func (t *timerTask) Dispose() {
	if t.getDisposed() {
		return
	}
	t.setDisposed(true)
	t.scheduler.Remove(t)
}

// Identify return the struct identify id
func (t *timerTask) Identify() string {
	return "" // t.identifyID
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
	t.doIntervalSchedule()
}

func (t *timerTask) doIntervalSchedule() {
	if t.intervalInMs <= 0 {
		t.Dispose()
		return
	}
	interval := time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	go func() {
		for !t.getDisposed() {
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
	if t.getDisposed() {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}

func (t *timerTask) getDisposed() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.disposed
}

func (t *timerTask) setDisposed(r bool) {
	t.lock.Lock()
	t.disposed = r
	t.lock.Unlock()
}
