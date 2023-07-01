package robin

import (
	"sync"
	"testing"
)

func Test_defaultExecutor_ExecuteTasks(t *testing.T) {
	var wg sync.WaitGroup
	type args struct {
		tasks []task
	}
	tests := []struct {
		name string
		d    defaultExecutor
		args args
	}{
		{"TestExecuteTasks", defaultExecutor{}, args{tasks: []task{newTask(func(s string) {
			//t.Logf("s:%v", s)
			wg.Done()
		}, "ExecuteTasks")}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)
			tt.d.executeTasks(tt.args.tasks)
			wg.Add(1)
			tt.d.executeTasksWithGoroutine(tt.args.tasks)
			wg.Add(1)
			for _, value := range tt.args.tasks {
				tt.d.executeTaskWithGoroutine(value)
			}
			wg.Wait()
		})
	}
}
