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
			return &Job{lock: sync.Mutex{}}
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
	everyUnit    unit
	delayUnit    delayUnit
	interval     int64
	nextTime     time.Time
	timingMode   timingAfterOrBeforeExecuteTask
	lock         sync.Mutex
	runTimes     int64
	maximumTimes int64
	disposed     bool
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
	c.setEveryUnit(delay)
	// Delay model default just execute once.
	c.setMaximumTimes(1)
	//c.everyUnit = delay
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

// Everyday The job will execute every day
func Everyday() *Job {
	return ec.Every(1).Days()
}

// Every The job will execute every N everyUnit(ex hour、minute、second、milliseconds etc..).
func Every(interval int64) *Job {
	return ec.Every(interval)
}

// Every The job will execute every N everyUnit(ex hour、minute、second、milliseconds etc..).
func (c *cronEvery) Every(interval int64) *Job {
	return NewJob(interval, c.fiber, delayNone)
}

func newEveryJob(weekday time.Weekday) *Job {
	c := NewJob(1, ec.fiber, delayNone)
	c.setEveryUnit(weeks)
	//c.everyUnit = weeks
	c.weekday = weekday
	return c
}

// return Job Constructors
func NewJob(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	return jobPool.Get().(*Job).init(intervel, fiber, delayUnit)
}

func (j *Job) init(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	j.lock.Lock()
	j.disposed = false
	j.runTimes = 0
	j.maximumTimes = -1
	j.hour = -1
	j.minute = -1
	j.second = -1
	j.fiber = fiber
	j.loc = time.Local
	j.interval = intervel
	j.delayUnit = delayUnit
	j.timingMode = beforeExecuteTask
	j.identifyId = fmt.Sprintf("%p-%p", &j, &fiber)
	j.lock.Unlock()
	return j
}

// Dispose Job's Dispose
func (j *Job) Dispose() {
	// j.lock.Lock()
	// defer j.lock.Unlock()
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
	return j.identifyId
}

// Days sTime everyUnit of execution
func (j *Job) Days() *Job {
	if j.delayUnit == delayNone {
		j.setEveryUnit(days)
		//j.everyUnit = days
	} else {
		j.delayUnit = delayDays
	}
	return j
}

// Hours Time everyUnit of execution
func (j *Job) Hours() *Job {
	if j.delayUnit == delayNone {
		j.setEveryUnit(hours)
		//j.everyUnit = hours
	} else {
		j.delayUnit = delayHours
	}
	return j
}

// Minutes Time everyUnit of execution
func (j *Job) Minutes() *Job {
	if j.delayUnit == delayNone {
		j.setEveryUnit(minutes)
		//j.everyUnit = minutes
	} else {
		j.delayUnit = delayMinutes
	}
	return j
}

// Seconds Time everyUnit of execution
func (j *Job) Seconds() *Job {
	if j.delayUnit == delayNone {
		j.setEveryUnit(seconds)
		//j.everyUnit = seconds
	} else {
		j.delayUnit = delaySeconds
	}
	return j
}

// MilliSeconds Time everyUnit of execution
func (j *Job) MilliSeconds() *Job {
	if j.delayUnit == delayNone {
		j.setEveryUnit(milliseconds)
		//j.everyUnit = milliseconds
	} else {
		j.delayUnit = delayMilliseconds
	}
	return j
}

// At sThe time specified at execution time
func (j *Job) At(hour int, minute int, second int) *Job {
	j.hour = Abs(j.hour)
	j.minute = Abs(j.minute)
	j.second = Abs(j.second)
	if j.getEveryUnit() != hours {
		j.hour = hour % 24
	}
	j.minute = minute % 60
	j.second = second % 60
	return j
}

// AfterExecuteTask Start timing after the Task is executed
func (j *Job) AfterExecuteTask() *Job {
	//if j.delayUnit == delayNone {
	j.timingMode = afterExecuteTask
	//}
	return j
}

// BeforeExecuteTask Start timing before the Task is executed
func (j *Job) BeforeExecuteTask() *Job {
	//if j.delayUnit == delayNone {
	j.timingMode = beforeExecuteTask
	//}
	return j
}

