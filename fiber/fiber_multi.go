package fiber

import (
	"sync"

	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/core"
)

type GoroutineMulti struct {
	queue          core.TaskQueue
	scheduler      core.IScheduler
	executor       core.Executor
	executionState executionState
	lock           *sync.Mutex
	subscriptions  *robin.Disposer
	flushPending   bool
}

func (g *GoroutineMulti) init() *GoroutineMulti {
	g.queue = core.NewDefaultQueue()
	g.executionState = created
	g.scheduler = core.NewScheduler(g)
	g.executor = core.NewDefaultExecutor()
	g.subscriptions = robin.NewDisposer()
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
	g.EnqueueWithTask(core.NewTask(taskFun, params...))
}

func (g *GoroutineMulti) EnqueueWithTask(task core.Task) {
	if g.executionState != running {
		return
	}
	//g.lock.Lock()
	//defer g.lock.Unlock()
	g.queue.Enqueue(task)
	if g.flushPending {
		return
	}
	g.flushPending = true
	go g.flush()
}

func (g *GoroutineMulti) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement SubscriptionRegistry.RegisterSubscription */
func (g *GoroutineMulti) RegisterSubscription(toAdd robin.Disposable) {
	g.subscriptions.Add(toAdd)
}

/*implement SubscriptionRegistry.DeregisterSubscription */
func (g *GoroutineMulti) DeregisterSubscription(toRemove robin.Disposable) {
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
		//Task is empty
		g.flushPending = false
	}
}
