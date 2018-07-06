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

type pendingTask struct {
	identifyId string
	task       task
	cancelled  bool
}

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

func newPendingTask(task task) *pendingTask {
	return new(pendingTask).init(task)
}

func (p *pendingTask) init(task task) *pendingTask {
	p.task = task
	p.cancelled = false
	p.identifyId = fmt.Sprintf("%p-%p", &p, &task)
	return p
}

func (p *pendingTask) Dispose() {
	p.cancelled = true
}

func (p *pendingTask) Identify() string {
	return p.identifyId
}

func (p pendingTask) execute() {
	if p.cancelled {
		return
	}
	p.task.run()
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
	if nil != t.first {
		t.first.Stop()
		t.first = nil
	}

	if nil != t.interval {
		t.interval.Stop()
		t.interval = nil
	}

	if nil != t.scheduler {
		t.scheduler.Remove(t)
		t.scheduler = nil
	}
}

func (t timerTask) Identify() string {
	return t.identifyId
}

/*
func (t *timerTask) schedule() {
    timeInMs := t.firstInMs
    if timeInMs <= 0 {
        t.scheduler.EnqueueWithTask(t.task)
        if t.intervalInMs <= 0 {
            t.Dispose()
            return
        }
        timeInMs = t.intervalInMs
        t.firstRan = true
    }
    t.interval = time.NewTicker(time.Duration(timeInMs) * time.Millisecond)
    go func() {
        for !t.cancelled {
            select {
            case <-t.interval.C:
                if t.cancelled {
                    break
                }
                t.scheduler.EnqueueWithTask(t.task)
                if t.intervalInMs <= 0 {
                    t.Dispose()
                    break
                }
                if t.firstRan {
                    continue
                }
                t.interval = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
                t.firstRan = true
            }
        }
    }()
}
*/

func (t *timerTask) schedule() {
    if t.firstInMs <= 0 {
        t.doFirstSchedule()
    } else {
        t.first = time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
        go func() {
            select {
            case <-t.first.C:
                t.first.Stop()
                if t.cancelled {
                    break
                }
                t.doFirstSchedule()
            }
        }()
    }
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
            select {
            case <-t.interval.C:
                t.executeOnFiber()
            }
        }
    }()
}

func (t timerTask) executeOnFiber() {
    if t.cancelled {
        return
    }
    t.scheduler.EnqueueWithTask(t.task)
}

