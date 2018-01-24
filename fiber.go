package robin

type Fiber interface {
	Start()
	Stop()
	Dispose()
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d Disposable)
}

type executionState int

const (
	created executionState = iota
	running
	stopped
)
