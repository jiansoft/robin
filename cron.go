package robin

import (
	"sync"
	"time"
)

// Worker is a scheduling endpoint that can register a task for future execution.
// Worker 是可註冊未來執行任務的排程端點介面。
type Worker interface {
	Do(taskFunc any, params ...any) Disposable
}

// intervalUnit expresses a time unit as a millisecond multiplier.
// intervalUnit 以毫秒倍數表示時間單位。
type intervalUnit int64

// Time unit constants, expressed as millisecond multipliers.
// 時間單位常數，皆以毫秒倍數表示。
const (
	millisecond intervalUnit = 1
	second                   = 1000 * millisecond
	minute                   = 60 * second
	hour                     = 60 * minute
	day                      = 24 * hour
	week                     = 7 * day
)

// jobModel identifies how a job computes and repeats schedule points.
// jobModel 用來表示 Job 的排程模式與重複策略。
type jobModel int

// Job scheduling modes.
// Job 排程模式。
const (
	jobDelay jobModel = iota // one-shot delayed execution
	jobEvery                 // recurring execution at fixed intervals
	jobUntil                 // one-shot execution at a specific time
)

// fiber is the global shared fiber used by default cron scheduling and Channel.
// defaultCron is the package-level CronScheduler that backs all global convenience functions.
// fiber 是預設 cron 與 Channel 共用的全域 fiber。
// defaultCron 是所有全域便利函式背後使用的 CronScheduler。
var (
	fiber       = NewGoroutineMulti()
	defaultCron = &CronScheduler{fiber: fiber}
)

// normalizePositiveInterval clamps interval values to a positive minimum of 1.
// It prevents zero/negative interval values from creating busy-loop schedules.
// normalizePositiveInterval 會把間隔值正規化為至少 1。
// 這可避免零或負值間隔造成忙迴圈排程。
func normalizePositiveInterval(interval int64) int64 {
	if interval <= 0 {
		return 1
	}
	return interval
}

// CronScheduler provides an instance-based cron scheduling API.
// Each CronScheduler owns its own fiber for independent lifecycle management.
// CronScheduler 提供實例化的 cron 排程 API。
// 每個 CronScheduler 都擁有獨立 fiber，可獨立管理生命週期。
type CronScheduler struct {
	fiber Fiber
}

// NewCronScheduler creates an isolated scheduler with its own fiber.
// NewCronScheduler 會建立具有獨立 fiber 的排程器實例。
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{fiber: NewGoroutineMulti()}
}

// Dispose releases all resources owned by this scheduler instance.
// Dispose 會釋放此 scheduler 實例所持有的所有資源。
func (cs *CronScheduler) Dispose() {
	cs.fiber.Dispose()
}

// Every creates a recurring job whose interval defaults to milliseconds.
// Non-positive intervals are normalized to 1 to avoid busy-loop schedules.
// Every 會建立重複型 Job，預設單位為毫秒。
// 非正間隔會正規化為 1，以避免忙迴圈排程。
func (cs *CronScheduler) Every(interval int64) *Job {
	j := cs.newJob()
	j.interval = normalizePositiveInterval(interval)
	j.intervalUnit = millisecond
	return j
}

// Delay creates a one-shot job that runs after delayInMs milliseconds.
// Delay 會建立一次性 Job，在 delayInMs 毫秒後執行。
func (cs *CronScheduler) Delay(delayInMs int64) *Job {
	j := cs.newJob()
	j.jobModel = jobDelay
	j.interval = delayInMs
	j.maximumTimes = 1
	j.intervalUnit = millisecond
	return j
}

// RightNow creates a one-shot job that executes immediately.
// RightNow 會建立立即執行的一次性 Job。
func (cs *CronScheduler) RightNow() *Job {
	return cs.Delay(0)
}

// Until creates a one-shot worker targeting a specific absolute time.
// Until 會建立在指定絕對時間觸發的一次性 Worker。
func (cs *CronScheduler) Until(t time.Time) Worker {
	return &UntilJob{untilTime: t, fiber: cs.fiber}
}

// Everyday creates a daily recurring job.
// Everyday 會建立每日重複執行的 Job。
func (cs *CronScheduler) Everyday() *Job {
	return cs.Every(1).Days()
}

// EverySunday creates a weekly job for Sunday.
// EverySunday 會建立每週日執行的 Job。
func (cs *CronScheduler) EverySunday() *Job { return cs.newJob().week(time.Sunday) }

