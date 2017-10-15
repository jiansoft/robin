package fiber

import (
	"sync"

	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/core"
)

type GoroutineSingle struct {
	queue          core.IQueue
	scheduler      core.IScheduler
	executor       core.IExecutor
	executionState executionState
	lock           *sync.Mutex
	cond           *sync.Cond
	subscriptions  *core.Disposer
}

func (g *GoroutineSingle) init() *GoroutineSingle {
	g.queue = core.NewDefaultQueue()
	g.executionState = created
	g.subscriptions = core.NewDisposer()
	g.scheduler = core.NewScheduler(g)
	g.executor = core.NewDefaultExecutor()
	g.lock = new(sync.Mutex)
	g.cond = sync.NewCond(g.lock)
	return g
}

func NewGoroutineSingle() *GoroutineSingle {
	return new(GoroutineSingle).init()
}

func (g *GoroutineSingle) Start() {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.executionState == running {
		return
	}
	g.executionState = running
	g.scheduler.Start()
	go func() {
		for g.executeNextBatch() {
		}
	}()
}
func (g *GoroutineSingle) Stop() {
	g.lock.Lock()
	g.executionState = stopped
	g.cond.Broadcast()
	g.lock.Unlock()
	g.scheduler.Stop()
}

func (g *GoroutineSingle) Dispose() {
	g.lock.Lock()
	g.executionState = stopped
	g.cond.Broadcast()
	g.lock.Unlock()
	g.scheduler.Dispose()
	g.subscriptions.Dispose()
	g.queue.Dispose()
}

func (g *GoroutineSingle) Enqueue(taskFun interface{}, params ...interface{}) {
	g.EnqueueWithTask(core.NewTask(taskFun, params...))
}

func (g *GoroutineSingle) EnqueueWithTask(task core.Task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	//喚醒等待中的 Goroutine
	g.cond.Broadcast()
}

func (g GoroutineSingle) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (g GoroutineSingle) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement ISubscriptionRegistry.RegisterSubscription */
func (g *GoroutineSingle) RegisterSubscription(toAdd robin.Disposable) {
	g.subscriptions.Add(toAdd)
}

/*implement ISubscriptionRegistry.DeregisterSubscription */
func (g *GoroutineSingle) DeregisterSubscription(toRemove robin.Disposable) {
	g.subscriptions.Remove(toRemove)
}

func (g *GoroutineSingle) NumSubscriptions() int {
	return g.subscriptions.Count()
}

func (g *GoroutineSingle) executeNextBatch() bool {
	tasks, ok := g.dequeueAll()
	if ok {
		g.executor.ExecuteTasks(tasks)
	}
	return ok
}

func (g *GoroutineSingle) dequeueAll() ([]core.Task, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()
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
