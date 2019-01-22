package robin

import (
	"fmt"
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

type afterOrBeforeExecuteTask int

const (
	beforeExecuteTask afterOrBeforeExecuteTask = iota
	afterExecuteTask
)

type jobModel int

const (
	JobDelayModel jobModel = iota
	JobEveryModel
)

var (
	jobPool = sync.Pool{
		New: func() interface{} {
			return &Job{lock: sync.Mutex{}}
		},
	}
	dc = newDelay()
	ec = newEvery()
)

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
	atHour       int
	atMinute     int
	atSecond     int
	interval     int64
	nextTime     time.Time
	timingMode   afterOrBeforeExecuteTask
	lock         sync.Mutex
	runTimes     int64
	maximumTimes int64
	disposed     bool
	jobModel
	intervalUnit
}

// The job executes immediately.
func RightNow() *Job {
	return Delay(0)
}

// The job executes will delay N Milliseconds.
func Delay(delayInMs int64) *Job {
	return dc.Delay(delayInMs)
}

// newDelay Constructors
func newDelay() *cronDelay {
	return new(cronDelay).init()
}

func newDelayJob(delayInMs int64) *Job {
	return NewJob(dc.fiber).setInterval(delayInMs).Times(1).setJobModel(JobDelayModel).MilliSeconds()
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
func newEvery() *cronEvery {
	return new(cronEvery).init()
}

func (c *cronEvery) init() *cronEvery {
	c.fiber = NewGoroutineMulti()
	c.fiber.Start()
	return c
}

// EverySunday The job will execute every Sunday .
func EverySunday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Sunday)
}

// EveryMonday The job will execute every Monday
func EveryMonday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Monday)
}

// EveryTuesday The job will execute every Tuesday
func EveryTuesday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Tuesday)
}

// EveryWednesday The job will execute every Wednesday
func EveryWednesday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Wednesday)
}

// EveryThursday The job will execute every Thursday
func EveryThursday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Thursday)
}

// EveryFriday The job will execute every Friday
func EveryFriday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Friday)
}

// EverySaturday The job will execute every Saturday
func EverySaturday() *Job {
	return NewJob(ec.fiber).setInterval(1).Week(time.Saturday)
}

// Everyday The job will execute every day
func Everyday() *Job {
	return ec.Every(1).Days()
}

// Every The job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func Every(interval int64) *Job {
	return ec.Every(interval)
}

// Every The job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func (c *cronEvery) Every(interval int64) *Job {
	return NewJob(c.fiber).setInterval(interval)
}

// return Job Constructors
func NewJob(fiber Fiber) *Job {
	j := jobPool.Get().(*Job)
	j.lock.Lock()
	j.jobModel = JobEveryModel
	j.disposed = false
	j.runTimes = 0
	j.maximumTimes = -1
	j.atHour = -1
	j.atMinute = -1
	j.atSecond = -1
	j.fiber = fiber
	j.loc = time.Local
	j.timingMode = beforeExecuteTask
	j.identifyId = fmt.Sprintf("%p-%p", &j, &fiber)
	j.lock.Unlock()
	return j
}

// Dispose Job's Dispose
func (j *Job) Dispose() {
	if j.getDisposed() {
		return
	}
	j.setDisposed(true)
	j.taskDisposer.Dispose()
	j.task.release()
	jobPool.Put(j)
}

// Identify Job's Identify
func (j Job) Identify() string {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.identifyId
}

// Week a time interval of execution
func (j *Job) Week(dayOfWeek time.Weekday) *Job {
	return j.setIntervalUnit(week).setWeekday(dayOfWeek)
}

// Days a time interval of execution
func (j *Job) Days() *Job {
	return j.setIntervalUnit(day)
}

// Hours a time interval of execution
func (j *Job) Hours() *Job {
	return j.setIntervalUnit(hour)
}

// Minutes a time interval of execution
func (j *Job) Minutes() *Job {
	return j.setIntervalUnit(minute)
}

// Seconds a time interval of execution
func (j *Job) Seconds() *Job {
	return j.setIntervalUnit(second)
}

// MilliSeconds a time interval of execution
func (j *Job) MilliSeconds() *Job {
	return j.setIntervalUnit(millisecond)
}

// At the time specified at execution time
func (j *Job) At(hh int, mm int, ss int) *Job {
	j.lock.Lock()
	j.atHour = Abs(hh) % 24
	j.atMinute = Abs(mm) % 60
	j.atSecond = Abs(ss) % 60
	j.lock.Unlock()
	return j
}

// AfterExecuteTask Start timing after the Task is executed
func (j *Job) AfterExecuteTask() *Job {
	return j.setTimingMode(afterExecuteTask)
}

// BeforeExecuteTask Start timing before the Task is executed
func (j *Job) BeforeExecuteTask() *Job {
	return j.setTimingMode(beforeExecuteTask)
}

// Times set the job maximum number of executed times
func (j *Job) Times(times int64) *Job {
	return j.setMaximumTimes(times)
}

