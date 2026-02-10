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

	type fields struct {
		gs Fiber
		gm Fiber
	}

	tests := []struct {
		fields fields
		name   string
	}{
		{fields{gs: gs, gm: gm}, "Test_Scheduler_ScheduleOnInterval"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedulerTest(t, tt.fields.gs)
			schedulerTest(t, tt.fields.gm)
			timeout := time.NewTimer(10 * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestSchedulerRemove(t *testing.T) {
	gm := NewGoroutineMulti()
	defer func() {
		gm.Dispose()
		time.Sleep(50 * time.Millisecond)
	}()

	s := newScheduler(gm)
	d := s.ScheduleOnInterval(0, 100, func() {})

	// Verify the timerTask is stored in the sync.Map
	found := false
	s.tasks.Range(func(k, v any) bool {
		if k == d {
			found = true
			return false
		}
		return true
	})
	if !found {
		t.Fatal("timerTask should be stored in scheduler before Remove")
	}

	// Remove it
	s.Remove(d)

	// Verify it's gone
	found = false
	s.tasks.Range(func(k, v any) bool {
		if k == d {
			found = true
			return false
		}
		return true
	})
	if found {
		t.Fatal("timerTask should be removed from scheduler after Remove")
	}

	d.Dispose()
}

func schedulerTest(t *testing.T, f Fiber) {
	wg := sync.WaitGroup{}
	wg.Add(4)
	loop := int32(0)
	s := newScheduler(f)

	taskFun := func(s string, t *testing.T) {
		atomic.AddInt32(&loop, 1)
		wg.Done()
	}

	s.Enqueue(taskFun, "Enqueue", t)
	s.enqueueTask(newTask(taskFun, "enqueueTask", t))
	s.Schedule(0, taskFun, "Schedule", t)
	interval := s.ScheduleOnInterval(0, 100, taskFun, "first", t)
	wg.Wait()
	interval.Dispose()

	if int32(4) != atomic.LoadInt32(&loop) {
		t.Fatal("they should be equal")
	}

	wg.Add(4)
	s.ScheduleOnInterval(0, 10, taskFun, "second", t)
	wg.Wait()

	// Dispose stops all running tasks and sets disposed=true.
	s.Dispose()

	// ScheduleOnInterval after Dispose returns without scheduling.
	d := s.ScheduleOnInterval(0, 10, taskFun, "remove", t)
	d.Dispose()
}
