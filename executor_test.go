package robin

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestDefaultExecutorExecuteTask verifies executeTask runs exactly once in synchronous mode.
// TestDefaultExecutorExecuteTask 驗證 executeTask 在同步模式下只會執行一次。
func TestDefaultExecutorExecuteTask(t *testing.T) {
	var called int32
	de := newDefaultExecutor()
	tk := newTask(func() {
		atomic.AddInt32(&called, 1)
	})

	de.executeTask(tk)

	if got := atomic.LoadInt32(&called); got != 1 {
		t.Errorf("executeTask() called %d times, want 1", got)
	}
}

// Test_defaultExecutor_ExecuteTasks verifies sequential and goroutine execution paths both invoke tasks.
// Test_defaultExecutor_ExecuteTasks 驗證序列與 goroutine 路徑都會正確觸發任務。
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
