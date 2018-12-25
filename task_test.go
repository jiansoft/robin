package robin

import (
	"testing"
	"time"
)

func TestTask_run(t *testing.T) {
	type args struct {
		task Task
	}
	params := []args{
		{newTask(func(s string) { t.Logf("s:%v", s) }, "run 1")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "run 2")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "run 3")}}
	tests := []struct {
		name string
		args []args
	}{
		{"Test_Task_run", params},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, ttt := range tt.args {
				ttt.task.run()
				ttt.task.release()
			}
		})
	}
}

/*
func TestTask_Run(t *testing.T) {
	type args struct {
		task Task
		want int64
	}
	params := []args{
		{newTask(func(s string) {
			t.Logf("s:%v", s)
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		}, "run 1"), 100},
		{newTask(func(s string) {
			t.Logf("s:%v", s)
			timeout := time.NewTimer(time.Duration(200) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		}, "run 2"), 200},
		{newTask(func(s string) {
			t.Logf("s:%v", s)
			timeout := time.NewTimer(time.Duration(300) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		}, "run 3"), 300}}
	tests := []struct {
		name string
		args []args
	}{
		{"Test_Task_Run", params},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, ttt := range tt.args {
				timeDuration := ttt.task.Run()
				ttt.task.release()
				executedTime := timeDuration.Nanoseconds() / 10000
				if executedTime < ttt.want {
					t.Logf("executed time error %v", timeDuration/time.Nanosecond)
				}
			}
		})
	}
}*/

func Test_timerTask_Dispose(t *testing.T) {
	runCount := 0
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_Dispose", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 10, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.schedule()
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			tt.fields.Dispose()
		})
	}
}

func Test_timerTask_Identify(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_Identify", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) { t.Logf("s:%v", s) }, time.Now().UnixNano()), 10, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Identify(); got == "" {
				t.Errorf("timerTask.Identify() = %v", got)
			}
			t.Logf("timerTask.Identify() = %v", tt.fields.Identify())
		})
	}
}

func Test_timerTask_schedule(t *testing.T) {
	runCount := 0
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_schedule_1", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("schedule_1 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 0, 5)},
		{"Test_timerTask_schedule_2", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("schedule_2 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 5, 10)},
		{"Test_timerTask_schedule_3", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("schedule_3 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 10, 20)},
		{"Test_timerTask_schedule_4", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("schedule_4 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 20, 20)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Test_timerTask_schedule_4" {
				tt.fields.cancelled = true
			}

			tt.fields.schedule()
			for true {
				if tt.name == "Test_timerTask_schedule_1" && runCount >= 10 {
					tt.fields.Dispose()
					return
				}
				if tt.name == "Test_timerTask_schedule_2" && runCount >= 20 {
					tt.fields.Dispose()
					return
				}
				if tt.name == "Test_timerTask_schedule_3" && runCount >= 31 {
					return
				}
				if tt.name == "Test_timerTask_schedule_3" && runCount >= 30 {
					tt.fields.cancelled = true
					runCount++
				}
				if tt.name == "Test_timerTask_schedule_4" {
					timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
					select {
					case <-timeout.C:
					}
					return
				}
			}
		})
	}
}

func Test_timerTask_doFirstSchedule(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_doFirstSchedule", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			t.Logf("s:%v", s)
		}, time.Now().UnixNano()), 10, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.doFirstSchedule()
			timeout := time.NewTimer(time.Duration(10) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func Test_timerTask_doIntervalSchedule(t *testing.T) {
	runCount := 0
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_doIntervalSchedule_1", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("Schedule_1 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 0, 0)},
		{"Test_timerTask_doIntervalSchedule_2", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("Schedule_2 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 0, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.doIntervalSchedule()
			timeout := time.NewTimer(time.Duration(101) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func Test_timerTask_executeOnFiber(t *testing.T) {
	runCount := 0
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_executeOnFiber_1", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("executeOnFiber_1 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 0, 0)},
		{"Test_timerTask_executeOnFiber_2", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("executeOnFiber_2 count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 0, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.executeOnFiber()
			timeout := time.NewTimer(time.Duration(101) * time.Millisecond)
			select {
			case <-timeout.C:
			}
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
		{"Test_Schedule_Mix_2"},
		{"Test_Schedule_Mix_3"},
		{"Test_Schedule_Mix_4"},
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

			g.Schedule(40, func(s int64) {
				t.Logf("run_4 s:%v", s)
			}, int64(40))

			g.Schedule(70, func(s int64) {
				t.Logf("run_7 s:%v", s)
			}, int64(70))

			timeout := time.NewTimer(time.Duration(150) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}
