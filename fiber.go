package robin

import (
	"sync"
)

// Fiber define some function
type Fiber interface {
	// Deprecated: This method is no longer used.
	Start()
	// Deprecated: This method is no longer used.
	Stop()
	Dispose()
	Enqueue(taskFunc any, params ...any)
	EnqueueWithTask(task Task)
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
}

type fiberCommon struct {
	queue      taskQueue
	scheduler  IScheduler
	executor   executor
	mu         sync.Mutex
	isDisposed bool
}

// GoroutineMulti a fiber backed by more goroutine. Each job is executed by a new goroutine.
type GoroutineMulti struct {
	fiberCommon
	flushPending bool
}

// GoroutineSingle a fiber backed by a dedicated goroutine. Every job is executed by a goroutine.
type GoroutineSingle struct {
	cond *sync.Cond
	fiberCommon
}

// NewGoroutineMulti create a GoroutineMulti instance
func NewGoroutineMulti() *GoroutineMulti {
	g := new(GoroutineMulti)
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()

	return g
}

// Deprecated: This method is no longer used.
func (g *GoroutineMulti) Start() {
}

// Deprecated: This method is no longer used.
func (g *GoroutineMulti) Stop() {
}

// Dispose stop the fiber and release resource
func (g *GoroutineMulti) Dispose() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.isDisposed = true
	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue use the fiber to execute a task
func (g *GoroutineMulti) Enqueue(taskFunc any, params ...any) {
	g.EnqueueWithTask(newTask(taskFunc, params...))
}

// EnqueueWithTask use the fiber to execute a task
func (g *GoroutineMulti) EnqueueWithTask(task Task) {
	flushTask := newTask(g.flush)

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isDisposed {
		return
	}

	g.queue.enqueue(task)

	if g.flushPending {
		return
	}

	g.flushPending = true

	g.executor.executeTaskWithGoroutine(flushTask)
}

// Schedule execute the task once at the specified time
// that depends on parameter firstInMs.
func (g *GoroutineMulti) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval execute the task once at the specified time
// that depends on parameters both firstInMs and regularInMs.
func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

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

		if g.queue.count() == 0 {
			g.flushPending = false
			return
		}
	}
}

// NewGoroutineSingle create a GoroutineSingle instance
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

// Deprecated: This method is no longer used.
func (g *GoroutineSingle) Start() {
}

// Deprecated: This method is no longer used.
func (g *GoroutineSingle) Stop() {
}

// Dispose stop the fiber and release resource
func (g *GoroutineSingle) Dispose() {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.isDisposed = true
	g.cond.Signal()
	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue use the fiber to execute a task
func (g *GoroutineSingle) Enqueue(taskFunc any, params ...any) {
	g.EnqueueWithTask(newTask(taskFunc, params...))
}

// EnqueueWithTask enqueue the parameter task
// into the queue waiting for executing.
func (g *GoroutineSingle) EnqueueWithTask(task Task) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	if g.isDisposed {
		return
	}

	g.queue.enqueue(task)
	//Wake up the waiting goroutine
	g.cond.Signal()
}

// Schedule execute the task once at the specified time
// that depends on parameter firstInMs.
func (g *GoroutineSingle) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval execute the task once at the specified time
// that depends on parameters both firstInMs and regularInMs.
func (g *GoroutineSingle) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

func (g *GoroutineSingle) executeNextBatch() bool {
	tasks, ok := g.dequeueAll()
	if ok {
		g.executor.executeTasks(tasks)
	}

	return ok
}

func (g *GoroutineSingle) dequeueAll() ([]Task, bool) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.waitingForTask()

	if g.isDisposed {
		return nil, false
	}

	return g.queue.dequeueAll()
}

func (g *GoroutineSingle) waitingForTask() {
	//若貯列中已沒有任務要執行時，將 Goroutine 進入等侯訊號的狀態
	for g.queue.count() == 0 {
		if g.isDisposed {
			return
		}

		g.cond.Wait()
	}
}