// EveryMonday creates a weekly job for Monday.
// EveryMonday 會建立每週一執行的 Job。
func (cs *CronScheduler) EveryMonday() *Job { return cs.newJob().week(time.Monday) }

// EveryTuesday creates a weekly job for Tuesday.
// EveryTuesday 會建立每週二執行的 Job。
func (cs *CronScheduler) EveryTuesday() *Job { return cs.newJob().week(time.Tuesday) }

// EveryWednesday creates a weekly job for Wednesday.
// EveryWednesday 會建立每週三執行的 Job。
func (cs *CronScheduler) EveryWednesday() *Job { return cs.newJob().week(time.Wednesday) }

// EveryThursday creates a weekly job for Thursday.
// EveryThursday 會建立每週四執行的 Job。
func (cs *CronScheduler) EveryThursday() *Job { return cs.newJob().week(time.Thursday) }

// EveryFriday creates a weekly job for Friday.
// EveryFriday 會建立每週五執行的 Job。
func (cs *CronScheduler) EveryFriday() *Job { return cs.newJob().week(time.Friday) }

// EverySaturday creates a weekly job for Saturday.
// EverySaturday 會建立每週六執行的 Job。
func (cs *CronScheduler) EverySaturday() *Job { return cs.newJob().week(time.Saturday) }

// newJob builds a Job bound to this scheduler's fiber with default recurring settings.
// newJob 會建立綁定於當前 scheduler fiber 的 Job，並套用預設重複參數。
func (cs *CronScheduler) newJob() *Job {
	return &Job{fiber: cs.fiber, jobModel: jobEvery, maximumTimes: -1, atHour: -1, atMinute: -1, atSecond: -1}
}

// --- Global convenience functions (delegate to defaultCron) ---
// --- 全域便利函式（委派給 defaultCron）---

// RightNow returns a global one-shot job that executes immediately.
// RightNow 回傳使用全域 defaultCron 的立即執行一次性 Job。
func RightNow() *Job { return defaultCron.RightNow() }

// Delay returns a global one-shot job that runs after delayInMs milliseconds.
// Delay 回傳使用全域 defaultCron 的延遲一次性 Job。
func Delay(delayInMs int64) *Job { return defaultCron.Delay(delayInMs) }

// Every returns a global recurring job builder.
// Every 回傳使用全域 defaultCron 的重複 Job 建構器。
func Every(interval int64) *Job { return defaultCron.Every(interval) }

// Everyday returns a global daily recurring job.
// Everyday 回傳使用全域 defaultCron 的每日重複 Job。
func Everyday() *Job { return defaultCron.Everyday() }

// Until returns a global one-shot worker that runs at the specified time.
// Until 回傳使用全域 defaultCron、在指定時間觸發的一次性 Worker。
func Until(t time.Time) Worker { return defaultCron.Until(t) }

// EverySunday returns a global weekly job for Sunday.
// EverySunday 回傳使用全域 defaultCron 的每週日 Job。
func EverySunday() *Job { return defaultCron.EverySunday() }

// EveryMonday returns a global weekly job for Monday.
// EveryMonday 回傳使用全域 defaultCron 的每週一 Job。
func EveryMonday() *Job { return defaultCron.EveryMonday() }

// EveryTuesday returns a global weekly job for Tuesday.
// EveryTuesday 回傳使用全域 defaultCron 的每週二 Job。
func EveryTuesday() *Job { return defaultCron.EveryTuesday() }

// EveryWednesday returns a global weekly job for Wednesday.
// EveryWednesday 回傳使用全域 defaultCron 的每週三 Job。
func EveryWednesday() *Job { return defaultCron.EveryWednesday() }

// EveryThursday returns a global weekly job for Thursday.
// EveryThursday 回傳使用全域 defaultCron 的每週四 Job。
func EveryThursday() *Job { return defaultCron.EveryThursday() }

// EveryFriday returns a global weekly job for Friday.
// EveryFriday 回傳使用全域 defaultCron 的每週五 Job。
func EveryFriday() *Job { return defaultCron.EveryFriday() }

// EverySaturday returns a global weekly job for Saturday.
// EverySaturday 回傳使用全域 defaultCron 的每週六 Job。
func EverySaturday() *Job { return defaultCron.EverySaturday() }

// --- UntilJob ---
// --- UntilJob（指定時間的一次性工作）---

