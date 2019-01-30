package robin

import (
	"testing"
	"time"
)

func TestScheduler_ScheduleOnInterval(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	type fields struct {
		fiber       executionContext
		running     bool
		isDispose   bool
		disposabler *disposer
	}
	type args struct {
		firstInMs   int64
		regularInMs int64
		taskFun     interface{}
		params      []interface{}
	}
	strs := []string{"first"}
	names := make([]interface{}, len(strs))
	for i, s := range strs {
		names[i] = s
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test_Scheduler_ScheduleOnInterval",
			fields{fiber: g, running: true, isDispose: false, disposabler: NewDisposer()},
			args{firstInMs: 0, regularInMs: 10, taskFun: func(s string) { t.Logf("s:%v", s) }, params: names}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				fiber:       tt.fields.fiber,
				running:     tt.fields.running,
				isDispose:   tt.fields.isDispose,
				disposabler: tt.fields.disposabler,
			}
			s.isDispose = true
			s.ScheduleOnInterval(tt.args.firstInMs, tt.args.regularInMs, tt.args.taskFun, tt.args.params...)
			s.isDispose = false
			s.ScheduleOnInterval(tt.args.firstInMs, tt.args.regularInMs, tt.args.taskFun, tt.args.params...)
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestScheduler_Enqueue(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	type fields struct {
		fiber       executionContext
		running     bool
		isDispose   bool
		disposabler *disposer
	}
	type args struct {
		taskFun interface{}
		params  []interface{}
	}
	strs := []string{"first"}
	names := make([]interface{}, len(strs))
	for i, s := range strs {
		names[i] = s
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test_Scheduler_Enqueue",
			fields{fiber: g, running: true, isDispose: false, disposabler: NewDisposer()},
			args{taskFun: func(s string) { t.Logf("s:%v", s) }, params: names}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				fiber:       tt.fields.fiber,
				running:     tt.fields.running,
				isDispose:   tt.fields.isDispose,
				disposabler: tt.fields.disposabler,
			}
			s.Enqueue(tt.args.taskFun, tt.args.params...)
		})
	}
}

func TestScheduler_Stop(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	type fields struct {
		fiber       executionContext
		running     bool
		isDispose   bool
		disposabler *disposer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"Test_Scheduler_Stop", fields{fiber: g, running: true, isDispose: false, disposabler: NewDisposer()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				fiber:       tt.fields.fiber,
				running:     tt.fields.running,
				isDispose:   tt.fields.isDispose,
				disposabler: tt.fields.disposabler,
			}
			s.Stop()
		})
	}
}

func TestScheduler_Dispose(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	type fields struct {
		fiber       executionContext
		running     bool
		isDispose   bool
		disposabler *disposer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"Test_Scheduler_Dispose", fields{fiber: g, running: true, isDispose: false, disposabler: NewDisposer()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				fiber:       tt.fields.fiber,
				running:     tt.fields.running,
				isDispose:   tt.fields.isDispose,
				disposabler: tt.fields.disposabler,
			}
			s.Dispose()
		})
	}
}
