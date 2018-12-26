package robin

import (
	"fmt"
	"sync"
	"time"
)

const (
	delay unit = iota
	weeks
	days
	hours
	minutes
	seconds
	milliseconds
)

const (
	beforeExecuteTask timingAfterOrBeforeExecuteTask = iota
	afterExecuteTask
)

const (
	delayNone delayUnit = iota
	//delayWeeks
	delayDays
	delayHours
	delayMinutes
	delaySeconds
	delayMilliseconds
)

var (
	jobPool = sync.Pool{
		New: func() interface{} {
			return new(Job)
		},
	}
	dc = newCronDelay()
	ec = newEveryCron()
)

type unit int
type delayUnit int
type timingAfterOrBeforeExecuteTask int

type cronDelay struct {
	fiber Fiber
}

type cronEvery struct {
	fiber Fiber
}

type Job struct {
	fiber        Fiber
	identifyId   string
	loc          *time.Location
	task         Task
	taskDisposer Disposable
	weekday      time.Weekday
	hour         int
	minute       int
	second       int
	unit         unit
	delayUnit    delayUnit
	interval     int64
	nextTime     time.Time
	timingMode   timingAfterOrBeforeExecuteTask
}

// The job executes immediately.
func RightNow() *Job {
	return Delay(0)
}

// The job executes will delay N Milliseconds.
func Delay(delayInMs int64) *Job {
	return dc.Delay(delayInMs)
}

// newCronDelay Constructors
func newCronDelay() *cronDelay {
	return new(cronDelay).init()
}

func newDelayJob(delayInMs int64) *Job {
	c := NewJob(delayInMs, dc.fiber, delayMilliseconds)
	c.unit = delay
	return c
}

func (c *cronDelay) init() *cronDelay {
	c.fiber = NewGoroutineMulti()
	c.fiber.Start()
	return c
}

// The job executes will delay N Milliseconds.
func (c *cronDelay) Delay(delayInMs int64) *Job {
	return newDelayJob(delayInMs)
}

// EveryCron Constructors
func newEveryCron() *cronEvery {
	return new(cronEvery).init()
}

func (c *cronEvery) init() *cronEvery {
	c.fiber = NewGoroutineMulti()
	c.fiber.Start()
	return c
}

// EverySunday The job will execute every Sunday .
func EverySunday() *Job {
	return newEveryJob(time.Sunday)
}

// EveryMonday The job will execute every Monday
func EveryMonday() *Job {
	return newEveryJob(time.Monday)
}

// EveryTuesday The job will execute every Tuesday
func EveryTuesday() *Job {
	return newEveryJob(time.Tuesday)
}

// EveryWednesday The job will execute every Wednesday
func EveryWednesday() *Job {
	return newEveryJob(time.Wednesday)
}

// EveryThursday The job will execute every Thursday
func EveryThursday() *Job {
	return newEveryJob(time.Thursday)
}

// EveryFriday The job will execute every Friday
func EveryFriday() *Job {
	return newEveryJob(time.Friday)
}

// EverySaturday The job will execute every Saturday
func EverySaturday() *Job {
	return newEveryJob(time.Saturday)
}

// Every The job will execute every N unit(ex hour、minute、second、milliseconds etc..).
func Every(interval int64) *Job {
	return ec.Every(interval)
}

// Every The job will execute every N unit(ex hour、minute、second、milliseconds etc..).
func (c *cronEvery) Every(interval int64) *Job {
	return NewJob(interval, c.fiber, delayNone)
}

func newEveryJob(weekday time.Weekday) *Job {
	c := NewJob(1, ec.fiber, delayNone)
	c.unit = weeks
	c.weekday = weekday
	return c
}

// return Job Constructors
func NewJob(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	return jobPool.Get().(*Job).init(intervel, fiber, delayUnit)
}

func (c *Job) init(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	c.hour = -1
	c.minute = -1
	c.second = -1
	c.fiber = fiber
	c.loc = time.Local
	c.interval = intervel
	c.delayUnit = delayUnit
	c.timingMode = beforeExecuteTask
	c.identifyId = fmt.Sprintf("%p-%p-%d", &c, &fiber, time.Now().UnixNano())
	return c
}

// Dispose Job's Dispose
func (c *Job) Dispose() {
	c.taskDisposer.Dispose()
	jobPool.Put(c)
}

// Identify Job's Identify
func (c Job) Identify() string {
	return c.identifyId
}

// Days sTime unit of execution
func (c *Job) Days() *Job {
	if c.delayUnit == delayNone {
		c.unit = days
	} else {
		c.delayUnit = delayDays
	}
	return c
}

// Hours Time unit of execution
func (c *Job) Hours() *Job {
	if c.delayUnit == delayNone {
		c.unit = hours
	} else {
		c.delayUnit = delayHours
	}
	return c
}

// Minutes Time unit of execution
func (c *Job) Minutes() *Job {
	if c.delayUnit == delayNone {
		c.unit = minutes
	} else {
		c.delayUnit = delayMinutes
	}
	return c
}

// Seconds Time unit of execution
func (c *Job) Seconds() *Job {
	if c.delayUnit == delayNone {
		c.unit = seconds
	} else {
		c.delayUnit = delaySeconds
	}
	return c
}

// MilliSeconds Time unit of execution
func (c *Job) MilliSeconds() *Job {
	if c.delayUnit == delayNone {
		c.unit = milliseconds
	} else {
		c.delayUnit = delayMilliseconds
	}
	return c
}