// UntilJob is a Worker implementation that schedules one execution at a specific time.
// UntilJob 是 Worker 的實作之一，用來在指定時間執行一次任務。
type UntilJob struct {
	untilTime time.Time
	fiber     Fiber
}

// Do binds the callback to a transient Job configured for one-time execution.
// Do 會建立一次性 Job 並綁定回呼，於 untilTime 觸發執行。
func (u *UntilJob) Do(fun any, params ...any) Disposable {
	j := &Job{fiber: u.fiber, jobModel: jobUntil, maximumTimes: 1, atHour: -1, atMinute: -1, atSecond: -1}
	j.nextTime = u.untilTime
	return j.Do(fun, params...)
}

// --- Job ---
// --- Job（排程狀態與執行策略）---

// Job stores schedule state and execution options for cron-style tasks.
// Job 保存 cron 任務的排程狀態與執行選項。
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

// Dispose stops the job and releases its scheduled timer handle.
// Dispose 會停止此 Job，並釋放其排程計時資源。
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

// week configures this job as a weekly schedule for a specific weekday.
// week 會把 Job 設為每週指定星期執行。
func (j *Job) week(dayOfWeek time.Weekday) *Job {
	j.intervalUnit = week
	j.weekday = dayOfWeek
	j.interval = 1
	return j
}

// Days sets interval unit to days.
// Days 設定間隔單位為天。
func (j *Job) Days() *Job {
	j.intervalUnit = day
	return j
}

// Hours sets interval unit to hours.
// Hours 設定間隔單位為小時。
func (j *Job) Hours() *Job {
	j.intervalUnit = hour
	return j
}

// Minutes sets interval unit to minutes.
// Minutes 設定間隔單位為分鐘。
func (j *Job) Minutes() *Job {
	j.intervalUnit = minute
	return j
}

// Seconds sets interval unit to seconds.
// Seconds 設定間隔單位為秒。
func (j *Job) Seconds() *Job {
	j.intervalUnit = second
	return j
}

// Milliseconds sets interval unit to milliseconds.
// Milliseconds 設定間隔單位為毫秒。
func (j *Job) Milliseconds() *Job {
	j.intervalUnit = millisecond
	return j
}

// At sets the specific time of day (hour, minute, second) for job execution.
// Values are clamped: hh mod 24, mm mod 60, ss mod 60. Negative values are made absolute.
// At 可設定每天（或週期對齊）的時分秒；輸入值會被正規化至合法範圍。
func (j *Job) At(hh int, mm int, ss int) *Job {
	j.atHour = Abs(hh) % 24
	j.atMinute = Abs(mm) % 60
	j.atSecond = Abs(ss) % 60
	return j
}

// AfterExecuteTask configures the job to calculate its next execution time after the task finishes.
// This ensures no overlap between consecutive runs. Only applies to delay, second, and millisecond intervals.
// AfterExecuteTask 會在任務完成後才計算下一次執行時間，
// 可避免連續執行重疊（僅適用 delay/second/millisecond）。
func (j *Job) AfterExecuteTask() *Job {
	if j.jobModel == jobDelay || j.intervalUnit == second || j.intervalUnit == millisecond {
		j.afterCalculate = true
	}
	return j
}

// BeforeExecuteTask configures the job to calculate its next execution time immediately,
// without waiting for the current task to finish. This is the default behavior.
// BeforeExecuteTask 會在任務觸發時立即計算下一次時間，不等待本次任務完成（預設模式）。
func (j *Job) BeforeExecuteTask() *Job {
	j.afterCalculate = false
	return j
}

// Times sets the maximum number of times the job will execute.
// Use -1 (default) for unlimited executions.
// Times 設定最多執行次數；-1（預設）代表無限次。
func (j *Job) Times(times int64) *Job {
	j.maximumTimes = times
	return j
}

