package robin

import (
	"sync"
	"time"
)

type intervalUnit int64

const (
	millisecond intervalUnit = 1
	second                   = 1000 * millisecond
	minute                   = 60 * second
	hour                     = 60 * minute
	day                      = 24 * hour
	week                     = 7 * day
)

type jobModel int

const (
	jobDelay jobModel = iota
	jobEvery
)

var (
	fiber *GoroutineMulti
	dc    = newDelay()
	ec    = newEvery()
)

type cronDelay struct {
}

type cronEvery struct {
}

type Job struct {
	task           Task
	taskDisposer   Disposable
	weekday        time.Weekday
	atHour         int
	atMinute       int
	atSecond       int
	interval       int64
	nextTime       time.Time
	afterCalculate bool
	maximumTimes   int64
	disposed       bool
	duration       time.Duration
	jobModel
	intervalUnit
	sync.Mutex
}

func init() {
	fiber = NewGoroutineMulti()
	fiber.Start()
}

// newDelay Constructors
func newDelay() *cronDelay {
	return new(cronDelay)
}

// newEvery Constructors
func newEvery() *cronEvery {
	return new(cronEvery)
}

// The job executes immediately.
func RightNow() *Job {
	return Delay(0)
}

// The job executes will delay N interval.
func Delay(interval int64) *Job {
	return dc.Delay(interval)
}

// Delay the job executes will delay N interval.
func (c *cronDelay) Delay(interval int64) *Job {
	return newDelayJob(interval)
}

func newDelayJob(delayInMs int64) *Job {
	j := newJob()
	j.jobModel = jobDelay
	j.interval = delayInMs
	j.maximumTimes = 1
	j.intervalUnit = millisecond
	return j
}

// EverySunday the job will execute every Sunday .
func EverySunday() *Job {
	return newJob().everyWeek(time.Sunday)
}

// EveryMonday the job will execute every Monday
func EveryMonday() *Job {
	return newJob().everyWeek(time.Monday)
}

// EveryTuesday the job will execute every Tuesday
func EveryTuesday() *Job {
	return newJob().everyWeek(time.Tuesday)
}

// EveryWednesday the job will execute every Wednesday
func EveryWednesday() *Job {
	return newJob().everyWeek(time.Wednesday)
}

// EveryThursday the job will execute every Thursday
func EveryThursday() *Job {
	return newJob().everyWeek(time.Thursday)
}

// EveryFriday the job will execute every Friday
func EveryFriday() *Job {
	return newJob().everyWeek(time.Friday)
}

// EverySaturday the job will execute every Saturday
func EverySaturday() *Job {
	return newJob().everyWeek(time.Saturday)
}

// Everyday the job will execute every day
func Everyday() *Job {
	return ec.every(1).Days()
}

// every the job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func Every(interval int64) *Job {
	return ec.every(interval)
}

// every the job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func (c *cronEvery) every(interval int64) *Job {
	j := newJob()
	j.interval = interval
	j.intervalUnit = millisecond
	return j
}

// newJob create a Job struct and return it
func newJob() *Job {
	j := &Job{jobModel: jobEvery, maximumTimes: -1, atHour: -1, atMinute: -1, atSecond: -1}
	return j
}

// Dispose Job's Dispose
func (j *Job) Dispose() {
	if j.getDisposed() {
		return
	}
	j.setDisposed(true)
	j.taskDisposer.Dispose()
}

// everyWeek a time interval of execution
func (j *Job) everyWeek(dayOfWeek time.Weekday) *Job {
	j.intervalUnit = week
	j.weekday = dayOfWeek
	j.interval = 1
	return j
}

// Days a time interval of execution
func (j *Job) Days() *Job {
	j.intervalUnit = day
	return j
}

// Hours a time interval of execution
func (j *Job) Hours() *Job {
	j.intervalUnit = hour
	return j
}

// Minutes a time interval of execution
func (j *Job) Minutes() *Job {
	j.intervalUnit = minute
	return j
}

// Seconds a time interval of execution
func (j *Job) Seconds() *Job {
	j.intervalUnit = second
	return j
}

