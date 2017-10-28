package robin

var delayExecutor = NewDelaySchedulerExecutor()

type delaySchedulerExecutor struct {
	fiber Fiber
}

func NewDelaySchedulerExecutor() *delaySchedulerExecutor {
	return new(delaySchedulerExecutor).init()
}

func (c *delaySchedulerExecutor) init() *delaySchedulerExecutor {
	c.fiber = NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func (c *delaySchedulerExecutor) Delay(delayInMs int64) *Job {
	return newCronDelay(delayInMs)
}

func Delay(delayInMs int64) *Job {
	return delayExecutor.Delay(delayInMs)
}

func newCronDelay(delayInMs int64) *Job {
	c := NewJob(delayInMs, delayExecutor.fiber)
	c.unit = delay
	return c
}
