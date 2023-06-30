package robin

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_timerTask_schedule(t *testing.T) {
	g := NewGoroutineMulti()
	tests := []struct {
		timerTask *timerTask
		name      string
	}{
		{newTimerTask(g.scheduler.(*scheduler), newTask(func() {
			fmt.Printf("go... %v\n", time.Now())
		}), -100, 1000), "Test_timerTask_schedule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.timerTask.schedule()
			<-time.After(time.Duration(3) * time.Second)
			tt.timerTask.Dispose()
			t1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {
				fmt.Printf("go... %v\n", time.Now())
			}), 1000, 1000)
			t1.Dispose()
			<-time.After(time.Duration(1) * time.Second)

		})
	}
}

func Test_timerTask_schedule_runFirst(t *testing.T) {
	g := NewGoroutineMulti()
	tests := []struct {
		timerTask *timerTask
		name      string
	}{
		{newTimerTask(g.scheduler.(*scheduler), newTask(func() {
			fmt.Printf("go... %v\n", time.Now())
		}), 100000, 10000), "Test_timerTask_schedule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.timerTask.schedule()
			<-time.After(time.Duration(1) * time.Second)
			tt.timerTask.Dispose()
			<-time.After(time.Duration(1) * time.Second)

		})
	}
}

func Test_timerTask(t *testing.T) {
	var runCount int32
	g := NewGoroutineMulti()
	tests := []struct {
		fields *timerTask
		name   string
	}{
		{newTimerTask(g.scheduler.(*scheduler), newTask(func() {
			atomic.AddInt32(&runCount, 1)
		}), 0, 5), "Test_timerTask_schedule_1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.schedule()
			for true {
				if atomic.LoadInt32(&runCount) >= 10 {
					tt.fields.Dispose()
					break
				}
			}

			tt.fields.executeOnFiber()
			if int32(10) != atomic.LoadInt32(&runCount) {
				t.Fatal("they should be equal")
			}

			wg := sync.WaitGroup{}
			wg.Add(2)
			var runT1Count int32
			t1 := g.ScheduleOnInterval(0, 10, func() {
				atomic.AddInt32(&runT1Count, 1)
				wg.Done()
			})
			wg.Wait()
			t1.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			if int32(2) != atomic.LoadInt32(&runT1Count) {
				t.Fatal("they should be equal")
			}

			var runT2Count int32
			t2 := g.ScheduleOnInterval(1000, 10, func() {
				atomic.AddInt32(&runT2Count, 1)
			})
			t2.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			if int32(0) != atomic.LoadInt32(&runT2Count) {
				t.Fatal("they should be equal")
			}

			var runT3Count int32
			t3 := g.Schedule(1000, func() {
				atomic.AddInt32(&runT3Count, 1)
			})
			t3.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			if int32(0) != atomic.LoadInt32(&runT3Count) {
				t.Fatal("they should be equal")
			}

			var runT4Count int32
			t4 := g.Schedule(0, func() {
				atomic.AddInt32(&runT4Count, 1)
			})
			t4.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			if int32(1) != atomic.LoadInt32(&runT4Count) {
				t.Fatal("they should be equal")
			}
		})
	}
}

func Test_T(t *testing.T) {

	tests := []struct {
		name string
		task Task
	}{
		{task: newTask(exampleFunc1, "QQQQ"), name: "Test_timerTask_schedule"},
		{task: newTask(exampleFunc2), name: "Test_timerTask_schedule"},
		//		{task: newTask(exampleFunc2, "AAAA"), name: "Test_timerTask_schedule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.task.execute()
			//tt.timerTask.schedule()
			//<-time.After(time.Duration(6) * time.Second)
			//tt.timerTask.Dispose()
			//<-time.After(time.Duration(10) * time.Second)

		})
	}
}

func exampleFunc1(args ...string) {
	fmt.Printf("exampleFunc1 %v", args)
}
func exampleFunc2() {
	fmt.Printf("exampleFunc2 ")
}
