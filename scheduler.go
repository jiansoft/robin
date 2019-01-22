package robin

type IScheduler interface {
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	Start()
	Stop()
	Dispose()
}

type SchedulerRegistry interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
	Remove(d Disposable)
}

//Allows for the registration and deregistration of subscriptions /*The IFiber has implemented*/
type SubscriptionRegistry interface {
	//Register subscription to be unsubcribed from when the scheduler is disposed.
	RegisterSubscription(Disposable)
	//Deregister a subscription.
	DeregisterSubscription(Disposable)
}

type ExecutionContext interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
}

type Scheduler struct {
	fiber       ExecutionContext
	running     bool
	isDispose   bool
	disposabler *Disposer
}

func (s *Scheduler) init(executionState ExecutionContext) *Scheduler {
	s.fiber = executionState
	s.running = true
	s.disposabler = NewDisposer()
	return s
}

func NewScheduler(executionState ExecutionContext) *Scheduler {
	return new(Scheduler).init(executionState)
}

// Schedule delay n millisecond then execute once the function
func (s *Scheduler) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	return s.ScheduleOnInterval(firstInMs, -1, taskFun, params...)
}

// ScheduleOnInterval first time delay N millisecond then execute once the function,
// then interval N millisecond repeat execute the function.
func (s *Scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable) {
	pending := newTimerTask(s, newTask(taskFun, params...), firstInMs, regularInMs)
	if s.isDispose {
		return pending
	}
	pending.schedule()
	if firstInMs > 0 {
		s.disposabler.Add(pending)
	}
	/*if firstInMs <= 0 {
		//pending.executeOnFiber()
		//pending.Dispose()
	} else {
		//pending.schedule()
		s.disposabler.Add(pending)
	}*/
	return pending
}

//Implement SchedulerRegistry.Enqueue
func (s *Scheduler) Enqueue(taskFun interface{}, params ...interface{}) {
	s.EnqueueWithTask(newTask(taskFun, params...))
}

func (s *Scheduler) EnqueueWithTask(task Task) {
	s.fiber.EnqueueWithTask(task)
}

//Implement SchedulerRegistry.Remove
func (s *Scheduler) Remove(d Disposable) {
	s.disposabler.Remove(d)
	//s.fiber.Enqueue(s.disposabler.Remove, d)
}

func (s *Scheduler) Start() {
	s.running = true
	s.isDispose = false
}

func (s *Scheduler) Stop() {
	s.running = false
}

func (s *Scheduler) Dispose() {
	s.Stop()
	s.isDispose = true
	s.disposabler.Dispose()
}
