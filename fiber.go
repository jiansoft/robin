package robin

import (
	"sync"
)

// Fiber defines the interface for task execution fibers.
type Fiber interface {
	Dispose()
	Enqueue(taskFunc any, params ...any)
	enqueueTask(task task)
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
}

// fiberCommon holds shared state for all fiber implementations.
type fiberCommon struct {
	queue      taskQueue
	scheduler  Scheduler
	executor   executor
	mu         sync.Mutex
	isDisposed bool
}

// GoroutineMulti is a fiber backed by multiple goroutines. Each task batch is dispatched concurrently.
type GoroutineMulti struct {
	fiberCommon
	flushPending bool
	flushFn      func() // cached bound method, avoids per-enqueue closure allocation
}

// GoroutineSingle is a fiber backed by a single dedicated goroutine. Tasks execute serially in order.
type GoroutineSingle struct {
	cond *sync.Cond
	fiberCommon
}

// NewGoroutineMulti creates a new GoroutineMulti fiber instance.
func NewGoroutineMulti() *GoroutineMulti {
	g := new(GoroutineMulti)
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	g.flushFn = g.flush

	return g
}

// Dispose stops the fiber, cancels all scheduled tasks, and releases resources.
func (g *GoroutineMulti) Dispose() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.isDisposed = true
	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue submits a task for concurrent execution on this fiber.
func (g *GoroutineMulti) Enqueue(taskFunc any, params ...any) {
	g.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask adds a task to the queue and triggers a flush if none is pending.
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

// Schedule executes the task once after firstInMs milliseconds.
func (g *GoroutineMulti) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval executes the task after firstInMs milliseconds, then repeats every regularInMs milliseconds.
func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

// flush drains all pending tasks from the queue and dispatches them via goroutines.
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

// NewGoroutineSingle creates a new GoroutineSingle fiber instance with a dedicated goroutine loop.
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

// Dispose stops the fiber, signals the goroutine to exit, and releases resources.
func (g *GoroutineSingle) Dispose() {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.isDisposed = true
	g.cond.Signal()
	g.scheduler.Dispose()
	g.queue.dispose()
}

// Enqueue submits a task for serial execution on this fiber's dedicated goroutine.
func (g *GoroutineSingle) Enqueue(taskFunc any, params ...any) {
	g.enqueueTask(newTask(taskFunc, params...))
}

// enqueueTask adds a task to the queue and signals the dedicated goroutine.
func (g *GoroutineSingle) enqueueTask(task task) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	if g.isDisposed {
		return
	}

	g.queue.enqueue(task)
	g.cond.Signal()
}

// Schedule executes the task once after firstInMs milliseconds.
func (g *GoroutineSingle) Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFunc, params...)
}

// ScheduleOnInterval executes the task after firstInMs milliseconds, then repeats every regularInMs milliseconds.
func (g *GoroutineSingle) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFunc, params...)
}

// executeNextBatch dequeues and executes one batch of tasks. Returns false when disposed.
func (g *GoroutineSingle) executeNextBatch() bool {
	tasks, ok := g.dequeueAll()
	if ok {
		g.executor.executeTasks(tasks)
	}

	return ok
}

// dequeueAll waits for tasks and returns them all. Returns (nil, false) when disposed.
func (g *GoroutineSingle) dequeueAll() ([]task, bool) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.waitingForTask()

	if g.isDisposed {
		return nil, false
	}

	return g.queue.dequeueAll()
}

// waitingForTask blocks on sync.Cond until tasks are available or the fiber is disposed.
func (g *GoroutineSingle) waitingForTask() {
	for g.queue.count() == 0 {
		if g.isDisposed {
			return
		}

		g.cond.Wait()
	}
}
