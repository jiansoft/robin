package robin

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Test_timerTask_schedule validates immediate schedule path and disposal behavior for timerTask.
// Test_timerTask_schedule 驗證 timerTask 的立即排程路徑與 Dispose 行為。
func Test_timerTask_schedule(t *testing.T) {
	g := NewGoroutineMulti()
	tests := []struct {
		timerTask *timerTask
		name      string
	}{
		{newTimerTask(g.scheduler, newTask(func() {
			fmt.Printf("go... %v\n", time.Now())
		}), -100, 1000), "Test_timerTask_schedule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.timerTask.schedule()
			<-time.After(3 * time.Second)
			tt.timerTask.Dispose()
			t1 := newTimerTask(g.scheduler, newTask(func() {
				fmt.Printf("go... %v\n", time.Now())
			}), 1000, 1000)
			t1.Dispose()
			<-time.After(1 * time.Second)
		})
	}
}

// Test_timerTask_schedule_withDelay validates delayed schedule path can be disposed before first fire.
// Test_timerTask_schedule_withDelay 驗證延遲排程可在首次觸發前安全 Dispose。
func Test_timerTask_schedule_withDelay(t *testing.T) {
	g := NewGoroutineMulti()
	tests := []struct {
		timerTask *timerTask
		name      string
	}{
		{newTimerTask(g.scheduler, newTask(func() {
			fmt.Printf("go... %v\n", time.Now())
		}), 100000, 10000), "Test_timerTask_schedule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.timerTask.schedule()
			<-time.After(1 * time.Second)
			tt.timerTask.Dispose()
			<-time.After(1 * time.Second)
		})
	}
}

// Test_timerTask validates interval execution count and dispose semantics for schedule APIs.
// Test_timerTask 驗證週期執行次數與各種 schedule API 的 Dispose 語義。
func Test_timerTask(t *testing.T) {
	var runCount int32
	g := NewGoroutineMulti()
	tests := []struct {
		fields *timerTask
		name   string
	}{
		{newTimerTask(g.scheduler, newTask(func() {
			atomic.AddInt32(&runCount, 1)
		}), 0, 5), "Test_timerTask_schedule_1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.schedule()
			for {
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
			<-time.After(30 * time.Millisecond)
			if int32(2) != atomic.LoadInt32(&runT1Count) {
				t.Fatal("they should be equal")
			}

			var runT2Count int32
			t2 := g.ScheduleOnInterval(1000, 10, func() {
				atomic.AddInt32(&runT2Count, 1)
			})
			t2.Dispose()
			<-time.After(30 * time.Millisecond)
			if int32(0) != atomic.LoadInt32(&runT2Count) {
				t.Fatal("they should be equal")
			}

			var runT3Count int32
			t3 := g.Schedule(1000, func() {
				atomic.AddInt32(&runT3Count, 1)
			})
			t3.Dispose()
			<-time.After(30 * time.Millisecond)
			if int32(0) != atomic.LoadInt32(&runT3Count) {
				t.Fatal("they should be equal")
			}

			var runT4Count int32
			t4 := g.Schedule(0, func() {
				atomic.AddInt32(&runT4Count, 1)
			})
			t4.Dispose()
			<-time.After(30 * time.Millisecond)
			if int32(1) != atomic.LoadInt32(&runT4Count) {
				t.Fatal("they should be equal")
			}
		})
	}
}

// Test_Task_execute validates reflection and fast-path task execution entry points.
// Test_Task_execute 驗證反射路徑與快速路徑的 task 執行入口。
func Test_Task_execute(t *testing.T) {
	tests := []struct {
		name string
		task task
	}{
		{task: newTask(exampleFunc1, "QQQQ", "aaa"), name: "Test_newTask_exampleFunc1"},
		{task: newTask(exampleFunc2), name: "Test_newTask_exampleFunc2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.task.execute()
		})
	}
}

// exampleFunc1 is a variadic helper used to exercise reflective invocation.
// exampleFunc1 是用來觸發反射呼叫路徑的可變參數輔助函式。
func exampleFunc1(args ...string) {
	fmt.Printf("exampleFunc1 %v", args)
}

// exampleFunc2 is a zero-argument helper used to exercise fast-path invocation.
// exampleFunc2 是用來觸發無參數快速路徑的輔助函式。
func exampleFunc2() {
	fmt.Printf("exampleFunc2 ")
}

// TestPanicHandler_recover verifies panic is captured and routed to custom PanicHandler.
// TestPanicHandler_recover 驗證 panic 會被攔截並交給自訂 PanicHandler。
func TestPanicHandler_recover(t *testing.T) {
	var captured any
	var capturedStack []byte
	original := getPanicHandler()
	SetPanicHandler(func(r any, stack []byte) {
		captured = r
		capturedStack = stack
	})
	defer SetPanicHandler(original)

	tk := newTask(func() { panic("boom") })
	tk.execute() // should not crash

	if captured != "boom" {
		t.Fatalf("PanicHandler got %v, want 'boom'", captured)
	}
	if len(capturedStack) == 0 {
		t.Fatal("PanicHandler received empty stack trace")
	}
}

