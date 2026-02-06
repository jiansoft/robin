package robin

import (
	"sync"
	"time"
)

// Worker is an interface for types that can schedule task execution.
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
	fiber       = NewGoroutineMulti()
	defaultCron = &CronScheduler{fiber: fiber}
)

// Abs Returns the absolute value of a specified number.
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](v T) T {
	if v < 0 {
		return -v
	}

	return v
}

// CronScheduler provides an instance-based cron scheduling API.
// Each CronScheduler owns its own fiber for independent lifecycle management.
type CronScheduler struct {
	fiber Fiber
}

// NewCronScheduler creates a new CronScheduler with its own fiber.
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{fiber: NewGoroutineMulti()}
}

// Dispose releases the CronScheduler's fiber resources.
func (cs *CronScheduler) Dispose() {
	cs.fiber.Dispose()
}

// Every creates a job that executes every N units.
func (cs *CronScheduler) Every(interval int64) *Job {
	j := cs.newJob()
	j.interval = interval
	j.intervalUnit = millisecond
	return j
}

// Delay creates a job that executes once after N milliseconds.
func (cs *CronScheduler) Delay(delayInMs int64) *Job {
	j := cs.newJob()
	j.jobModel = jobDelay
	j.interval = delayInMs
	j.maximumTimes = 1
	j.intervalUnit = millisecond
	return j
}

// RightNow creates a job that executes immediately.
func (cs *CronScheduler) RightNow() *Job {
	return cs.Delay(0)
}

// Until creates a job that executes once at the specified time.
func (cs *CronScheduler) Until(t time.Time) Worker {
	return &UntilJob{untilTime: t, fiber: cs.fiber}
}

// Everyday creates a job that executes every day.
func (cs *CronScheduler) Everyday() *Job {
	return cs.Every(1).Days()
}

// EverySunday creates a job that executes every Sunday.
func (cs *CronScheduler) EverySunday() *Job { return cs.newJob().week(time.Sunday) }

// EveryMonday creates a job that executes every Monday.
func (cs *CronScheduler) EveryMonday() *Job { return cs.newJob().week(time.Monday) }

// EveryTuesday creates a job that executes every Tuesday.
func (cs *CronScheduler) EveryTuesday() *Job { return cs.newJob().week(time.Tuesday) }

// EveryWednesday creates a job that executes every Wednesday.
func (cs *CronScheduler) EveryWednesday() *Job { return cs.newJob().week(time.Wednesday) }

// EveryThursday creates a job that executes every Thursday.
func (cs *CronScheduler) EveryThursday() *Job { return cs.newJob().week(time.Thursday) }

// EveryFriday creates a job that executes every Friday.
func (cs *CronScheduler) EveryFriday() *Job { return cs.newJob().week(time.Friday) }

// EverySaturday creates a job that executes every Saturday.
func (cs *CronScheduler) EverySaturday() *Job { return cs.newJob().week(time.Saturday) }

func (cs *CronScheduler) newJob() *Job {
	return &Job{fiber: cs.fiber, jobModel: jobEvery, maximumTimes: -1, atHour: -1, atMinute: -1, atSecond: -1}
}

// --- Global convenience functions (delegate to defaultCron) ---

// RightNow The job executes immediately.
func RightNow() *Job { return defaultCron.RightNow() }

// Delay The job executes will delay N interval.
func Delay(delayInMs int64) *Job { return defaultCron.Delay(delayInMs) }

// Every the job will execute every N everyUnit.
func Every(interval int64) *Job { return defaultCron.Every(interval) }

// Everyday the job will execute every day.
func Everyday() *Job { return defaultCron.Everyday() }

// Until creates a job that executes once at the specified time.
func Until(t time.Time) Worker { return defaultCron.Until(t) }

// EverySunday the job will execute every Sunday.
func EverySunday() *Job { return defaultCron.EverySunday() }

// EveryMonday the job will execute every Monday.
func EveryMonday() *Job { return defaultCron.EveryMonday() }

// EveryTuesday the job will execute every Tuesday.
func EveryTuesday() *Job { return defaultCron.EveryTuesday() }

// EveryWednesday the job will execute every Wednesday.
func EveryWednesday() *Job { return defaultCron.EveryWednesday() }

// EveryThursday the job will execute every Thursday.
func EveryThursday() *Job { return defaultCron.EveryThursday() }

// EveryFriday the job will execute every Friday.
func EveryFriday() *Job { return defaultCron.EveryFriday() }

// EverySaturday the job will execute every Saturday.
func EverySaturday() *Job { return defaultCron.EverySaturday() }

// --- UntilJob ---

// UntilJob schedules a one-shot job at a specific time.
type UntilJob struct {
	untilTime time.Time
	fiber     Fiber
}

// Do schedules the task to execute at the until time.
func (u *UntilJob) Do(fun any, params ...any) Disposable {
	j := &Job{fiber: u.fiber, jobModel: jobUntil, maximumTimes: 1, atHour: -1, atMinute: -1, atSecond: -1}
	j.nextTime = u.untilTime
	return j.Do(fun, params...)
}

// --- Job ---

// Job stores scheduling information for cron use.
type Job struct {
	fiber        Fiber
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
	mu             sync.Mutex
	afterCalculate bool
	disposed       bool
}

// Dispose releases the Job's resources.
func (j *Job) Dispose() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.disposed {
		return
	}
	j.disposed = true
	if j.taskDisposer != nil {
		j.taskDisposer.Dispose()
	}
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

// Between the job will be executed only between an assigned period (from f to t time HH:mm:ss.ff).
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
	j.mu.Lock()
	if j.disposed {
		j.mu.Unlock()
		return
	}

	adjustTime := j.remainTime()
	if adjustTime < 0 {
		canRun := (!j.toTime.IsZero() && j.toTime.UnixNano() >= time.Now().UnixNano()) ||
			j.fromTime.IsZero() || j.toTime.IsZero()

		if canRun {
			if j.afterCalculate {
				t := j.task
				j.mu.Unlock()
				s := time.Now()
				t.execute()
				elapsed := time.Since(s)
				j.mu.Lock()
				if j.disposed {
					j.mu.Unlock()
					return
				}
				j.nextTime = j.nextTime.Add(elapsed)
			} else {
				go j.task.execute()
			}

			j.maximumTimes--
			if j.maximumTimes == 0 {
				j.disposed = true
				d := j.taskDisposer
				j.mu.Unlock()
				if d != nil {
					d.Dispose()
				}
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

	diff := j.remainTime()
	j.taskDisposer = j.fiber.Schedule(diff, j.run)
	j.mu.Unlock()
}

func (j *Job) remainTime() (remainMs int64) {
	diff := j.nextTime.Sub(time.Now())
	remainMs = int64(diff) / 1e6
	return
}

func (j *Job) schedule() {
	j.mu.Lock()
	diff := j.remainTime()
	j.taskDisposer = j.fiber.Schedule(diff, j.run)
	j.mu.Unlock()
}
