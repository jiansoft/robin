package robin

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_timerTask_schedule(t *testing.T) {
	var runCount int32
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_schedule_1", newTimerTask(g.scheduler.(*scheduler), newTask(func(s int64) {
			atomic.AddInt32(&runCount, 1)
		}, time.Now().UnixNano()), 0, 5)},
		/*{"Test_timerTask_schedule_2", newTimerTask(g.scheduler.(*scheduler), newTask(func(s int64) {
		    atomic.AddInt32(&runCount, 1)
		}, time.Now().UnixNano()), 5, 10)},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Test_timerTask_schedule_4" {
				tt.fields.disposed = 1
			}
			tt.fields.schedule()
			for true {
				saveRunCount := atomic.LoadInt32(&runCount)
				if tt.name == "Test_timerTask_schedule_1" && saveRunCount >= 10 {
					tt.fields.Dispose()
					return
				}
				if tt.name == "Test_timerTask_schedule_2" && saveRunCount >= 20 {
					tt.fields.Dispose()
					return
				}
				if tt.name == "Test_timerTask_schedule_3" && saveRunCount >= 31 {
					return
				}
				if tt.name == "Test_timerTask_schedule_3" && saveRunCount >= 30 {
					atomic.SwapInt32(&tt.fields.disposed, 1)
					atomic.AddInt32(&runCount, 1)
				}
				if tt.name == "Test_timerTask_schedule_4" {
					timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
					select {
					case <-timeout.C:
					}
					return
				}
			}
			tt.fields.Dispose()
		})
	}
}

func Test_Schedule_Mix(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name string
	}{
		{"Test_Schedule_Mix_1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := g.Schedule(1000, func(s int64) {
				t.Logf("run_1 s:%v", s)
			}, int64(1000))
			d.Dispose()

			d = g.Schedule(20, func(s int64) {
				t.Logf("run_2 s:%v", s)
				g.Schedule(30, func(s int64) {
					t.Logf("run_3 s:%v", s)
					g.Schedule(60, func(s int64) {
						t.Logf("run_6 s:%v", s)
					}, int64(60+30+20))
				}, int64(30+20))
			}, int64(20))

			g40 := g.Schedule(40, func(s int64) {
				t.Logf("run_4 s:%v", s)
			}, int64(40))

			g70 := g.Schedule(70, func(s int64) {
				t.Logf("run_7 s:%v", s)
			}, int64(70))

			timeout := time.NewTimer(time.Duration(150) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			d.Dispose()
			g40.Dispose()
			g70.Dispose()
		})
	}
	timeout := time.NewTimer(time.Duration(150) * time.Millisecond)
	select {
	case <-timeout.C:
	}
}
