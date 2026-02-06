package robin

import (
	"sync/atomic"
	"testing"
	"time"
)

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

func TestBetween(t *testing.T) {
	t.Run("valid between range executes", func(t *testing.T) {
		var count int32
		now := time.Now()
		from := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second()-1, 0, time.Local)
		to := time.Date(0, 1, 1, now.Hour(), now.Minute()+5, 0, 0, time.Local)

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
		from := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second()-1, 0, time.Local)
		to := time.Date(0, 1, 1, now.Hour(), now.Minute()+5, 0, 0, time.Local)

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

func TestAbs(t *testing.T) {
	type args struct {
		a int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test_Abs_1", args{a: 1}, 1},
		{"Test_Abs_2", args{a: -1}, 1},
		{"Test_Abs_3", args{a: -123}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.a); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
