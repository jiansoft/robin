package robin

import (
	"fmt"
	"reflect"
	"time"
)

type task struct {
	doFunc      interface{}
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func newTask(t interface{}, p ...interface{}) task {
	task := task{doFunc: t}
	task.funcCache = reflect.ValueOf(t)
	task.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t task) run() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}

func (t task) Run() time.Duration {
	s := time.Now()
	t.run()
	e := time.Now()
	return e.Sub(s)
}

/*func (t *task) Dispose() {
	t.doFunc = nil
	t.paramsCache = t.paramsCache[:0]
	t.paramsCache = nil
}*/

type timerTask struct {
	identifyId   string
	scheduler    SchedulerRegistry
	firstInMs    int64
	intervalInMs int64
	first        *time.Timer
	interval     *time.Ticker
	task         task
	cancelled    bool
}

func newTimerTask(fiber SchedulerRegistry, task task, firstInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) init(scheduler SchedulerRegistry, task task, firstInMs int64, intervalInMs int64) *timerTask {
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	return t
}

func (t *timerTask) Dispose() {
	t.cancelled = true

	if nil != t.scheduler {
		t.scheduler.Remove(t)
		t.scheduler = nil
	}
}

func (t timerTask) Identify() string {
	return t.identifyId
}

func (t *timerTask) schedule() {
	if t.firstInMs <= 0 {
		t.doFirstSchedule()
		return
	}

	t.first = time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
	go func() {
		select {
		case <-t.first.C:
			if t.cancelled {
				break
			}
			t.doFirstSchedule()
		}
		t.first = nil
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
	t.interval = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	go func() {
		for !t.cancelled {
			/*select {
			case <-t.interval.C:
				t.executeOnFiber()
			}*/
			<-t.interval.C
			t.executeOnFiber()
		}
		t.interval.Stop()
		t.interval = nil
	}()
}

func (t timerTask) executeOnFiber() {
	if t.cancelled {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}
