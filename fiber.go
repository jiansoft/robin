package robin

import "sync"

const (
	created executionState = iota
	running
	stopped
)

type executionState int

type Fiber interface {
	Start()
	Stop()
	Dispose()
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
}

type GoroutineMulti struct {
	queue          taskQueue
	scheduler      IScheduler
	executor       executor
	executionState executionState
	lock           *sync.Mutex
	subscriptions  *Disposer
	flushPending   bool
}

type GoroutineSingle struct {
	queue          taskQueue
	scheduler      IScheduler
	executor       executor
	executionState executionState
	lock           *sync.Mutex
	cond           *sync.Cond
	subscriptions  *Disposer
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

func (g *GoroutineMulti) EnqueueWithTask(task Task) {
	if g.executionState != running {
		return
	}
	g.queue.Enqueue(task)
	g.lock.Lock()
	defer g.lock.Unlock()
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
		//It has new Task enqueue when clear tasks
		go g.flush()
	} else {
		//Task is empty
		g.flushPending = false
	}
}

func (g *GoroutineSingle) init() *GoroutineSingle {
	g.queue = NewDefaultQueue()
	g.executionState = created
	g.subscriptions = NewDisposer()
	g.scheduler = NewScheduler(g)
	g.executor = newDefaultExecutor()
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

// EnqueueWrap from parameters taskFun and params
// to a task and into to the queue waiting for executing.
func (g *GoroutineSingle) Enqueue(taskFun interface{}, params ...interface{}) {
	g.EnqueueWithTask(newTask(taskFun, params...))
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
func (g GoroutineSingle) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

// Schedule execute the task once at the specified time
// that depends on parameters both firstInMs and regularInMs.
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

func (g *GoroutineSingle) dequeueAll() ([]Task, bool) {
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
