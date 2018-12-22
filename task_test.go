package robin

import (
	"reflect"
	"testing"
	"time"
)

func Test_newTask(t *testing.T) {
	type args struct {
		t interface{}
		p []interface{}
	}
	tests := []struct {
		name string
		args args
		want Task
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTask(tt.args.t, tt.args.p...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			}
		})
	}
}

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
				executedTime := timeDuration.Nanoseconds() / 10000
				if executedTime < ttt.want {
					t.Logf("executed time error %v", timeDuration/time.Nanosecond)
				}
			}
		})
	}
}

func Test_newTimerTask(t *testing.T) {
	type args struct {
		fiber        SchedulerRegistry
		task         Task
		firstInMs    int64
		intervalInMs int64
	}
	tests := []struct {
		name string
		args args
		want *timerTask
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTimerTask(tt.args.fiber, tt.args.task, tt.args.firstInMs, tt.args.intervalInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTimerTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timerTask_init(t *testing.T) {
	type fields struct {
		identifyID   string
		scheduler    SchedulerRegistry
		firstInMs    int64
		intervalInMs int64
		first        *time.Timer
		interval     *time.Ticker
		task         Task
		cancelled    bool
	}
	type args struct {
		scheduler    SchedulerRegistry
		task         Task
		firstInMs    int64
		intervalInMs int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *timerTask
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &timerTask{
				identifyID:   tt.fields.identifyID,
				scheduler:    tt.fields.scheduler,
				firstInMs:    tt.fields.firstInMs,
				intervalInMs: tt.fields.intervalInMs,
				first:        tt.fields.first,
				interval:     tt.fields.interval,
				task:         tt.fields.task,
				cancelled:    tt.fields.cancelled,
			}
			if got := ta.init(tt.args.scheduler, tt.args.task, tt.args.firstInMs, tt.args.intervalInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timerTask.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"Test_timerTask_schedule", newTimerTask(g.scheduler.(*Scheduler), newTask(func(s int64) {
			runCount++
			t.Logf("count:%v s:%v", runCount, s)
		}, time.Now().UnixNano()), 10, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.schedule()
			for true {
				if runCount >= 10 {
					tt.fields.Dispose()
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