// Milliseconds a time interval of execution
func (j *Job) Milliseconds() *Job {
	j.intervalUnit = millisecond
	return j
}

// At the time specified at execution time
func (j *Job) At(hh int, mm int, ss int) *Job {
	j.atHour = Abs(hh) % 24
	j.atMinute = Abs(mm) % 60
	j.atSecond = Abs(ss) % 60
	return j
}

// AfterExecuteTask waiting for the job execute finish then calculating the job next execution time
// just for delay model、every N second and every N millisecond
// If you want some job every N minute、hour or day do once and want to calculate next execution time by after the job executed.
// Please use interval unit that Seconds or Milliseconds
func (j *Job) AfterExecuteTask() *Job {
	if j.jobModel == jobDelay || j.intervalUnit == second || j.intervalUnit == millisecond {
		j.afterCalculate = true
	}
	return j
}

// BeforeExecuteTask to calculate next execution time immediately don't wait
func (j *Job) BeforeExecuteTask() *Job {
	j.afterCalculate = false
	return j
}

// Times set the job maximum number of executed times
func (j *Job) Times(times int64) *Job {
	j.maximumTimes = times
	return j
}

// Do some job needs to execute.
func (j *Job) Do(fun interface{}, params ...interface{}) Disposable {
	j.task = newTask(fun, params...)
	j.duration = time.Duration(j.interval*int64(j.intervalUnit)) * time.Millisecond
	now := time.Now()
	if j.jobModel == jobDelay || j.intervalUnit == second || j.intervalUnit == millisecond {
		j.nextTime = now
	} else {
		switch j.checkAtTime(now).intervalUnit {
		case week:
			j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, now.Nanosecond(), time.Local)
			i := (7 - (int(now.Weekday() - j.weekday))) % 7
			if i > 0 {
				j.nextTime = j.nextTime.AddDate(0, 0, int(i))
			}
		case day:
			j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, now.Nanosecond(), time.Local)
		case hour:
			j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), j.atMinute, j.atSecond, now.Nanosecond(), time.Local)
		case minute:
			j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), j.atSecond, now.Nanosecond(), time.Local)
		}
	}

	diff := j.nextTime.UnixNano() - now.UnixNano()
	if diff <= 0 {
		j.nextTime = j.nextTime.Add(j.duration)
	}

	firstInMs := j.nextTime.Sub(now).Nanoseconds() / time.Millisecond.Nanoseconds()
	j.schedule(firstInMs)
	return j
}

// canDo the job can be execute or not
func (j *Job) canDo() {
	adjustTime := j.nextTime.Sub(time.Now()).Nanoseconds() / time.Millisecond.Nanoseconds()
	if adjustTime <= 0 {
		if j.afterCalculate {
			s := time.Now()
			fiber.executor.ExecuteTask(j.task)
			d := time.Now().Sub(s)
			j.nextTime = j.nextTime.Add(d)
		} else {
			fiber.executor.ExecuteTaskWithGoroutine(j.task)
		}

		j.maximumTimes += -1
		if j.maximumTimes == 0 {
			j.Dispose()
			return
		}

		j.nextTime = j.nextTime.Add(j.duration)
		adjustTime = j.nextTime.Sub(time.Now()).Nanoseconds() / time.Millisecond.Nanoseconds()
	}

	j.schedule(adjustTime)
}

func (j *Job) schedule(firstInMs int64) {
	j.Lock()
	j.taskDisposer = fiber.Schedule(firstInMs, j.canDo)
	j.Unlock()
}

func (j *Job) getDisposed() bool {
	j.Lock()
	defer j.Unlock()
	return j.disposed
}

func (j *Job) setDisposed(r bool) {
	j.Lock()
	j.disposed = r
	j.Unlock()
}

func (j *Job) checkAtTime(now time.Time) *Job {
	if j.atHour < 0 {
		j.atHour = now.Hour()
	}

	if j.atMinute < 0 {
		j.atMinute = now.Minute()
	}

	if j.atSecond < 0 {
		j.atSecond = now.Second()
	}
	return j
}