// TestPanicHandler_nil_propagates verifies panic propagates when PanicHandler is disabled.
// TestPanicHandler_nil_propagates 驗證在停用 PanicHandler 時 panic 會向外傳遞。
func TestPanicHandler_nil_propagates(t *testing.T) {
	original := getPanicHandler()
	SetPanicHandler(nil)
	defer SetPanicHandler(original)

	defer func() {
		r := recover()
		if r != "boom" {
			t.Fatalf("expected panic 'boom', got %v", r)
		}
	}()

	tk := newTask(func() { panic("boom") })
	tk.execute() // should panic through
	t.Fatal("should not reach here")
}

// TestPanicHandler_fiber_survives verifies fiber continues processing subsequent tasks after one panic.
// TestPanicHandler_fiber_survives 驗證單一任務 panic 後 fiber 仍可繼續處理後續任務。
func TestPanicHandler_fiber_survives(t *testing.T) {
	original := getPanicHandler()
	var wg sync.WaitGroup
	wg.Add(2) // one for panic handler, one for normal task

	SetPanicHandler(func(r any, stack []byte) {
		wg.Done()
	})
	defer SetPanicHandler(original)

	g := NewGoroutineMulti()
	defer g.Dispose()

	// first task panics — PanicHandler calls wg.Done
	g.Enqueue(func() { panic("bad task") })
	// second task should still execute
	g.Enqueue(func() { wg.Done() })

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success: both panic handler and normal task completed
	case <-time.After(2 * time.Second):
		t.Fatal("timed out: fiber did not survive panic")
	}
}

// noOpDisposable is a stub Disposable used by mock scheduler in lock-order tests.
// noOpDisposable 是鎖順序測試中 mock scheduler 使用的空實作 Disposable。
type noOpDisposable struct{}

// Dispose implements Disposable with no side effects for test stubbing.
// Dispose 提供無副作用的 Disposable 實作供測試替身使用。
func (noOpDisposable) Dispose() {}

// lockCheckScheduler is a scheduler stub that asserts mutex unlock ordering before enqueue.
// lockCheckScheduler 是 scheduler 測試替身，用於斷言 enqueue 前已釋放 mutex。
type lockCheckScheduler struct {
	t            *testing.T
	timerTask    *timerTask
	enqueueCalls int32
}

// Schedule returns a no-op disposable to satisfy Scheduler in tests.
// Schedule 回傳 no-op disposable，用於滿足測試中的 Scheduler 介面。
func (s *lockCheckScheduler) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return noOpDisposable{}
}

// ScheduleOnInterval returns a no-op disposable to satisfy Scheduler in tests.
// ScheduleOnInterval 回傳 no-op disposable，用於滿足測試中的 Scheduler 介面。
func (s *lockCheckScheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return noOpDisposable{}
}

// Enqueue is unused in this stub because tests call enqueueTask directly.
// Enqueue 在此替身中不使用，因測試直接呼叫 enqueueTask。
func (s *lockCheckScheduler) Enqueue(taskFunc any, params ...any) {}

// enqueueTask asserts timerTask mutex is not held when scheduler enqueue is invoked.
// enqueueTask 會斷言呼叫 scheduler enqueue 時不應仍持有 timerTask 的 mutex。
func (s *lockCheckScheduler) enqueueTask(task task) {
	atomic.AddInt32(&s.enqueueCalls, 1)
	if !s.timerTask.mu.TryLock() {
		s.t.Fatal("executeOnFiber called enqueueTask while holding timerTask mutex")
	}
	s.timerTask.mu.Unlock()
}

// Remove is a no-op in this scheduler stub.
// Remove 在此 scheduler 替身中為 no-op。
func (s *lockCheckScheduler) Remove(d Disposable) {}

// Dispose is a no-op in this scheduler stub.
// Dispose 在此 scheduler 替身中為 no-op。
func (s *lockCheckScheduler) Dispose() {}

// Test_timerTask_executeOnFiber_unlockBeforeEnqueue verifies lock inversion fix on timerTask enqueue path.
// Test_timerTask_executeOnFiber_unlockBeforeEnqueue 驗證 timerTask enqueue 路徑的鎖反轉修正。
func Test_timerTask_executeOnFiber_unlockBeforeEnqueue(t *testing.T) {
	mock := &lockCheckScheduler{t: t}
	tk := newTimerTask(mock, newTask(func() {}), 0, -1)
	mock.timerTask = tk

	tk.executeOnFiber()

	if got := atomic.LoadInt32(&mock.enqueueCalls); got != 1 {
		t.Fatalf("enqueueTask called %d times, want 1", got)
	}
}
