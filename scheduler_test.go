package robin

import (
	"sync"
	"sync/atomic"
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
		{"Test_Scheduler_ScheduleOnInterval", fields{gs: gs, gm: gm}},
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
	wg := sync.WaitGroup{}
	wg.Add(4)
	loop := int32(0)
	s := newScheduler(fiber)

	taskFun := func(s string, t *testing.T) {
		atomic.AddInt32(&loop, 1)
		//t.Logf("loop:%v msg:%v",atomic.LoadInt32(&loop), s)
		wg.Done()
	}

	s.Enqueue(taskFun, "Enqueue", t)
	s.EnqueueWithTask(newTask(taskFun, "EnqueueWithTask", t))
	s.Schedule(0, taskFun, "Schedule", t)
	interval := s.ScheduleOnInterval(0, 100, taskFun, "first", t)
	wg.Wait()
	interval.Dispose()

	if int32(4) != atomic.LoadInt32(&loop) {
		t.Fatal("they should be equal")
	}

	wg.Add(4)
	s.ScheduleOnInterval(0, 10, taskFun, "second", t)
	s.isDispose = true
	d := s.ScheduleOnInterval(0, 10, taskFun, "remove", t)
	d.Dispose()
	//s.Remove(d)

	wg.Wait()
	s.Dispose()
}
