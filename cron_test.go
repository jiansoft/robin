package robin

import (
	"reflect"
	"testing"
	"time"
)

func TestRightNow(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RightNow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RightNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelay(t *testing.T) {
	type args struct {
		delayInMs int64
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delay(tt.args.delayInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCronDelay(t *testing.T) {
	tests := []struct {
		name string
		want *cronDelay
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCronDelay(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCronDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDelayJob(t *testing.T) {
	type args struct {
		delayInMs int64
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDelayJob(tt.args.delayInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDelayJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cronDelay_init(t *testing.T) {
	type fields struct {
		fiber Fiber
	}
	tests := []struct {
		name   string
		fields fields
		want   *cronDelay
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cronDelay{
				fiber: tt.fields.fiber,
			}
			if got := c.init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cronDelay.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cronDelay_Delay(t *testing.T) {
	type fields struct {
		fiber Fiber
	}
	type args struct {
		delayInMs int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cronDelay{
				fiber: tt.fields.fiber,
			}
			if got := c.Delay(tt.args.delayInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cronDelay.Delay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEveryCron(t *testing.T) {
	tests := []struct {
		name string
		want *cronEvery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEveryCron(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEveryCron() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cronEvery_init(t *testing.T) {
	type fields struct {
		fiber Fiber
	}
	tests := []struct {
		name   string
		fields fields
		want   *cronEvery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cronEvery{
				fiber: tt.fields.fiber,
			}
			if got := c.init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cronEvery.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEverySunday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EverySunday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EverySunday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryMonday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryMonday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EveryMonday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryTuesday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryTuesday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EveryTuesday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryWednesday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryWednesday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EveryWednesday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryThursday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryThursday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EveryThursday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryFriday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryFriday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EveryFriday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEverySaturday(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EverySaturday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EverySaturday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvery(t *testing.T) {
	type args struct {
		interval int64
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Every(tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Every() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cronEvery_Every(t *testing.T) {
	type fields struct {
		fiber Fiber
	}
	type args struct {
		interval int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cronEvery{
				fiber: tt.fields.fiber,
			}
			if got := c.Every(tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cronEvery.Every() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newEveryJob(t *testing.T) {
	type args struct {
		weekday time.Weekday
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEveryJob(tt.args.weekday); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEveryJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJob(t *testing.T) {
	type args struct {
		intervel  int64
		fiber     Fiber
		delayUnit delayUnit
	}
	tests := []struct {
		name string
		args args
		want *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJob(tt.args.intervel, tt.args.fiber, tt.args.delayUnit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_init(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		intervel  int64
		fiber     Fiber
		delayUnit delayUnit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.init(tt.args.intervel, tt.args.fiber, tt.args.delayUnit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Dispose(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.Dispose()
		})
	}
}

func TestJob_Identify(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Identify(); got != tt.want {
				t.Errorf("Job.Identify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Days(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Days(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Days() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Hours(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Hours(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Hours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Minutes(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Minutes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Minutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Seconds(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Seconds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Seconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_MilliSeconds(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.MilliSeconds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.MilliSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_At(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		hour   int
		minute int
		second int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.At(tt.args.hour, tt.args.minute, tt.args.second); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.At() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_AfterExecuteTask(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.AfterExecuteTask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.AfterExecuteTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_BeforeExecuteTask(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
		want   *Job
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.BeforeExecuteTask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.BeforeExecuteTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_firstTimeSetDelayNextTime(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.firstTimeSetDelayNextTime(tt.args.now)
		})
	}
}

func TestJob_firstTimeSetWeeksNextTime(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.firstTimeSetWeeksNextTime(tt.args.now)
		})
	}
}

func TestJob_firstTimeSetDaysNextTime(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.firstTimeSetDaysNextTime(tt.args.now)
		})
	}
}

func TestJob_firstTimeSetHoursNextTime(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.firstTimeSetHoursNextTime(tt.args.now)
		})
	}
}

func TestJob_firstTimeSetMinutesNextTime(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.firstTimeSetMinutesNextTime(tt.args.now)
		})
	}
}

func TestJob_Do(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	type args struct {
		fun    interface{}
		params []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Disposable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			if got := c.Do(tt.args.fun, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_canDo(t *testing.T) {
	type fields struct {
		fiber        Fiber
		identifyId   string
		loc          *time.Location
		task         Task
		taskDisposer Disposable
		weekday      time.Weekday
		hour         int
		minute       int
		second       int
		unit         unit
		delayUnit    delayUnit
		interval     int64
		nextTime     time.Time
		timingMode   timingAfterOrBeforeExecuteTask
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Job{
				fiber:        tt.fields.fiber,
				identifyId:   tt.fields.identifyId,
				loc:          tt.fields.loc,
				task:         tt.fields.task,
				taskDisposer: tt.fields.taskDisposer,
				weekday:      tt.fields.weekday,
				hour:         tt.fields.hour,
				minute:       tt.fields.minute,
				second:       tt.fields.second,
				unit:         tt.fields.unit,
				delayUnit:    tt.fields.delayUnit,
				interval:     tt.fields.interval,
				nextTime:     tt.fields.nextTime,
				timingMode:   tt.fields.timingMode,
			}
			c.canDo()
		})
	}
}
