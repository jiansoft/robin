package robin

import (
	"sync"
	"time"
)

// Worker is an interface for types that can schedule task execution.
type Worker interface {
	Do(taskFunc any, params ...any) Disposable
}

// intervalUnit represents a time unit as a multiplier in milliseconds.
type intervalUnit int64

// Time unit constants, expressed as millisecond multipliers.
const (
	millisecond intervalUnit = 1
	second                   = 1000 * millisecond
	minute                   = 60 * second
	hour                     = 60 * minute
	day                      = 24 * hour
	week                     = 7 * day
)

// jobModel represents the scheduling mode of a Job.
type jobModel int

// Job scheduling modes.
const (
	jobDelay jobModel = iota // one-shot delayed execution
	jobEvery                 // recurring execution at fixed intervals
	jobUntil                 // one-shot execution at a specific time
)

// fiber is the global shared fiber used by default cron scheduling and Channel.
// defaultCron is the package-level CronScheduler that backs all global convenience functions.
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

// newJob creates a new Job bound to this CronScheduler's fiber with default settings.
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

// week sets the job to execute weekly on the specified day.
func (j *Job) week(dayOfWeek time.Weekday) *Job {
	j.intervalUnit = week
	j.weekday = dayOfWeek
	j.interval = 1
	return j
}

// Days sets the job's interval unit to days.
func (j *Job) Days() *Job {
	j.intervalUnit = day
	return j
}

// Hours sets the job's interval unit to hours.
func (j *Job) Hours() *Job {
	j.intervalUnit = hour
	return j
}

// Minutes sets the job's interval unit to minutes.
func (j *Job) Minutes() *Job {
	j.intervalUnit = minute
	return j
}

// Seconds sets the job's interval unit to seconds.
func (j *Job) Seconds() *Job {
	j.intervalUnit = second
	return j
}

// Milliseconds sets the job's interval unit to milliseconds.
func (j *Job) Milliseconds() *Job {
	j.intervalUnit = millisecond
	return j
}

// At sets the specific time of day (hour, minute, second) for job execution.
// Values are clamped: hh mod 24, mm mod 60, ss mod 60. Negative values are made absolute.
func (j *Job) At(hh int, mm int, ss int) *Job {
	j.atHour = Abs(hh) % 24
	j.atMinute = Abs(mm) % 60
	j.atSecond = Abs(ss) % 60
	return j
}

// AfterExecuteTask configures the job to calculate its next execution time after the task finishes.
// This ensures no overlap between consecutive runs. Only applies to delay, second, and millisecond intervals.
func (j *Job) AfterExecuteTask() *Job {
	if j.jobModel == jobDelay || j.intervalUnit == second || j.intervalUnit == millisecond {
		j.afterCalculate = true
	}
	return j
}

// BeforeExecuteTask configures the job to calculate its next execution time immediately,
// without waiting for the current task to finish. This is the default behavior.
func (j *Job) BeforeExecuteTask() *Job {
	j.afterCalculate = false
	return j
}

// Times sets the maximum number of times the job will execute.
// Use -1 (default) for unlimited executions.
func (j *Job) Times(times int64) *Job {
	j.maximumTimes = times
	return j
}

// Between restricts job execution to a daily time window from f to t (only the time-of-day portion is used).
// Has no effect on delay-mode jobs or if f/t is zero or t <= f.
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

// Do registers the task function and starts scheduling the job.
// It returns a Disposable that can be used to cancel the job.
//
// 設計要點：
//   - Do 是 builder chain 的終端方法（Every(1).Seconds().At(10,0,0).Do(fn)），
//     呼叫後 Job 才真正開始排程，之前的 builder 方法只是設定參數。
//   - nextTime 的計算策略是「錨定到最近的符合條件時間點」：
//     先用 At() 指定的時分秒（未指定的欄位填入 now 的值）組合出今天/這小時/這分鐘的目標時刻，
//     再由底部的 !nextTime.After(now) 判斷是否已過期，過期則順延一個 duration。
//     這確保首次觸發總是在「未來最近的對齊點」，而非從 now 開始偏移。
//   - jobUntil 若目標時間已過，直接返回不排程（返回的 Job 本身就是 Disposable，
//     呼叫 Dispose 是安全的 no-op）。這避免了對過期時間排負數 timer。
//   - week 的 weekday 差值用 (target - now + 7) % 7，加 7 防止負數取模，
//     結果 0 表示「就是今天」，交由底部 !nextTime.After(now) 處理是否已過今日時間。
//   - second/millisecond 級別搭配 Between 時間窗時，若 now 在 fromTime 之前（例如凌晨排程
//     但窗口是 09:00-17:00），nextTime 跳到 fromTime；若連 fromTime 都已超過 toTime
//     （代表今天的窗口已結束），整個窗口順延 24 小時到明天。
func (j *Job) Do(fun any, params ...any) Disposable {
	j.task = newTask(fun, params...)
	j.duration = time.Duration(j.interval*int64(j.intervalUnit)) * time.Millisecond
	now := time.Now()

	switch j.jobModel {
	case jobDelay:
		// Delay/RightNow: nextTime = now，實際延遲由 duration 控制
		// （Delay(500) → duration=500ms，schedule 排 500ms 後的 timer）。
		j.nextTime = now
	case jobUntil:
		// 目標時間已過 → 不排程，直接返回 Job 作為惰性 Disposable。
		if j.nextTime.Before(now) {
			return j
		}
	case jobEvery:
		// 未透過 At() 指定的欄位，填入 now 的對應值，
		// 使得 Every(2).Hours().Do(fn) 首次觸發在「當前分秒的下一個 2 小時週期」。
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
			// (target - now + 7) % 7: 加 7 確保非負；結果 0 表示今天就是目標星期。
			i := (j.weekday - now.Weekday() + 7) % 7
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
			// Between 時間窗處理：now 在 fromTime 之前 → 跳到 fromTime 等待窗口開啟。
			if !j.fromTime.IsZero() && j.nextTime.Before(j.fromTime) {
				j.nextTime = j.fromTime
				// 若 fromTime 已超過 toTime，代表今日窗口已過，整體順延到明天。
				if !j.toTime.IsZero() && j.nextTime.After(j.toTime) {
					j.fromTime = j.fromTime.Add(24 * time.Hour)
					j.toTime = j.toTime.Add(24 * time.Hour)
					j.nextTime = j.fromTime
				}
			}
		}
	}

	// 若算出的 nextTime 已經不在未來（今天的 At 時間已過、當前分鐘的秒數已過等），
	// 順延一個 duration，確保首次 timer 一定是正值，不會立即觸發歷史時間點。
	if !j.nextTime.After(now) {
		j.nextTime = j.nextTime.Add(j.duration)
	}

	j.schedule()
	return j
}

