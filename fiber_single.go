package robin

import (
	"sync"
)

type GoroutineSingle struct {
	queue          TaskQueue
	scheduler      IScheduler
	executor       Executor
	executionState executionState
	lock           *sync.Mutex
	cond           *sync.Cond
	subscriptions  *Disposer
}

func (g *GoroutineSingle) init() *GoroutineSingle {
	g.queue = NewDefaultQueue()
	g.executionState = created
	g.subscriptions = NewDisposer()
	g.scheduler = NewScheduler(g)
	g.executor = NewDefaultExecutor()
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
	g.EnqueueWithTask(newTask(taskFun, params...))
}

func (g *GoroutineSingle) EnqueueWithTask(task task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	//喚醒等待中的 Goroutine
	g.cond.Broadcast()
}

func (g GoroutineSingle) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (g GoroutineSingle) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement SubscriptionRegistry.RegisterSubscription */
func (g *GoroutineSingle) RegisterSubscription(toAdd Disposable) {
	g.subscriptions.Add(toAdd)
}

/*implement SubscriptionRegistry.DeregisterSubscription */
func (g *GoroutineSingle) DeregisterSubscription(toRemove Disposable) {
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

func (g *GoroutineSingle) dequeueAll() ([]task, bool) {
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
