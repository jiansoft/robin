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
	firstTimer   *time.Timer
	timer        *time.Ticker
	task         Task
	cancelled    bool
}

func (t *timerTask) init(scheduler SchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	return t
}

func newTimerTask(fiber SchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) Dispose() {
	t.cancelled = true
	if nil != t.firstTimer {
		t.firstTimer.Stop()
		t.firstTimer = nil
	}
	if nil != t.timer {
		t.timer.Stop()
		t.timer = nil
	}
	if nil != t.scheduler {
		t.scheduler.Remove(t)
		t.scheduler = nil
	}
}

func (t *timerTask) Identify() string {
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
    t.timer = time.NewTicker(time.Duration(timeInMs) * time.Millisecond)
    go func() {
        for !t.cancelled {
            select {
            case <-t.timer.C:
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
                t.timer = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
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
		t.firstTimer = time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
		go func() {
			select {
			case <-t.firstTimer.C:
				t.firstTimer.Stop()
				if t.cancelled {
					break
				}
				t.doFirstSchedule()
			}
		}()
	}
}

func (t *timerTask) doFirstSchedule() {
	t.scheduler.Enqueue(t.executeOnFiber)
	t.doIntervalSchedule()
}

func (t *timerTask) doIntervalSchedule() {
	if t.intervalInMs <= 0 {
		t.Dispose()
		return
	}
	t.timer = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	go func() {
		for !t.cancelled {
			select {
			case <-t.timer.C:
				if t.cancelled {
					break
				}
				t.scheduler.Enqueue(t.executeOnFiber)
			}
		}
	}()
}

func (t timerTask) executeOnFiber() {
	if t.cancelled {
		return
	}
	t.task.Run()
}
