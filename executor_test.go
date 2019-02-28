package robin

import (
	"sync"
	"testing"
)

func Test_defaultExecutor_ExecuteTasks(t *testing.T) {
	var wg sync.WaitGroup
	type args struct {
		tasks []Task
	}
	tests := []struct {
		name string
		d    defaultExecutor
		args args
	}{
		{"TestExecuteTasks", defaultExecutor{}, args{tasks: []Task{newTask(func(s string) {
			//t.Logf("s:%v", s)
			wg.Done()
		}, "ExecuteTasks")}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)
			tt.d.ExecuteTasks(tt.args.tasks)
			wg.Add(1)
			tt.d.ExecuteTasksWithGoroutine(tt.args.tasks)
			wg.Add(1)
			for _, value := range tt.args.tasks {
				tt.d.ExecuteTaskWithGoroutine(value)
			}
			wg.Wait()
		})
	}
}