// Between restricts job execution to a daily time window from f to t (only the time-of-day portion is used).
// Has no effect on delay-mode jobs or if f/t is zero or t <= f.
// Between 可限制每天允許執行的時間視窗（僅使用時分秒，不使用日期）。
// 對 delay 模式、零值時間、或 t<=f 的情況不生效。
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
// Do 會把 builder 參數封裝成實際排程並啟動執行流程，回傳可取消控制代碼。
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
//   - Between 時間窗會在初次排程時統一套用下界對齊（nextTime >= fromTime），
//     並在超過 toTime 時順延窗口到隔天。
func (j *Job) Do(fun any, params ...any) Disposable {
	j.task = newTask(fun, params...)
	if j.jobModel == jobEvery {
		j.interval = normalizePositiveInterval(j.interval)
	}
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
		}
	}

	// 若算出的 nextTime 已經不在未來（今天的 At 時間已過、當前分鐘的秒數已過等），
	// 順延一個 duration，確保首次 timer 一定是正值，不會立即觸發歷史時間點。
	if !j.nextTime.After(now) {
		j.nextTime = j.nextTime.Add(j.duration)
	}

	j.alignNextTimeAtOrAfter(j.fromTime)
	if !j.toTime.IsZero() && j.nextTime.After(j.toTime) {
		j.fromTime = j.fromTime.Add(24 * time.Hour)
		j.toTime = j.toTime.Add(24 * time.Hour)
		j.alignNextTimeAtOrAfter(j.fromTime)
	}

	j.schedule()
	return j
}

// run is the scheduled callback invoked by fiber's timer.
// run 是由 fiber timer 觸發的核心排程回呼。
//
// 設計要點：
//   - run 本身不是週期性 ticker，而是「一次性 timer → 執行 → 再排下一次 timer」的遞迴排程。
//     這樣做的好處是每次都能根據實際執行耗時重新校準下次時間，避免 ticker 的漂移累積。
//   - remainTime < 0 表示「已到期或過期」，才進入執行邏輯；>= 0 表示 timer 提前喚醒
//     （OS 排程抖動），此時直接重新排程補償差值，不執行任務。
//   - afterCalculate 模式下必須先釋放鎖再同步執行 task，因為 task 可能耗時長，
//     持鎖會阻塞 Dispose 等外部操作。執行完後重新上鎖並再檢查 disposed，
//     防止在無鎖期間被外部 Dispose 的情況下繼續排程。
//   - Between 時間窗的判斷分成兩層：canRun 決定「現在能不能跑（需同時滿足 from/to）」，
//     窗口過期檢查（nextTime > toTime）決定「下一次該順延窗口」。
//     兩者獨立運作，確保窗口外不執行但仍持續排程。
func (j *Job) run() {
	j.mu.Lock()
	if j.disposed {
		j.mu.Unlock()
		return
	}

	adjustTime := j.remainTime()
	if adjustTime < 0 {
		now := time.Now()
		canRun := true
		if !j.fromTime.IsZero() && now.Before(j.fromTime) {
			canRun = false
		}
		if !j.toTime.IsZero() && now.After(j.toTime) {
			canRun = false
		}

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

		if !j.fromTime.IsZero() && now.Before(j.fromTime) {
			j.alignNextTimeAtOrAfter(j.fromTime)
		} else {
			// 推進 nextTime 到下一個週期。
			j.nextTime = j.nextTime.Add(j.duration)
		}

		// 若 nextTime 超出 Between 時間窗，將整個窗口順延 24 小時。
		if !j.toTime.IsZero() && j.nextTime.After(j.toTime) {
			j.fromTime = j.fromTime.Add(24 * time.Hour)
			j.toTime = j.toTime.Add(24 * time.Hour)
			j.alignNextTimeAtOrAfter(j.fromTime)
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
// remainTime 回傳距離 nextTime 的毫秒差；負值表示已到期（或已過期）。
func (j *Job) remainTime() int64 {
	return j.nextTime.Sub(time.Now()).Milliseconds()
}

// alignNextTimeAtOrAfter advances nextTime so it is not earlier than the given bound.
// This keeps schedule alignment consistent when a fromTime window is configured.
// alignNextTimeAtOrAfter 會推進 nextTime，確保不早於指定下界時間。
// 此方法可在設定 fromTime 視窗時維持對齊一致性。
func (j *Job) alignNextTimeAtOrAfter(bound time.Time) {
	if bound.IsZero() || !j.nextTime.Before(bound) {
		return
	}
	if j.duration <= 0 {
		j.nextTime = bound
		return
	}

	delta := bound.Sub(j.nextTime)
	steps := delta / j.duration
	if delta%j.duration != 0 {
		steps++
	}
	j.nextTime = j.nextTime.Add(time.Duration(steps) * j.duration)
}

// schedule arms the initial timer for this job.
// schedule 會為此 Job 啟動初次 timer 排程。
func (j *Job) schedule() {
	j.mu.Lock()
	diff := j.remainTime()
	j.taskDisposer = j.fiber.Schedule(diff, j.run)
	j.mu.Unlock()
}
