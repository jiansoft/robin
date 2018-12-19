package robin

import (
	"reflect"
	"testing"
	"time"
)

func Test_newDefaultExecutor(t *testing.T) {
	tests := []struct {
		name string
		want defaultExecutor
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDefaultExecutor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDefaultExecutor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultExecutor_ExecuteTasks(t *testing.T) {
	type args struct {
		tasks []Task
	}
	tests := []struct {
		name string
		d    defaultExecutor
		args args
	}{
		{"TestExecuteTasks", defaultExecutor{}, args{tasks: []Task{newTask(func(s string) { t.Logf("s:%v", s) }, "ExecuteTasks")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.ExecuteTasks(tt.args.tasks)
		})
	}
}

func Test_defaultExecutor_ExecuteTasksWithGoroutine(t *testing.T) {
	type args struct {
		tasks []Task
	}
	tests := []struct {
		name string
		d    defaultExecutor
		args args
	}{
		{"TestExecuteTasksWithGoroutine", defaultExecutor{}, args{tasks: []Task{newTask(func(s string) { t.Logf("s:%v", s) }, "ExecuteTasksWithGoroutine")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.ExecuteTasksWithGoroutine(tt.args.tasks)
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}
