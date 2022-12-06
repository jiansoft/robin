package robin

import "sync"

const (
	created executionState = iota
	running
	stopped
)

type executionState int

// Fiber define some function
type Fiber interface {
	Start()
	Stop()
	Dispose()
	Enqueue(taskFunc any, params ...any)
	EnqueueWithTask(task Task)
	Schedule(firstInMs int64, taskFunc any, params ...any) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFunc any, params ...any) (d Disposable)
}

// GoroutineMulti a fiber backed by more goroutine. Each job is executed by a new goroutine.
type GoroutineMulti struct {
	queue          taskQueue
	scheduler      IScheduler
	executor       executor
	executionState executionState
	mu             sync.Mutex
	flushPending   bool
}

// GoroutineSingle a fiber backed by a dedicated goroutine. Every job is executed by a goroutine.
type GoroutineSingle struct {
	queue          taskQueue
	scheduler      IScheduler
	executor       executor
	executionState executionState
	mu             sync.Mutex
	cond           *sync.Cond
}

// NewGoroutineMulti create a GoroutineMulti instance
func NewGoroutineMulti() *GoroutineMulti {
	g := new(GoroutineMulti)
	g.queue = newDefaultQueue()
	g.executionState = created
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	return g
}

// Start the fiber work now
func (g *GoroutineMulti) Start() {
	if g.executionState == running {
		return
	}
	g.executionState = running
}

// Stop the fiber work
func (g *GoroutineMulti) Stop() {
	g.executionState = stopped
}

// Dispose stop the fiber and release resource
func (g *GoroutineMulti) Dispose() {
	g.Stop()
	g.scheduler.Dispose()
	g.queue.Dispose()
}

// Enqueue use the fiber to execute a task
func (g *GoroutineMulti) Enqueue(taskFunc any, params ...any) {
	g.EnqueueWithTask(newTask(taskFunc, params...))
}

// EnqueueWithTask use the fiber to execute a task
func (g *GoroutineMulti) EnqueueWithTask(task Task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.flushPending {
		return
	}
	g.flushPending = true
	g.executor.ExecuteTaskWithGoroutine(newTask(g.flush))
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
	toDoTasks, ok := g.queue.DequeueAll()
	if !ok {
		g.flushPending = false
		return
	}
	g.executor.ExecuteTasksWithGoroutine(toDoTasks)
	if g.queue.Count() > 0 {
		//It has new Task enqueue when clear tasks
		g.executor.ExecuteTaskWithGoroutine(newTask(g.flush))
		//go g.flush()
	} else {
		//Task is empty
		g.flushPending = false
	}
}

// NewGoroutineSingle create a GoroutineSingle instance
func NewGoroutineSingle() *GoroutineSingle {
	g := new(GoroutineSingle)
	g.executionState = created
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	g.cond = sync.NewCond(&g.mu)
	return g
}

// Start the fiber work now
func (g *GoroutineSingle) Start() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.executionState == running {
		return
	}
	g.executionState = running
	go func() {
		for g.executeNextBatch() {
		}
	}()
}

// Stop the fiber work
func (g *GoroutineSingle) Stop() {
	g.mu.Lock()
	g.executionState = stopped
	g.cond.Broadcast()
	g.mu.Unlock()
}

// Dispose stop the fiber and release resource
func (g *GoroutineSingle) Dispose() {
	g.mu.Lock()
	g.executionState = stopped
	g.cond.Broadcast()
	g.mu.Unlock()
	g.scheduler.Dispose()
	g.queue.Dispose()
}

// Enqueue use the fiber to execute a task
func (g *GoroutineSingle) Enqueue(taskFunc any, params ...any) {
	g.EnqueueWithTask(newTask(taskFunc, params...))
}

// EnqueueWithTask enqueue the parameter task
// into the queue waiting for executing.
func (g *GoroutineSingle) EnqueueWithTask(task Task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	//Wake up the waiting goroutine
	g.cond.Broadcast()
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
		g.executor.ExecuteTasks(tasks)
	}
	return ok
}

func (g *GoroutineSingle) dequeueAll() ([]Task, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.readyToDequeue() {
		return nil, false
	}
	return g.queue.DequeueAll()
}

func (g *GoroutineSingle) readyToDequeue() bool {
	//若貯列中已沒有任務要執行時，將 Goroutine 進入等侯訊號的狀態
	tmp := g.executionState == running
	for g.queue.Count() == 0 && tmp {
		g.cond.Wait()
	}
	return tmp
}
