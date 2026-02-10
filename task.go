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

// panicHandler stores the current panic handler using atomic.Value for concurrent safety.
// The stored type is always func(any, []byte); a typed nil represents "no handler".
var panicHandler atomic.Value

func init() {
	SetPanicHandler(func(r any, stack []byte) {
		_, _ = fmt.Fprintf(os.Stderr, "robin: task panic: %v\n%s\n", r, stack)
	})
}

// SetPanicHandler sets the handler called when a task panics during execution.
// The default handler prints the panic value and stack trace to stderr.
// Pass nil to let panics propagate normally (crash the process).
func SetPanicHandler(h func(r any, stack []byte)) {
	if h == nil {
		panicHandler.Store((func(any, []byte))(nil))
		return
	}
	panicHandler.Store(h)
}

// getPanicHandler returns the current panic handler, or nil if not set.
func getPanicHandler() func(any, []byte) {
	h, _ := panicHandler.Load().(func(any, []byte))
	return h
}

// task represents a unit of work to be executed by a Fiber.
// It supports three execution paths:
//   - Fast path: direct function call via fn (for func() with no parameters)
//   - Typed path: invoke(invokeArg) dispatch via any (for TypedChannel, no reflection)
//   - Slow path: reflection-based call via funcCache and paramsCache (for arbitrary signatures)
type task struct {
	fn          func() // fast path: no reflection needed
	invoke      func(any)
	invokeArg   any
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

// newTask returns a task instance.
// If f is func() and there are no parameters, it uses the fast path (no reflection).
func newTask(f any, p ...any) task {
	if fn, ok := f.(func()); ok && len(p) == 0 {
		return task{fn: fn}
	}
	t := task{funcCache: reflect.ValueOf(f)}
	t.params(p...)

	return t
}

// params converts variadic arguments to reflect.Value slice for reflection-based invocation.
func (t *task) params(p ...any) {
	t.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		t.paramsCache[k] = reflect.ValueOf(param)
	}
}

// execute runs the task. If a PanicHandler is set via SetPanicHandler,
// panics are recovered and forwarded to the handler; otherwise panics propagate normally.
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

// timerTask wraps a task with timer-based scheduling.
// It supports both one-shot (firstInMs only) and recurring (intervalInMs) execution,
// using time.AfterFunc for efficient timer management.
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
func (tk *timerTask) executeOnFiber() {
	if tk.disposed.Load() {
		return
	}

	tk.mu.Lock()
	if tk.scheduler != nil {
		tk.scheduler.enqueueTask(tk.task)
	}
	tk.mu.Unlock()
}