// run is the scheduled callback invoked by fiber's timer.
//
// 設計要點：
//   - run 本身不是週期性 ticker，而是「一次性 timer → 執行 → 再排下一次 timer」的遞迴排程。
//     這樣做的好處是每次都能根據實際執行耗時重新校準下次時間，避免 ticker 的漂移累積。
//   - remainTime < 0 表示「已到期或過期」，才進入執行邏輯；>= 0 表示 timer 提前喚醒
//     （OS 排程抖動），此時直接重新排程補償差值，不執行任務。
//   - afterCalculate 模式下必須先釋放鎖再同步執行 task，因為 task 可能耗時長，
//     持鎖會阻塞 Dispose 等外部操作。執行完後重新上鎖並再檢查 disposed，
//     防止在無鎖期間被外部 Dispose 的情況下繼續排程。
//   - Between 時間窗的判斷分成兩層：canRun 決定「現在能不能跑」，
//     窗口過期檢查（nextTime > toTime）決定「下一次該跳到明天的 fromTime」。
//     兩者獨立運作，確保窗口外不執行但仍持續排程。
func (j *Job) run() {
	j.mu.Lock()
	if j.disposed {
		j.mu.Unlock()
		return
	}

	adjustTime := j.remainTime()
	if adjustTime < 0 {
		// 判斷是否在 Between 時間窗內：
		// 無設定 fromTime/toTime → 無限制，可執行；
		// 有設定 toTime 且尚未超過 → 在窗口內，可執行。
		canRun := (!j.toTime.IsZero() && !time.Now().After(j.toTime)) ||
			j.fromTime.IsZero() || j.toTime.IsZero()

		if canRun {
			if j.afterCalculate {
				// AfterExecuteTask 模式：同步執行，確保「上一次跑完才算下一次間隔」。
				// 必須釋放鎖，否則長時間任務會阻塞 Dispose。
				t := j.task
				j.mu.Unlock()
				s := time.Now()
				t.execute()
				elapsed := time.Since(s)
				j.mu.Lock()
				// 無鎖期間可能被 Dispose，必須再檢查。
				if j.disposed {
					j.mu.Unlock()
					return
				}
				// 將執行耗時加入 nextTime，使下次間隔從「執行結束」起算。
				j.nextTime = j.nextTime.Add(elapsed)
			} else {
				// 預設模式：非同步執行，間隔從「排程時間」起算，不等任務完成。
				go j.task.execute()
			}

			// 次數控制：maximumTimes == -1 表示無限次（遞減永遠不會到 0）。
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

		// 推進 nextTime 到下一個週期。
		j.nextTime = j.nextTime.Add(j.duration)

		// 若 nextTime 超出 Between 時間窗，將整個窗口順延 24 小時，
		// 並將 nextTime 重設為明天的 fromTime。
		if !j.toTime.IsZero() && j.nextTime.After(j.toTime) {
			j.fromTime = j.fromTime.Add(24 * time.Hour)
			j.toTime = j.toTime.Add(24 * time.Hour)
			j.nextTime = j.fromTime
		}
	}

	// 無論是否執行了任務，都重新排程下一次 run。
	// adjustTime >= 0（提前喚醒）時直接走到這裡，用剩餘差值補償。
	diff := j.remainTime()
	j.taskDisposer = j.fiber.Schedule(diff, j.run)
	j.mu.Unlock()
}

// remainTime returns the milliseconds remaining until the job's next scheduled execution.
// A negative value indicates the job is overdue.
func (j *Job) remainTime() int64 {
	return j.nextTime.Sub(time.Now()).Milliseconds()
}

// schedule starts the job by scheduling the first run via the fiber's timer.
func (j *Job) schedule() {
	j.mu.Lock()
	diff := j.remainTime()
	j.taskDisposer = j.fiber.Schedule(diff, j.run)
	j.mu.Unlock()
}
