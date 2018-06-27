package robin

import (
	"sync"
)

type GoroutineMulti struct {
	queue          TaskQueue
	scheduler      IScheduler
	executor       executor
	executionState executionState
	lock           *sync.Mutex
	subscriptions  *Disposer
	flushPending   bool
}

func (g *GoroutineMulti) init() *GoroutineMulti {
	g.queue = NewDefaultQueue()
	g.executionState = created
	g.scheduler = NewScheduler(g)
	g.executor = newDefaultExecutor()
	g.subscriptions = NewDisposer()
	g.lock = new(sync.Mutex)
	return g
}

func NewGoroutineMulti() *GoroutineMulti {
	return new(GoroutineMulti).init()
}

func (g *GoroutineMulti) Start() {
	if g.executionState == running {
		return
	}
	g.executionState = running
	g.scheduler.Start()
	g.Enqueue(func() {})
}

func (g *GoroutineMulti) Stop() {
	g.executionState = stopped
}

func (g *GoroutineMulti) Dispose() {
	g.Stop()
	g.scheduler.Dispose()
	g.subscriptions.Dispose()
	g.queue.Dispose()
}

func (g *GoroutineMulti) Enqueue(taskFun interface{}, params ...interface{}) {
	g.EnqueueWithTask(newTask(taskFun, params...))
}

func (g *GoroutineMulti) EnqueueWithTask(task task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	if g.flushPending {
		return
	}
	g.flushPending = true
	go g.flush()
}

func (g *GoroutineMulti) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement SubscriptionRegistry.RegisterSubscription */
func (g *GoroutineMulti) RegisterSubscription(toAdd Disposable) {
	g.subscriptions.Add(toAdd)
}

/*implement SubscriptionRegistry.DeregisterSubscription */
func (g *GoroutineMulti) DeregisterSubscription(toRemove Disposable) {
	g.subscriptions.Remove(toRemove)
}

func (g GoroutineMulti) NumSubscriptions() int {
	return g.subscriptions.Count()
}

func (g *GoroutineMulti) flush() {
	toDoTasks, ok := g.queue.DequeueAll()
	if !ok {
		g.flushPending = false
		return
	}
	g.executor.ExecuteTasksWithGoroutine(toDoTasks)
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.queue.Count() > 0 {
		//It has new task enqueue when clear tasks
		go g.flush()
	} else {
		//task is empty
		g.flushPending = false
	}
}
