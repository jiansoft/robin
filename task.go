package robin

import (
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// panicHandler stores the current panic handler in an atomic container for lock-free reads.
// The stored value type is func(any, []byte); typed nil means "panic should propagate".
// panicHandler 以 atomic 容器保存目前的 panic handler，提供無鎖讀取。
// 儲存型別為 func(any, []byte)；typed nil 代表「panic 應直接往外拋」。
var panicHandler atomic.Value

// init installs the default panic handler that logs panic details to stderr.
// init 會安裝預設 panic handler，將 panic 資訊與堆疊輸出到 stderr。
func init() {
	SetPanicHandler(func(r any, stack []byte) {
		_, _ = fmt.Fprintf(os.Stderr, "robin: task panic: %v\n%s\n", r, stack)
	})
}

// SetPanicHandler sets the callback used to handle panics raised by task execution.
// The default behavior prints panic value + stack trace to stderr.
// Passing nil disables recovery and lets panics propagate as normal.
// SetPanicHandler 用來設定任務執行發生 panic 時的處理函式。
// 預設行為會將 panic 值與 stack trace 輸出到 stderr。
// 傳入 nil 可停用 recover，讓 panic 依原生流程向外拋出。
func SetPanicHandler(h func(r any, stack []byte)) {
	if h == nil {
		panicHandler.Store((func(any, []byte))(nil))
		return
	}
	panicHandler.Store(h)
}

// getPanicHandler returns the current panic handler function.
// A nil return means panic recovery is disabled.
// getPanicHandler 會回傳目前的 panic handler。
// 若回傳 nil，表示不進行 recover。
func getPanicHandler() func(any, []byte) {
	h, _ := panicHandler.Load().(func(any, []byte))
	return h
}

// task represents a unit of work to be executed by a Fiber.
// It supports three execution paths:
//   - Fast path: direct function call via fn (for func() with no parameters)
//   - Typed path: invoke(invokeArg) dispatch via any (for TypedChannel, no reflection)
//   - Slow path: reflection-based call via funcCache and paramsCache (for arbitrary signatures)
//
// task 是 fiber 內部的任務單位，具備三種執行路徑：
//   - 快速路徑：fn 直接呼叫（func() 且無參數）
//   - 型別路徑：invoke(invokeArg)（TypedChannel 無反射）
//   - 反射路徑：funcCache + paramsCache（支援任意簽章）
type task struct {
	fn          func() // fast path: no reflection needed
	invoke      func(any)
	invokeArg   any
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

// newTask builds a task from an arbitrary function + parameters.
// It selects the no-reflection fast path when f is func() and params are empty.
// newTask 會由任意函式與參數建立 task。
// 當 f 為 func() 且無參數時，會自動走無反射的快速路徑。
func newTask(f any, p ...any) task {
	if fn, ok := f.(func()); ok && len(p) == 0 {
		return task{fn: fn}
	}
	t := task{funcCache: reflect.ValueOf(f)}
	t.params(p...)

	return t
}

// params converts variadic values into reflect.Value slice for the reflection path.
// params 會把可變參數轉成 reflect.Value 切片，供反射呼叫使用。
func (t *task) params(p ...any) {
	t.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		t.paramsCache[k] = reflect.ValueOf(param)
	}
}

// execute runs the task. If a PanicHandler is set via SetPanicHandler,
// panics are recovered and forwarded to the handler; otherwise panics propagate normally.
// execute 會執行 task；若已設定 PanicHandler 則會 recover 並回呼 handler，
// 否則 panic 會依預設行為向外傳遞。
func (t *task) execute() {
	defer func() {
		h := getPanicHandler()
		if h == nil {
			return
		}
		if r := recover(); r != nil {
			h(r, debug.Stack())
		}
	}()
	if t.fn != nil {
		t.fn()
		return
	}
	if t.invoke != nil {
		t.invoke(t.invokeArg)
		return
	}
	_ = t.funcCache.Call(t.paramsCache)
}

// timerTask wraps a task with timer-based scheduling state.
// It supports one-shot and repeating execution via time.AfterFunc.
// timerTask 封裝了任務的 timer 排程狀態，
// 支援一次性與重複性執行，底層使用 time.AfterFunc。
type timerTask struct {
	scheduler    Scheduler
	task         task
	timer        *time.Timer
	firstInMs    int64
	intervalInMs int64
	mu           sync.Mutex
	disposed     atomic.Bool
}

// newTimerTask creates a timerTask that will execute the given task
// after firstInMs milliseconds, then repeat every intervalInMs milliseconds.
// If intervalInMs is 0 or less, the task executes only once.
// newTimerTask 會建立 timerTask：
// firstInMs 後執行第一次，之後每 intervalInMs 重複；
// 若 intervalInMs <= 0 則僅執行一次。
func newTimerTask(scheduler Scheduler, task task, firstInMs int64, intervalInMs int64) *timerTask {
	return &timerTask{
		scheduler:    scheduler,
		task:         task,
		firstInMs:    firstInMs,
		intervalInMs: intervalInMs,
	}
}

// Dispose releases the resources associated with the timerTask.
// It stops the timer and removes the task from the scheduler.
// Safe to call multiple times; only the first call takes effect.
// Dispose 會釋放 timerTask 資源：停止 timer，並移出 scheduler 追蹤。
// 可重複呼叫，只有第一次會生效。
func (tk *timerTask) Dispose() {
	if !tk.disposed.CompareAndSwap(false, true) {
		return
	}

	tk.mu.Lock()
	if tk.timer != nil {
		tk.timer.Stop()
	}
	s := tk.scheduler
	tk.scheduler = nil
	tk.mu.Unlock()
	if s != nil {
		s.Remove(tk)
	}
}

// schedule starts the execution of the timerTask.
// If firstInMs is 0 or less, it executes the task immediately,
// otherwise it schedules the task to run after firstInMs milliseconds.
// schedule 會啟動 timerTask：
// firstInMs <= 0 時立即進入 next()，否則延遲指定毫秒後觸發。
func (tk *timerTask) schedule() {
	if tk.firstInMs <= 0 {
		tk.next()
		return
	}

	tk.mu.Lock()
	if !tk.disposed.Load() {
		tk.timer = time.AfterFunc(time.Duration(tk.firstInMs)*time.Millisecond, tk.next)
	}
	tk.mu.Unlock()
}

// next executes the task and schedules the next execution if recurring.
// If the task is not recurring (intervalInMs is 0 or less), it disposes the task.
// next 會執行目前任務，若為週期任務則安排下一次；
// 若不是週期任務（intervalInMs <= 0）則直接 Dispose。
func (tk *timerTask) next() {
	tk.executeOnFiber()
	if tk.intervalInMs <= 0 {
		tk.Dispose()
		return
	}

	tk.mu.Lock()
	if !tk.disposed.Load() {
		tk.timer = time.AfterFunc(time.Duration(tk.intervalInMs)*time.Millisecond, tk.next)
	}
	tk.mu.Unlock()
}

// executeOnFiber enqueues the task in the scheduler for execution.
// It is a no-op if the timerTask has been disposed.
// executeOnFiber 會把任務送回 scheduler 的 fiber 佇列中執行。
// 若 timerTask 已 dispose 則此方法不做任何事。
func (tk *timerTask) executeOnFiber() {
	if tk.disposed.Load() {
		return
	}

	tk.mu.Lock()
	s := tk.scheduler
	t := tk.task
	tk.mu.Unlock()

	// Do not call enqueueTask while holding tk.mu.
	// This keeps lock ordering consistent with fiber Dispose paths and avoids lock inversion.
	if s != nil && !tk.disposed.Load() {
		s.enqueueTask(t)
	}
}
