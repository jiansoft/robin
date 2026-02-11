package robin

import (
	"sync"
)

// Fiber defines the public execution contract for task dispatchers.
// A fiber can enqueue immediate tasks and schedule delayed/interval tasks.
// Fiber 定義任務派送器的公開執行契約。
// fiber 可接收立即任務，也可排程延遲/週期任務。
type Fiber interface {
	Dispose()
	Enqueue(taskFunc any, params ...any)
	enqueueTask(task task)
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
}

// fiberCommon contains shared fields reused by concrete fiber implementations.
// fiberCommon 保存各 fiber 實作共用的內部狀態欄位。
type fiberCommon struct {
	queue      taskQueue
	scheduler  Scheduler
	executor   executor
	mu         sync.Mutex
	isDisposed bool
}

// GoroutineMulti is a fiber backed by multi-goroutine execution.
// It batches queue tasks and dispatches each task concurrently.
// GoroutineMulti 是多 goroutine 的 fiber 實作。
// 它會以批次方式提取任務，並讓批次內任務併發執行。
type GoroutineMulti struct {
	fiberCommon
	flushPending bool
	flushFn      func() // cached bound method, avoids per-enqueue closure allocation
}

// GoroutineSingle is a fiber with one dedicated worker goroutine.
// Tasks are executed serially and preserve enqueue order.
// GoroutineSingle 是單一專用 worker goroutine 的 fiber。
// 任務以序列方式執行，並維持 enqueue 順序。
type GoroutineSingle struct {
	cond *sync.Cond
	fiberCommon
}

// NewGoroutineMulti creates a multi-goroutine fiber with default queue/scheduler/executor.
// NewGoroutineMulti 會建立使用預設 queue/scheduler/executor 的多 goroutine fiber。
func NewGoroutineMulti() *GoroutineMulti {
	g := new(GoroutineMulti)
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	g.flushFn = g.flush

	return g
}

// Dispose stops the fiber, cancels scheduled work, and releases queue resources.
// Dispose 可安全重複呼叫，後續呼叫為 no-op。
func (g *GoroutineMulti) Dispose() {
	g.mu.Lock()
	if g.isDisposed {
		g.mu.Unlock()
		return
	}
	g.isDisposed = true
	g.mu.Unlock()

	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue submits one task into this fiber for asynchronous execution.
// Enqueue 會把一個任務送入此 fiber 進行非同步執行。
func (g *GoroutineMulti) Enqueue(taskFunc any, params ...any) {
	g.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask appends a prepared task and lazily starts one flush loop.
// enqueueTask 會加入已建構任務，且在沒有 flush 進行時才啟動一次 flush。
func (g *GoroutineMulti) enqueueTask(t task) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isDisposed {
		return
	}

	g.queue.enqueue(t)

	if g.flushPending {
		return
	}

	g.flushPending = true

	g.executor.executeTaskWithGoroutine(newTask(g.flushFn))
}

// Schedule executes a one-shot task after firstInMs milliseconds.
// Schedule 會在 firstInMs 毫秒後執行一次任務。
func (g *GoroutineMulti) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval executes a task after firstInMs and then repeatedly every regularInMs.
// ScheduleOnInterval 會在 firstInMs 後首次執行，之後每 regularInMs 毫秒重複執行。
func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

// flush drains queue batches and dispatches each batch's tasks concurrently.
// flush 會持續清空佇列批次，並把批次任務以併發方式派送。
func (g *GoroutineMulti) flush() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for {
		toDoTasks, ok := g.queue.dequeueAll()
		if !ok {
			g.flushPending = false
			return
		}

		g.executor.executeTasksWithGoroutine(toDoTasks)
		clearTaskBatch(toDoTasks)

		if g.queue.count() == 0 {
			g.flushPending = false
			return
		}
	}
}

// NewGoroutineSingle creates a serial fiber with one dedicated worker loop.
// NewGoroutineSingle 會建立單一 worker 迴圈的序列 fiber。
func NewGoroutineSingle() *GoroutineSingle {
	g := new(GoroutineSingle)
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	g.cond = sync.NewCond(&g.mu)

	go func() {
		for g.executeNextBatch() {
		}
	}()

	return g
}

// Dispose signals the worker loop to exit, then disposes scheduler and queue.
// Dispose 會通知 worker 迴圈結束，接著釋放 scheduler 與 queue。
func (g *GoroutineSingle) Dispose() {
	g.cond.L.Lock()
	if g.isDisposed {
		g.cond.L.Unlock()
		return
	}

	g.isDisposed = true
	g.cond.Signal()
	g.cond.L.Unlock()

	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue submits a task for ordered serial execution.
// Enqueue 會把任務提交到序列執行流程，依序處理。
func (g *GoroutineSingle) Enqueue(taskFunc any, params ...any) {
	g.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask pushes one task into queue and wakes the waiting worker.
// enqueueTask 會把任務加入佇列，並喚醒等待中的 worker。
func (g *GoroutineSingle) enqueueTask(task task) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	if g.isDisposed {
		return
	}

	g.queue.enqueue(task)
	g.cond.Signal()
}

// Schedule executes a one-shot task after firstInMs milliseconds.
// Schedule 會在 firstInMs 毫秒後執行一次任務。
func (g *GoroutineSingle) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval executes a task after firstInMs and then repeatedly every regularInMs.
// ScheduleOnInterval 會在 firstInMs 後首次執行，之後每 regularInMs 毫秒重複執行。
func (g *GoroutineSingle) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

// executeNextBatch processes one queue batch.
// It returns false when the fiber has been disposed and should stop its loop.
// executeNextBatch 會處理一個任務批次；
// 當 fiber 已 dispose 時回傳 false 以停止主迴圈。
func (g *GoroutineSingle) executeNextBatch() bool {
	tasks, ok := g.dequeueAll()
	if ok {
		g.executor.executeTasks(tasks)
		clearTaskBatch(tasks)
	}

	return ok
}

// clearTaskBatch zeroes task slots so captured references can be released early.
// clearTaskBatch 會把 task 槽位清零，讓捕捉到的引用可提早釋放。
func clearTaskBatch(tasks []task) {
	for i := range tasks {
		tasks[i] = task{}
	}
}

// dequeueAll blocks until tasks are available or the fiber is disposed.
// It returns (nil, false) when disposed.
// dequeueAll 會阻塞等待任務或等待 dispose 發生。
// 當已 dispose 時回傳 (nil, false)。
func (g *GoroutineSingle) dequeueAll() ([]task, bool) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.waitingForTask()

	if g.isDisposed {
		return nil, false
	}

	return g.queue.dequeueAll()
}

// waitingForTask waits on sync.Cond while queue is empty and fiber is still active.
// waitingForTask 會在佇列為空且 fiber 仍有效時，透過 sync.Cond 持續等待。
func (g *GoroutineSingle) waitingForTask() {
	for g.queue.count() == 0 {
		if g.isDisposed {
			return
		}

		g.cond.Wait()
	}
}