// At sThe time specified at execution time
func (c *Job) At(hour int, minute int, second int) *Job {
	c.hour = Abs(c.hour)
	c.minute = Abs(c.minute)
	c.second = Abs(c.second)
	if c.unit != hours {
		c.hour = hour % 24
	}
	c.minute = minute % 60
	c.second = second % 60
	return c
}

// AfterExecuteTask Start timing after the Task is executed
func (c *Job) AfterExecuteTask() *Job {
	if c.delayUnit == delayNone {
		c.timingMode = afterExecuteTask
	}
	return c
}

// BeforeExecuteTask Start timing before the Task is executed
func (c *Job) BeforeExecuteTask() *Job {
	if c.delayUnit == delayNone {
		c.timingMode = beforeExecuteTask
	}
	return c
}

// firstTimeSetDelayNextTime
func (c *Job) firstTimeSetDelayNextTime(now time.Time) {
	switch c.delayUnit {
	/*case delayWeeks:
	c.nextTime = now.AddDate(0, 0, 7)*/
	case delayDays:
		c.nextTime = now.AddDate(0, 0, int(c.interval))
	case delayHours:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Hour)
	case delayMinutes:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Minute)
	case delaySeconds:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Second)
	case delayMilliseconds:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Millisecond)
	}
}

// firstTimeSetWeeksNextTime
func (c *Job) firstTimeSetWeeksNextTime(now time.Time) {
	i := (7 - (int(now.Weekday() - c.weekday))) % 7
	c.nextTime = time.Date(now.Year(), now.Month(), now.Day()+int(i), c.hour, c.minute, c.second, 0, c.loc)
	if c.nextTime.Before(now) {
		c.nextTime = c.nextTime.AddDate(0, 0, 7)
	}
}

// firstTimeSetDaysNextTime
func (c *Job) firstTimeSetDaysNextTime(now time.Time) {
	if c.second < 0 || c.minute < 0 || c.hour < 0 {
		c.nextTime = now.AddDate(0, 0, 1)
		c.second = c.nextTime.Second()
		c.minute = c.nextTime.Minute()
		c.hour = c.nextTime.Hour()
	} else {
		c.nextTime = time.Date(now.Year(), now.Month(), now.Day(), c.hour, c.minute, c.second, 0, c.loc)
		if c.interval > 1 {
			c.nextTime = c.nextTime.AddDate(0, 0, int(c.interval-1))
		}
		if c.nextTime.Before(now) {
			c.nextTime = c.nextTime.AddDate(0, 0, int(c.interval))
		}
	}
}

// firstTimeSetHoursNextTime
func (c *Job) firstTimeSetHoursNextTime(now time.Time) {
	if c.minute < 0 {
		c.minute = now.Minute()
	}
	if c.second < 0 {
		c.second = now.Second()
	}
	c.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), c.minute, c.second, 0, c.loc)
	c.nextTime.Add(time.Duration(c.interval-1) * time.Hour)
	if c.nextTime.Before(now) {
		c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Hour /*.Duration(60*60*1000000)*/)
	}
}

// firstTimeSetMinutesNextTime
func (c *Job) firstTimeSetMinutesNextTime(now time.Time) {
	if c.second < 0 {
		c.second = now.Second()
	}
	c.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), c.second, 0, c.loc)
	c.nextTime.Add(time.Duration(c.interval-1) * time.Hour)
	if c.nextTime.Before(now) {
		c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Minute /*.Duration(60*60*1000000)*/)
	}
}

// Do some job needs to execute.
func (c *Job) Do(fun interface{}, params ...interface{}) Disposable {
	c.task = newTask(fun, params...)
	now := time.Now()
	switch c.unit {
	case delay:
		c.firstTimeSetDelayNextTime(now)
	case weeks:
		c.firstTimeSetWeeksNextTime(now)
	case days:
		c.firstTimeSetDaysNextTime(now)
	case hours:
		c.firstTimeSetHoursNextTime(now)
	case minutes:
		c.firstTimeSetMinutesNextTime(now)
	case seconds:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Second)
	case milliseconds:
		c.nextTime = now.Add(time.Duration(c.interval) * time.Millisecond)
	}

	firstInMs := int64(c.nextTime.Sub(now) / time.Millisecond)
	c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
	return c
}

// Is the job can be executed
func (c *Job) canDo() {
	diff := int64(time.Now().Sub(c.nextTime) / time.Millisecond /*1000000*/)
	if diff >= 0 {
		if c.unit == delay || c.timingMode == beforeExecuteTask {
			c.fiber.EnqueueWithTask(c.task)
		} else {
			s := time.Now()
			c.task.run()
			e := time.Now()
			d := e.Sub(s)
			c.nextTime = c.nextTime.Add(d)
		}
		switch c.unit {
		case delay:
			jobPool.Put(c)
			return
		case weeks:
			c.nextTime = c.nextTime.AddDate(0, 0, 7)
		case days:
			c.nextTime = c.nextTime.AddDate(0, 0, int(c.interval))
		case hours:
			c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Hour)
		case minutes:
			c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Minute)
		case seconds:
			c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Second)
		case milliseconds:
			c.nextTime = c.nextTime.Add(time.Duration(c.interval) * time.Millisecond)
		}
	}

	c.taskDisposer.Dispose()
	adjustTime := int64(c.nextTime.Sub(time.Now()) / time.Millisecond /*1000000*/)
	c.taskDisposer = c.fiber.Schedule(adjustTime, c.canDo)
}