func (j *Job) Times(times int64) *Job {
	//if j.getDelayUnit() == delayNone {
	j.setMaximumTimes(times)
	//}
	return j
}

// setDelayNextTime
func (j *Job) setDelayNextTime(now time.Time) {
	switch j.delayUnit {
	/*case delayWeeks:
	j.nextTime = now.AddDate(0, 0, 7)*/
	case delayDays:
		j.setNextTime(now.AddDate(0, 0, int(j.getInterval())))
		//j.nextTime = now.AddDate(0, 0, int(j.interval))
	case delayHours:
		j.setNextTime(now.Add(time.Duration(j.getInterval()) * time.Hour))
		//j.nextTime = now.Add(time.Duration(j.interval) * time.Hour)
	case delayMinutes:
		j.setNextTime(now.Add(time.Duration(j.getInterval()) * time.Minute))
		//j.nextTime = now.Add(time.Duration(j.interval) * time.Minute)
	case delaySeconds:
		j.setNextTime(now.Add(time.Duration(j.getInterval()) * time.Second))
		//j.nextTime = now.Add(time.Duration(j.interval) * time.Second)
	case delayMilliseconds:
		j.setNextTime(now.Add(time.Duration(j.getInterval()) * time.Millisecond))
		//j.nextTime = now.Add(time.Duration(j.interval) * time.Millisecond)
	}
}

// firstTimeSetWeeksNextTime
func (j *Job) firstTimeSetWeeksNextTime(now time.Time) {
	i := (7 - (int(now.Weekday() - j.weekday))) % 7
	tmp := time.Date(now.Year(), now.Month(), now.Day(), j.hour, j.minute, j.second, 0, j.loc).AddDate(0, 0, int(i))
	j.setNextTime(tmp)
	//j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.hour, j.minute, j.second, 0, j.loc)
	//j.nextTime = j.nextTime.AddDate(0, 0, int(i))
	if j.getNextTime().Before(now) {
		j.setNextTime(j.getNextTime().AddDate(0, 0, 7))
		//j.nextTime = j.nextTime.AddDate(0, 0, 7)
	}
}

// firstTimeSetDaysNextTime
func (j *Job) firstTimeSetDaysNextTime(now time.Time) {
	if j.second < 0 || j.minute < 0 || j.hour < 0 {
		j.nextTime = now.AddDate(0, 0, 1)
		j.second = j.nextTime.Second()
		j.minute = j.nextTime.Minute()
		j.hour = j.nextTime.Hour()
	} else {
		j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.hour, j.minute, j.second, 0, j.loc)
		if j.interval > 1 {
			j.nextTime = j.nextTime.AddDate(0, 0, int(j.interval-1))
		}
		if j.nextTime.Before(now) {
			j.nextTime = j.nextTime.AddDate(0, 0, int(j.interval))
		}
	}
}

// firstTimeSetHoursNextTime
func (j *Job) firstTimeSetHoursNextTime(now time.Time) {
	if j.minute < 0 {
		j.minute = now.Minute()
	}
	if j.second < 0 {
		j.second = now.Second()
	}
	j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), j.minute, j.second, 0, j.loc)
	j.nextTime.Add(time.Duration(j.interval-1) * time.Hour)
	if j.nextTime.Before(now) {
		j.nextTime = j.nextTime.Add(time.Duration(j.interval) * time.Hour /*.Duration(60*60*1000000)*/)
	}
}

// firstTimeSetMinutesNextTime
func (j *Job) firstTimeSetMinutesNextTime(now time.Time) {
	if j.second < 0 {
		j.second = now.Second()
	}
	j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), j.second, 0, j.loc)
	j.nextTime.Add(time.Duration(j.interval-1) * time.Hour)
	if j.nextTime.Before(now) {
		j.nextTime = j.nextTime.Add(time.Duration(j.interval) * time.Minute /*.Duration(60*60*1000000)*/)
	}
}