// Do some job needs to execute.
func (j *Job) Do(fun interface{}, params ...interface{}) Disposable {
	j.setTask(newTask(fun, params...))
	now := time.Now()
	if j.getJobModel() == JobDelayModel {
		nextTime := now.Add(time.Duration(j.getInterval()*int64(j.getIntervalUnit())) * time.Millisecond)
		j.setNextTime(nextTime)
	} else {
		switch j.checkAtTime().getIntervalUnit() {
		case week:
			i := (7 - (int(now.Weekday() - j.weekday))) % 7
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, 0, j.loc).AddDate(0, 0, int(i))
			j.setNextTime(nextTime).checkNextTime(now)
		case day:
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, 0, j.loc)
			j.setNextTime(nextTime).checkNextTime(now)
		case hour:
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), j.atMinute, j.atSecond, 0, j.loc)
			j.setNextTime(nextTime).checkNextTime(now)
		case minute:
			nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), j.atSecond, 0, j.loc)
			j.setNextTime(nextTime).checkNextTime(now)
		default:
			nextTime := now.Add(time.Duration(j.getInterval()*int64(j.getIntervalUnit())) * time.Millisecond)
			j.setNextTime(nextTime)
		}
	}
	firstInMs := int64(j.getNextTime().Sub(now) / time.Millisecond)
	j.setTaskDisposer(firstInMs)
	return j
}

// canDo the job can be execute or not
func (j *Job) canDo() {
	diff := int64(time.Now().Sub(j.getNextTime()) / time.Millisecond /*1000000*/)
	if diff >= 0 {
		if j.getTimingMode() == beforeExecuteTask {
			j.fiber.EnqueueWithTask(j.getTask())
		} else {
			s := time.Now()
			j.task.run()
			e := time.Now()
			d := e.Sub(s)
			nextTime := j.getNextTime().Add(d)
			j.setNextTime(nextTime)
		}

		j.addRunTimes()

		if j.getMaximumTimes() > 0 && j.getRunTimes() >= j.getMaximumTimes() {
			j.setDisposed(true)
			j.task.release()
			jobPool.Put(j)
			return
		}

		nextTime := j.getNextTime().Add(time.Duration(j.getInterval()*int64(j.getIntervalUnit())) * time.Millisecond)
		j.setNextTime(nextTime)
	}

	adjustTime := int64(j.getNextTime().Sub(time.Now()) / time.Millisecond /*1000000*/)
	j.setTaskDisposer(adjustTime)
}

func (j *Job) setJobModel(jobModel jobModel) *Job {
	j.lock.Lock()
	j.jobModel = jobModel
	j.lock.Unlock()
	return j
}

func (j *Job) setTaskDisposer(firstInMs int64) {
	j.lock.Lock()
	j.taskDisposer = j.fiber.Schedule(firstInMs, j.canDo)
	j.lock.Unlock()
}

func (j *Job) setNextTime(nextTime time.Time) *Job {
	j.lock.Lock()
	j.nextTime = nextTime
	j.lock.Unlock()
	return j
}

func (j *Job) getNextTime() time.Time {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.nextTime
}

func (j *Job) setIntervalUnit(intervalUnit intervalUnit) *Job {
	j.lock.Lock()
	defer j.lock.Unlock()
	j.intervalUnit = intervalUnit
	return j
}

func (j *Job) getIntervalUnit() intervalUnit {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.intervalUnit
}

func (j *Job) getJobModel() jobModel {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.jobModel
}

func (j *Job) setInterval(interval int64) *Job {
	j.lock.Lock()
	defer j.lock.Unlock()
	j.interval = interval
	return j
}

func (j *Job) getInterval() int64 {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.interval
}

func (j *Job) addRunTimes() {
	j.lock.Lock()
	j.runTimes++
	j.lock.Unlock()
}

func (j *Job) getRunTimes() int64 {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.runTimes
}

func (j *Job) getDisposed() bool {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.disposed
}

func (j *Job) setDisposed(r bool) {
	j.lock.Lock()
	j.disposed = r
	j.lock.Unlock()
}

func (j *Job) setMaximumTimes(r int64) *Job {
	j.lock.Lock()
	j.maximumTimes = r
	j.lock.Unlock()
	return j
}

func (j *Job) getMaximumTimes() int64 {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.maximumTimes
}

func (j *Job) setTimingMode(timingMode afterOrBeforeExecuteTask) *Job {
	j.lock.Lock()
	j.timingMode = timingMode
	j.lock.Unlock()
	return j
}

func (j *Job) getTimingMode() afterOrBeforeExecuteTask {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.timingMode
}

func (j *Job) getTask() Task {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.task
}

func (j *Job) setTask(t Task) {
	j.lock.Lock()
	j.task = t
	j.lock.Unlock()
}

func (j *Job) setWeekday(weekday time.Weekday) *Job {
	j.lock.Lock()
	j.weekday = weekday
	j.lock.Unlock()
	return j
}

func (j *Job) checkAtTime() *Job {
	now := time.Now()
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

func (j *Job) checkNextTime(now time.Time) *Job {
	if j.getInterval() > 1 {
		nextTime := j.getNextTime().Add(time.Duration(j.getInterval()*int64(j.getIntervalUnit())) * time.Millisecond)
		j.setNextTime(nextTime)
	}

	if j.getNextTime().Before(now) {
		nextTime := j.getNextTime().Add(time.Duration(j.getInterval()*int64(j.getIntervalUnit())) * time.Millisecond)
		j.setNextTime(nextTime)
	}

	return j
}
