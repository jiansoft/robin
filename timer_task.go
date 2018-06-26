package robin

import (
	"fmt"
	"time"
)

type timerTask struct {
	identifyId   string
	scheduler    SchedulerRegistry
	firstInMs    int64
	intervalInMs int64
	firstTime    *time.Timer
	interval     *time.Ticker
	task         task
	cancelled    bool
}

func (t *timerTask) init(scheduler SchedulerRegistry, task task, firstInMs int64, intervalInMs int64) *timerTask {
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	return t
}

func newTimerTask(fiber SchedulerRegistry, task task, firstInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) Dispose() {
	t.cancelled = true
	if nil != t.firstTime {
		t.firstTime.Stop()
		t.firstTime = nil
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
	t.scheduler.Enqueue(t.runFirstSchedule)
}

func (t *timerTask) runFirstSchedule() {
	t.firstTime = time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
	select {
	case <-t.firstTime.C:
		if t.cancelled {
			return
		}

		t.scheduler.Enqueue(t.execute)

		if nil != t.firstTime {
			t.firstTime.Stop()
			t.firstTime = nil
		}

		if t.intervalInMs <= 0 {
			t.Dispose()
			return
		}

		t.scheduler.Enqueue(t.runIntervalSchedule)
	}
}

func (t *timerTask) runIntervalSchedule() {
	t.interval = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	for !t.cancelled {
		//select {
		//case <-t.interval.C:
		//	t.scheduler.Enqueue(t.execute)
		//}
		<-t.interval.C
		t.scheduler.Enqueue(t.execute)
	}
}

func (t timerTask) execute() {
	if t.cancelled {
		return
	}
	t.task.run()
}
