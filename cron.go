package robin

import (
    "sync"
    "time"
)

type Worker interface {
    Do(taskFunc any, params ...any) Disposable
}

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
    jobUntil
)

var (
    fiber = NewGoroutineMulti()
)

// Abs Returns the absolute value of a specified number.
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](v T) T {
    if v < 0 {
        return -v
    }

    return v
}

type UntilJob struct {
    untilTime time.Time
}

// Job store some information for cron use.
type Job struct {
    nextTime     time.Time
    toTime       time.Time
    fromTime     time.Time
    taskDisposer Disposable
    task         task
    duration     time.Duration
    interval     int64
    maximumTimes int64
    weekday      time.Weekday
    atSecond     int
    atMinute     int
    atHour       int
    jobModel
    intervalUnit
    sync.Mutex
    afterCalculate bool
    disposed       bool
}

// newUntilJob Constructors
func newUntilJob() *UntilJob {
    return new(UntilJob)
}

// Until
func Until(time time.Time) Worker {
    j := newUntilJob()
    j.untilTime = time
    return j
}

// Do
func (u *UntilJob) Do(fun any, params ...any) Disposable {
    j := newJob()
    j.jobModel = jobUntil
    j.nextTime = u.untilTime
    j.maximumTimes = 1
    return j.Do(fun, params...)
}

// RightNow The job executes immediately.
func RightNow() *Job {
    return Delay(0)
}

// Delay The job executes will delay N interval.
func Delay(delayInMs int64) *Job {
    j := newJob()
    j.jobModel = jobDelay
    j.interval = delayInMs
    j.maximumTimes = 1
    j.intervalUnit = millisecond
    return j
}

// EverySunday the job will execute every Sunday .
func EverySunday() *Job {
    return newJob().week(time.Sunday)
}

// EveryMonday the job will execute every Monday
func EveryMonday() *Job {
    return newJob().week(time.Monday)
}

// EveryTuesday the job will execute every Tuesday
func EveryTuesday() *Job {
    return newJob().week(time.Tuesday)
}

// EveryWednesday the job will execute every Wednesday
func EveryWednesday() *Job {
    return newJob().week(time.Wednesday)
}

// EveryThursday the job will execute every Thursday
func EveryThursday() *Job {
    return newJob().week(time.Thursday)
}

// EveryFriday the job will execute every Friday
func EveryFriday() *Job {
    return newJob().week(time.Friday)
}

// EverySaturday the job will execute every Saturday
func EverySaturday() *Job {
    return newJob().week(time.Saturday)
}

// Everyday the job will execute every day
func Everyday() *Job {
    return every(1).Days()
}

// Every the job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func Every(interval int64) *Job {
    return every(interval)
}

// every the job will execute every N everyUnit(ex atHour、atMinute、atSecond、millisecond etc..).
func every(interval int64) *Job {
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
    j.Lock()
    defer j.Unlock()
    if j.disposed {
        return
    }
    j.disposed = true
    j.taskDisposer.Dispose()
}

// week a time interval of execution
func (j *Job) week(dayOfWeek time.Weekday) *Job {
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

// Between the job will be executed only between an assigned period (from f to f time HH:mm:ss.ff).
func (j *Job) Between(f time.Time, t time.Time) *Job {
    if j.jobModel == jobDelay ||
        f.IsZero() ||
        t.IsZero() ||
        t.Unix() <= f.Unix() {
        return j
    }

    now := time.Now()
    y, m, d := now.Year(), now.Month(), now.Day()
    j.fromTime = time.Date(y, m, d, f.Hour(), f.Minute(), f.Second(), f.Nanosecond(), time.Local)
    j.toTime = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)

    return j
}

// Do some job needs to execute.
func (j *Job) Do(fun any, params ...any) Disposable {
    j.task = newTask(fun, params...)
    j.duration = time.Duration(j.interval*int64(j.intervalUnit)) * time.Millisecond
    now := time.Now()

    if j.jobModel == jobDelay {
        j.nextTime = now
    } else if j.jobModel == jobUntil {
        if j.nextTime.UnixNano() < now.UnixNano() {
            return j
        }
    } else if j.jobModel == jobEvery {
        if j.atHour < 0 {
            j.atHour = now.Hour()
        }

        if j.atMinute < 0 {
            j.atMinute = now.Minute()
        }

        if j.atSecond < 0 {
            j.atSecond = now.Second()
        }

        switch j.intervalUnit {
        case week:
            j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, 0, time.Local)
            i := (7 - (now.Weekday() - j.weekday)) % 7
            if i > 0 {
                j.nextTime = j.nextTime.AddDate(0, 0, int(i))
            }
        case day:
            j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), j.atHour, j.atMinute, j.atSecond, 0, time.Local)
        case hour:
            j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), j.atMinute, j.atSecond, 0, time.Local)
        case minute:
            j.nextTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), j.atSecond, 0, time.Local)
        case second, millisecond:
            j.nextTime = now
            if !j.fromTime.IsZero() && j.nextTime.UnixNano() < j.fromTime.UnixNano() {
                j.nextTime = j.fromTime
                if !j.toTime.IsZero() && j.nextTime.UnixNano() > j.toTime.UnixNano() {
                    j.fromTime = j.fromTime.Add(24 * time.Hour)
                    j.toTime = j.toTime.Add(24 * time.Hour)
                    j.nextTime = j.fromTime
                }
            }
        }
    }

    diff := j.nextTime.UnixNano() - now.UnixNano()
    if diff <= 0 {
        j.nextTime = j.nextTime.Add(j.duration)
    }

    j.schedule()
    return j
}

// run the job can be execute or not
func (j *Job) run() {
    adjustTime := j.remainTime()
    if adjustTime < 0 {
        if (!j.toTime.IsZero() && j.toTime.UnixNano() >= time.Now().UnixNano()) || (j.fromTime.IsZero() || j.toTime.IsZero()) {
            if j.afterCalculate {
                s := time.Now()
                fiber.executor.executeTask(j.task)
                d := time.Now().Sub(s)
                j.nextTime = j.nextTime.Add(d)
            } else {
                fiber.executor.executeTaskWithGoroutine(j.task)
            }

            j.maximumTimes += -1
            if j.maximumTimes == 0 {
                j.Dispose()
                return
            }
        }

        j.nextTime = j.nextTime.Add(j.duration)

        if !j.toTime.IsZero() && j.nextTime.UnixNano() > j.toTime.UnixNano() {
            j.fromTime = j.fromTime.Add(24 * time.Hour)
            j.toTime = j.toTime.Add(24 * time.Hour)
            j.nextTime = j.fromTime
        }
    }

    j.schedule()
}

func (j *Job) remainTime() (remainMs int64) {
    var diff = j.nextTime.Sub(time.Now())
    remainMs = int64(diff) / 1e6
    return
}

func (j *Job) schedule() {
    diff := j.remainTime()
    j.Lock()
    j.taskDisposer = fiber.Schedule(diff, j.run)
    j.Unlock()
}
