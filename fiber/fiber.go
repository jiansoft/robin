package fiber

import (
	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/core"
)

type Fiber interface {
	Start()
	Stop()
	Dispose()
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task core.Task)
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d robin.Disposable)
}
