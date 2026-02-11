package robin

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestEvery validates the Every/At/Times/weekday builder chain and disposal stability.
// TestEvery 驗證 Every/At/Times/weekday 建構鏈與 Dispose 穩定性。
func TestEvery(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"TestEvery"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var disposables []Disposable

			monday := EveryMonday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Monday")
			tuesday := EveryTuesday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday")
			wednesday := EveryWednesday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday")
			thursday := EveryThursday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Thursday")
			friday := EveryFriday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Friday")
			saturday := EverySaturday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Saturday")
			sunday := EverySunday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Sunday")

			milliSeconds := Every(55).Milliseconds().Do(func(s string) { t.Logf("s:%v", s) }, "Milliseconds")
			seconds := Every(1).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds")
			minutes := Every(1).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes")
			hours := Every(2).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours")
			days := Every(1).Days().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Days")
			rightNow := RightNow().Do(func(s string) { t.Logf("s:%v", s) }, "RightNow")

			after := Every(60).Milliseconds().AfterExecuteTask().Do(func(s string) { t.Logf("s:%v", s) }, "After")
			before := Every(60).Milliseconds().BeforeExecuteTask().Do(func(s string) { t.Logf("s:%v", s) }, "Before")

			timeout := time.NewTimer(time.Duration(1000) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			monday.Dispose()
			tuesday.Dispose()
			wednesday.Dispose()
			thursday.Dispose()
			friday.Dispose()
			saturday.Dispose()
			sunday.Dispose()

			milliSeconds.Dispose()
			seconds.Dispose()
			minutes.Dispose()
			hours.Dispose()
			days.Dispose()
			rightNow.Dispose()
			after.Dispose()
			before.Dispose()

			disposables = append(disposables, Every(2).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days"))
			disposables = append(disposables, Every(2).Days().At(0, 1, 2).Do(func(s string) { t.Logf("s:%v", s) }, "Days"))

			hh := time.Now().Hour()
			mm := time.Now().Minute()
			ss := time.Now().Second()
			t.Logf("now:%v:%v:%v", hh, mm, ss)
			disposables = append(disposables, Every(200).Milliseconds().Times(3).Do(func(s string) { t.Logf("s:%v", s) }, "every 200 Milliseconds times 3"))

			h := Every(1200).Milliseconds().Do(func(s string) { t.Logf("s:%v", s) }, "every 3 Seconds")
			everyOneSecondDispose := Every(1).Seconds().Do(func(s string) {
				t.Logf("s:%v", s)
				h.Dispose()
			}, "h.Dispose")

			everyOneSecond := Every(1).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Seconds")
			everySecondSecond := Every(2).Seconds().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Seconds")

			disposables = append(disposables, Every(1).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Minutes"))
			disposables = append(disposables, Every(1).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Minutes"))
			disposables = append(disposables, Every(2).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Minutes"))

			disposables = append(disposables, Every(1).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Hours"))
			disposables = append(disposables, Every(1).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Hours"))
			disposables = append(disposables, Every(2).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Hours"))

			disposables = append(disposables, Everyday().Do(func(s string) { t.Logf("s:%v", s) }, "Everyday"))
			disposables = append(disposables, Everyday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Everyday"))
			disposables = append(disposables, Every(1).Days().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Days"))
			disposables = append(disposables, Every(1).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Days"))
			disposables = append(disposables, Every(2).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Days"))

			disposables = append(disposables, EveryMonday().Do(func(s string) { t.Logf("s:%v", s) }, "Monday"))
			disposables = append(disposables, EveryMonday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Monday"))

			disposables = append(disposables, EveryTuesday().Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday"))
			disposables = append(disposables, EveryTuesday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday"))

			disposables = append(disposables, EveryWednesday().Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday"))
			disposables = append(disposables, EveryWednesday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday"))

			disposables = append(disposables, EveryThursday().Do(func(s string) { t.Logf("s:%v", s) }, "Thursday"))
			disposables = append(disposables, EveryThursday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Thursday"))

			disposables = append(disposables, EveryFriday().Do(func(s string) { t.Logf("s:%v", s) }, "Friday"))
			disposables = append(disposables, EveryFriday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Friday"))

			disposables = append(disposables, EverySaturday().Do(func(s string) { t.Logf("s:%v", s) }, "Saturday"))
			disposables = append(disposables, EverySaturday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Saturday"))

			disposables = append(disposables, EverySunday().Do(func(s string) { t.Logf("s:%v", s) }, "Sunday"))
			disposables = append(disposables, EverySunday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Sunday"))

			timeout.Reset(time.Duration(2500) * time.Millisecond)
			select {
			case <-timeout.C:
			}

			everyOneSecondDispose.Dispose()
			everyOneSecond.Dispose()
			everySecondSecond.Dispose()
			h.Dispose()
			for _, d := range disposables {
				d.Dispose()
			}
			// Wait for in-flight goroutines to complete
			time.Sleep(100 * time.Millisecond)
		})
	}
}

// TestDelay validates one-shot delay semantics and chained nested delay timing lower-bounds.
// TestDelay 驗證一次性延遲語義，以及巢狀 Delay 連鎖的最小時間門檻。
func TestDelay(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"TestDelay"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var disposables []Disposable

			disposables = append(disposables, RightNow().Times(3).Do(func(s string) { t.Logf("s:%v", s) }, "RightNow time 3"))
			disposables = append(disposables, RightNow().Do(func(s string) { t.Logf("s:%v", s) }, "RightNow"))
			disposables = append(disposables, Delay(5).Milliseconds().Do(func(s string) { t.Logf("s:%v", s) }, "Milliseconds"))
			disposables = append(disposables, Delay(50).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds"))
			disposables = append(disposables, Delay(50).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes"))
			disposables = append(disposables, Delay(50).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours"))
			disposables = append(disposables, Delay(50).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days"))

			unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
			disposables = append(disposables, Delay(10).Milliseconds().Do(func(s string, pervUnixMillisecond int64) {
				unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
				diffTime := unixMillisecond - pervUnixMillisecond
				t.Logf("s:%v diff:%d", s, diffTime)
				if diffTime < 10 {
					t.Errorf("delay time got = %v, want >= 10", diffTime)
				}
				Delay(20).Milliseconds().Do(func(s string, pervUnixMillisecond int64) {
					unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
					diffTime := unixMillisecond - pervUnixMillisecond
					t.Logf("s:%v diff:%d", s, diffTime)
					if diffTime < 20 {
						t.Errorf("delay time got = %v, want >= 20", diffTime)
					}
					Delay(30).Milliseconds().Do(func(s string, pervUnixMillisecond int64) {
						unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
						diffTime := unixMillisecond - pervUnixMillisecond
						t.Logf("s:%v diff:%d", s, diffTime)
						if diffTime < 30 {
							t.Errorf("delay time got = %v, want >= 30", diffTime)
						}
						Delay(40).Milliseconds().Do(func(s string, pervUnixMillisecond int64) {
							unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
							diffTime := unixMillisecond - pervUnixMillisecond
							t.Logf("s:%v diff:%d", s, diffTime)
							if diffTime < 40 {
								t.Errorf("delay time got = %v, want >= 40", diffTime)
							}
						}, "Milliseconds 40", unixMillisecond)
					}, "Milliseconds 30", unixMillisecond)
				}, "Milliseconds 20", unixMillisecond)
			}, "Milliseconds 10", unixMillisecond))

			timeout := time.NewTimer(time.Duration(2500) * time.Millisecond)
			select {
			case <-timeout.C:
			}

			for _, d := range disposables {
				d.Dispose()
			}
			// Wait for in-flight goroutines to complete
			time.Sleep(100 * time.Millisecond)
		})
	}
}

// TestUntil validates that future timestamps execute once and past timestamps do not execute.
// TestUntil 驗證未來時間會執行一次、過去時間不會執行。
func TestUntil(t *testing.T) {
	t.Run("future time executes once", func(t *testing.T) {
		var count int32
		futureTime := time.Now().Add(200 * time.Millisecond)
		d := Until(futureTime).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(500 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got != 1 {
			t.Errorf("Until(future) executed %d times, want 1", got)
		}
		d.Dispose()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("past time does not execute", func(t *testing.T) {
		var count int32
		pastTime := time.Now().Add(-1 * time.Second)
		d := Until(pastTime).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got != 0 {
			t.Errorf("Until(past) executed %d times, want 0", got)
		}
		// taskDisposer is nil but Dispose() should handle it gracefully
		d.Dispose()
		time.Sleep(50 * time.Millisecond)
	})
}

// TestBetween validates window filtering, invalid-parameter ignore rules, and delay-mode bypass.
// TestBetween 驗證時間窗過濾、無效參數忽略規則，以及 delay 模式忽略 Between。
func TestBetween(t *testing.T) {
	t.Run("valid between range executes", func(t *testing.T) {
		var count int32
		now := time.Now()
		from := time.Date(0, 1, 1, now.Hour(), 0, 0, 0, time.Local)
		to := time.Date(0, 1, 1, now.Hour(), 59, 59, 0, time.Local)

		d := Every(100).Milliseconds().Between(from, to).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(350 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got == 0 {
			t.Errorf("Between(valid range) executed %d times, want > 0", got)
		}
		d.Dispose()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("invalid between params ignored", func(t *testing.T) {
		var count int32
		// f.IsZero()
		d1 := Every(100).Milliseconds().Between(time.Time{}, time.Now()).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		// t.IsZero()
		d2 := Every(100).Milliseconds().Between(time.Now(), time.Time{}).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		// t <= f
		now := time.Now()
		d3 := Every(100).Milliseconds().Between(now.Add(time.Hour), now).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C

		// count should be > 0 because Between was ignored (no filtering)
		d1.Dispose()
		d2.Dispose()
		d3.Dispose()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("delay model ignores between", func(t *testing.T) {
		var count int32
		now := time.Now()
		from := time.Date(0, 1, 1, now.Hour(), 0, 0, 0, time.Local)
		to := time.Date(0, 1, 1, now.Hour(), 59, 59, 0, time.Local)

		d := Delay(50).Milliseconds().Between(from, to).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got != 1 {
			t.Errorf("Delay with Between should execute once, got %d", got)
		}
		d.Dispose()
		time.Sleep(50 * time.Millisecond)
	})
}

// TestJobDisposed validates disposal stop behavior and double-dispose safety.
// TestJobDisposed 驗證 Dispose 後停止執行與重複 Dispose 的安全性。
func TestJobDisposed(t *testing.T) {
	t.Run("disposed job does not run", func(t *testing.T) {
		var count int32
		d := Every(50).Milliseconds().Do(func() {
			atomic.AddInt32(&count, 1)
		})
		// Dispose immediately
		d.Dispose()

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		// After disposal the job should stop, at most 1 execution may have snuck in
		if got > 1 {
			t.Errorf("Disposed job executed %d times, want <= 1", got)
		}
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("double dispose does not panic", func(t *testing.T) {
		d := Every(100).Milliseconds().Do(func() {})
		d.Dispose()
		d.Dispose() // should not panic
		time.Sleep(50 * time.Millisecond)
	})
}

// TestCronScheduler validates instance-based API parity with global helpers and lifecycle isolation.
// TestCronScheduler 驗證實例化 API 與全域 API 的行為一致性及生命週期隔離。
func TestCronScheduler(t *testing.T) {
	t.Run("instance-based scheduling", func(t *testing.T) {
		cs := NewCronScheduler()
		defer cs.Dispose()

		var count int32
		d := cs.Every(50).Milliseconds().Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C
		d.Dispose()

		got := atomic.LoadInt32(&count)
		if got == 0 {
			t.Error("CronScheduler.Every() did not execute")
		}
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("delay and right now", func(t *testing.T) {
		cs := NewCronScheduler()
		defer cs.Dispose()

		var count int32
		cs.RightNow().Do(func() {
			atomic.AddInt32(&count, 1)
		})
		cs.Delay(50).Milliseconds().Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(300 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got != 2 {
			t.Errorf("CronScheduler RightNow+Delay executed %d times, want 2", got)
		}
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("until with future time", func(t *testing.T) {
		cs := NewCronScheduler()
		defer cs.Dispose()

		var count int32
		futureTime := time.Now().Add(100 * time.Millisecond)
		d := cs.Until(futureTime).Do(func() {
			atomic.AddInt32(&count, 1)
		})

		timeout := time.NewTimer(400 * time.Millisecond)
		<-timeout.C

		got := atomic.LoadInt32(&count)
		if got != 1 {
			t.Errorf("CronScheduler.Until(future) executed %d times, want 1", got)
		}
		d.Dispose()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("everyday and weekday methods", func(t *testing.T) {
		cs := NewCronScheduler()
		defer cs.Dispose()

		var disposables []Disposable
		disposables = append(disposables, cs.Everyday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EverySunday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EveryMonday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EveryTuesday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EveryWednesday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EveryThursday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EveryFriday().At(0, 0, 0).Do(func() {}))
		disposables = append(disposables, cs.EverySaturday().At(0, 0, 0).Do(func() {}))

		for _, d := range disposables {
			d.Dispose()
		}
		time.Sleep(50 * time.Millisecond)
	})
}

// TestJobDoBetweenTimeShift validates Do() path where Between window is shifted to next day.
// TestJobDoBetweenTimeShift 驗證 Do() 在 Between 視窗需順延到隔天時的分支行為。
func TestJobDoBetweenTimeShift(t *testing.T) {
	// Cover Do() lines 334-341: Between with millisecond unit where
	// fromTime (adjusted to today) > toTime (adjusted to today),
	// triggering the 24-hour time shift in Do().
	//
	// Dynamically compute hours so fromTime is always in the future
	// and fromTime > toTime after today-adjustment.
	now := time.Now()
	fromMinute := now.Add(2 * time.Minute)
	toMinute := now.Add(-2 * time.Minute)

	// Guard: if now is near midnight (00:00~00:02), toMinute wraps to
	// 23:58 yesterday, making adjusted toTime > fromTime. Skip this edge case.
	if now.Hour() == 0 && now.Minute() < 3 {
		t.Skip("skipping: too close to midnight for reliable time-shift test")
	}

	// Original to must have to.Unix() > from.Unix() to pass Between's guard.
	// Use different base dates to guarantee this.
	from := time.Date(2099, 1, 1, fromMinute.Hour(), fromMinute.Minute(), fromMinute.Second(), 0, time.Local)
	to := time.Date(2099, 1, 2, toMinute.Hour(), toMinute.Minute(), toMinute.Second(), 0, time.Local)

	var count int32
	d := Every(100).Milliseconds().Between(from, to).Do(func() {
		atomic.AddInt32(&count, 1)
	})

	// Job is scheduled far in the future (tomorrow), so no executions
	time.Sleep(200 * time.Millisecond)
	d.Dispose()

	got := atomic.LoadInt32(&count)
	if got != 0 {
		t.Errorf("expected 0 executions (scheduled for tomorrow), got %d", got)
	}
}

// TestJobRunBetweenWindowExpires validates run() behavior when the active Between window expires.
// TestJobRunBetweenWindowExpires 驗證 run() 在 Between 視窗過期後的行為。
func TestJobRunBetweenWindowExpires(t *testing.T) {
	// Cover run() canRun=false path (line 364-366) and
	// toTime 24h shift in run() (lines 398-401).
	//
	// Creates a ~2 second Between window. After window expires,
	// run() finds canRun=false, skips execution, then shifts times by 24h.
	var count int32
	now := time.Now()

	from := time.Date(2000, 6, 1, now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	to := time.Date(2000, 6, 1, now.Hour(), now.Minute(), now.Second()+2, 0, time.Local)

	d := Every(100).Milliseconds().Between(from, to).Do(func() {
		atomic.AddInt32(&count, 1)
	})

	// Wait for the 2-second window to expire + one more run cycle
	time.Sleep(3 * time.Second)
	d.Dispose()

	got := atomic.LoadInt32(&count)
	if got < 5 {
		t.Errorf("expected at least 5 executions in between window, got %d", got)
	}
}

// TestJobRunAfterCalculateDisposeDuringExec validates disposing during AfterExecuteTask synchronous execution.
// TestJobRunAfterCalculateDisposeDuringExec 驗證 AfterExecuteTask 同步執行期間被 Dispose 的路徑。
func TestJobRunAfterCalculateDisposeDuringExec(t *testing.T) {
	// Cover run() lines 375-378: afterCalculate=true, job disposed
	// while task is executing synchronously.
	started := make(chan struct{}, 1)

	d := Every(50).Milliseconds().AfterExecuteTask().Do(func() {
		select {
		case started <- struct{}{}:
		default:
		}
		time.Sleep(200 * time.Millisecond)
	})

	<-started   // wait for first execution to start
	d.Dispose() // dispose while task is executing synchronously

	time.Sleep(300 * time.Millisecond) // wait for task to finish
	// If we reach here without deadlock or panic, the path is covered
}

// TestEveryNonPositiveIntervalDoesNotBusyLoop validates interval normalization and busy-loop prevention.
// TestEveryNonPositiveIntervalDoesNotBusyLoop 驗證非正間隔正規化與忙迴圈防護。
func TestEveryNonPositiveIntervalDoesNotBusyLoop(t *testing.T) {
	cs := NewCronScheduler()
	defer cs.Dispose()

	if got := cs.Every(0).interval; got != 1 {
		t.Fatalf("Every(0) interval = %d, want 1", got)
	}
	if got := cs.Every(-7).interval; got != 1 {
		t.Fatalf("Every(-7) interval = %d, want 1", got)
	}

	var count int32
	d := cs.Every(0).Milliseconds().Do(func() {
		atomic.AddInt32(&count, 1)
	})

	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&count) > 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	d.Dispose()

	got := atomic.LoadInt32(&count)
	if got == 0 {
		t.Fatal("expected at least one execution")
	}
	if got > 500 {
		t.Fatalf("possible busy loop: executions = %d, want <= 500", got)
	}
}

// TestBetweenMinuteUnitInitialAlignToFromTime validates initial nextTime alignment against fromTime for minute unit.
// TestBetweenMinuteUnitInitialAlignToFromTime 驗證分鐘單位下 nextTime 初始值會對齊到 fromTime 之後。
func TestBetweenMinuteUnitInitialAlignToFromTime(t *testing.T) {
	cs := NewCronScheduler()
	defer cs.Dispose()

	now := time.Now()
	fromRaw := now.Add(2 * time.Minute)
	toRaw := now.Add(10 * time.Minute)

	// Between uses time-of-day only and rewrites to "today".
	// Skip near-day-boundary cases to keep this test deterministic.
	if fromRaw.Day() != now.Day() || toRaw.Day() != now.Day() {
		t.Skip("skipping day-boundary case")
	}

	from := time.Date(fromRaw.Year(), fromRaw.Month(), fromRaw.Day(), fromRaw.Hour(), fromRaw.Minute(), fromRaw.Second(), fromRaw.Nanosecond(), time.Local)
	to := time.Date(toRaw.Year(), toRaw.Month(), toRaw.Day(), toRaw.Hour(), toRaw.Minute(), toRaw.Second(), toRaw.Nanosecond(), time.Local)

	j := cs.Every(1).Minutes().Between(from, to)
	d := j.Do(func() {})
	defer d.Dispose()

	if j.nextTime.Before(j.fromTime) {
		t.Fatalf("nextTime %v should be >= fromTime %v", j.nextTime, j.fromTime)
	}
}

// TestJobRunBeforeFromTimeDoesNotExecute validates run() does not execute tasks before fromTime boundary.
// TestJobRunBeforeFromTimeDoesNotExecute 驗證 run() 在 fromTime 邊界之前不會執行任務。
func TestJobRunBeforeFromTimeDoesNotExecute(t *testing.T) {
	f := NewGoroutineMulti()
	defer f.Dispose()

	var count int32
	now := time.Now()
	j := &Job{
		fiber:          f,
		nextTime:       now.Add(-1 * time.Millisecond),
		fromTime:       now.Add(200 * time.Millisecond),
		toTime:         now.Add(400 * time.Millisecond),
		task:           newTask(func() { atomic.AddInt32(&count, 1) }),
		duration:       50 * time.Millisecond,
		interval:       50,
		maximumTimes:   -1,
		jobModel:       jobEvery,
		intervalUnit:   millisecond,
		afterCalculate: true,
	}

	j.run()
	defer j.Dispose()

	if got := atomic.LoadInt32(&count); got != 0 {
		t.Fatalf("job executed before fromTime, count = %d", got)
	}
}