// Do some job needs to execute.
func (j *Job) Do(fun interface{}, params ...interface{}) Disposable {
	j.setTask(newTask(fun, params...))
	now := time.Now()
	switch j.getEveryUnit() {
	case delay:
		j.setDelayNextTime(now)
	case weeks:
		j.firstTimeSetWeeksNextTime(now)
	case days:
		j.firstTimeSetDaysNextTime(now)
	case hours:
		j.firstTimeSetHoursNextTime(now)
	case minutes:
		j.firstTimeSetMinutesNextTime(now)
	case seconds:
		tmp := now.Add(time.Duration(j.getInterval()) * time.Second)
		j.setNextTime(tmp)
	case milliseconds:
		tmp := now.Add(time.Duration(j.getInterval()) * time.Millisecond)
		j.setNextTime(tmp)
	}

	firstInMs := int64(j.getNextTime().Sub(now) / time.Millisecond)
	j.setTaskDisposer(firstInMs)
	return j
}

// canDo the job can be execute or not
func (j *Job) canDo() {
	diff := int64(time.Now().Sub(j.getNextTime()) / time.Millisecond /*1000000*/)
	if diff >= 0 /*time.Now().After(j.getNextTime())*/ {
		if j.getTimingMode() == beforeExecuteTask {
			j.fiber.EnqueueWithTask(j.getTask())
		} else {
			s := time.Now()
			j.task.run()
			e := time.Now()
			d := e.Sub(s)
			tmp := j.getNextTime().Add(d)
			j.setNextTime(tmp)
		}

		j.addRunTimes()

		if j.getMaximumTimes() > 0 && j.getRunTimes() >= j.getMaximumTimes() /*|| j.getEveryUnit() == delay*/ {
			j.setDisposed(true)
			j.task.release()
			jobPool.Put(j)
			return
		}

		switch j.getEveryUnit() {
		case delay:
			j.setDelayNextTime(time.Now())
		case weeks:
			tmp := j.getNextTime().AddDate(0, 0, 7)
			j.setNextTime(tmp)
		case days:
			tmp := j.getNextTime().AddDate(0, 0, int(j.getInterval()))
			j.setNextTime(tmp)
		case hours:
			tmp := j.getNextTime().Add(time.Duration(j.getInterval()) * time.Hour)
			j.setNextTime(tmp)
		case minutes:
			tmp := j.getNextTime().Add(time.Duration(j.getInterval()) * time.Minute)
			j.setNextTime(tmp)
		case seconds:
			tmp := j.getNextTime().Add(time.Duration(j.getInterval()) * time.Second)
			j.setNextTime(tmp)
		case milliseconds:
			tmp := j.getNextTime().Add(time.Duration(j.getInterval()) * time.Millisecond)
			j.setNextTime(tmp)
		}
	}

	adjustTime := int64(j.getNextTime().Sub(time.Now()) / time.Millisecond /*1000000*/)
	j.setTaskDisposer(adjustTime)
}

func (j *Job) setTaskDisposer(firstInMs int64) {
	j.lock.Lock()
	j.taskDisposer = j.fiber.Schedule(firstInMs, j.canDo)
	j.lock.Unlock()
}

func (j *Job) setNextTime(nextTime time.Time) {
	j.lock.Lock()
	j.nextTime = nextTime
	j.lock.Unlock()
}

func (j *Job) getNextTime() time.Time {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.nextTime
}

func (j *Job) setEveryUnit(unit unit) {
	j.lock.Lock()
	j.everyUnit = unit
	j.lock.Unlock()
}

func (j *Job) getEveryUnit() unit {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.everyUnit
}
func (j *Job) getDelayUnit() delayUnit {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.delayUnit
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

func (j *Job) setMaximumTimes(r int64) {
	j.lock.Lock()
	j.maximumTimes = r
	j.lock.Unlock()
}

func (j *Job) getMaximumTimes() int64 {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.maximumTimes
}

func (j *Job) getTimingMode() timingAfterOrBeforeExecuteTask {
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
	defer j.lock.Unlock()
	j.task = t
}
