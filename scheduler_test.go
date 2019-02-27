package robin

import (
    "testing"
    "time"
)

func TestScheduler(t *testing.T) {
    gs := NewGoroutineSingle()
    gm := NewGoroutineMulti()
    gs.Start()
    gm.Start()

    type fields struct {
        gs Fiber
        gm Fiber
    }

    tests := []struct {
        name   string
        fields fields
    }{
        {"Test_Scheduler_ScheduleOnInterval", fields{gs: gs, gm: gm},},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            schedulerTest(t, tt.fields.gs)
            schedulerTest(t, tt.fields.gm)
            timeout := time.NewTimer(time.Duration(10) * time.Millisecond)
            select {
            case <-timeout.C:
            }
        })
    }
}

func schedulerTest(t *testing.T, fiber Fiber) {
    s := newScheduler(fiber)
    taskFun := func(s string, t *testing.T) {
        t.Logf("msg:%v", s)
    }
    s.Enqueue(taskFun, "Enqueue", t)
    s.EnqueueWithTask(newTask(taskFun, "EnqueueWithTask", t))
    s.Schedule(0, taskFun, "Schedule", t)
    s.ScheduleOnInterval(0, 10, taskFun, "first", t)
    <-time.After(time.Duration(10) * time.Millisecond)
    s.ScheduleOnInterval(0, 10, taskFun, "second", t)
    <-time.After(time.Duration(20) * time.Millisecond)
    s.isDispose = true
    d := s.ScheduleOnInterval(0, 10, taskFun, "remove", t)
    d.Dispose()
    s.Remove(d)
    s.Dispose()
}
