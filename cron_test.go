package robin

import (
	"testing"
	"time"
)

func TestEverySeries(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_EverySeries"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monday := EveryMonday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Monday")
			tuesday := EveryTuesday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday")
			wednesday := EveryWednesday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday")
			thursday := EveryThursday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Thursday")
			firday := EveryFriday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Friday")
			saturday := EverySaturday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Saturday")
			sunday := EverySunday().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Sunday")

			milliSeconds := Every(55).MilliSeconds().Do(func(s string) { t.Logf("s:%v", s) }, "MilliSeconds")
			seconds := Every(1).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds")
			minutes := Every(1).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes")
			hours := Every(2).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours")
			days := Every(1).Days().At(0, 0, 0).Do(func(s string) { t.Logf("s:%v", s) }, "Days")
			rightNow := RightNow().Do(func(s string) { t.Logf("s:%v", s) }, "RightNow")

			after := Every(60).MilliSeconds().AfterExecuteTask().Do(func(s string) { t.Logf("s:%v", s) }, "After")
			before := Every(60).MilliSeconds().BeforeExecuteTask().Do(func(s string) { t.Logf("s:%v", s) }, "Before")

			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			monday.Dispose()
			tuesday.Dispose()
			wednesday.Dispose()
			thursday.Dispose()
			firday.Dispose()
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

			Every(2).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")
			Every(2).Days().At(0, 1, 2).Do(func(s string) { t.Logf("s:%v", s) }, "Days")

			hh := time.Now().Hour()
			mm := time.Now().Minute()
			ss := time.Now().Second()
			t.Logf("now:%v:%v:%v", hh, mm, ss)
			Every(1).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Seconds")
			Every(2).Seconds().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 2 Seconds")

			Every(1).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Minutes")
			Every(1).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Minutes")
			Every(2).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 2 Minutes")

			Every(1).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Hours")
			Every(1).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Hours")
			Every(2).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 2 Hours")

			Every(1).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Days")
			Every(1).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 1 Days")
			Every(2).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Every 2 Days")

			EveryMonday().Do(func(s string) { t.Logf("s:%v", s) }, "Monday")
			EveryMonday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Monday")

			EveryTuesday().Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday")
			EveryTuesday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Tuesday")

			EveryWednesday().Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday")
			EveryWednesday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Wednesday")

			EveryThursday().Do(func(s string) { t.Logf("s:%v", s) }, "Thursday")
			EveryThursday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Thursday")

			EveryFriday().Do(func(s string) { t.Logf("s:%v", s) }, "Friday")
			EveryFriday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Friday")

			EverySaturday().Do(func(s string) { t.Logf("s:%v", s) }, "Saturday")
			EverySaturday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Saturday")

			EverySunday().Do(func(s string) { t.Logf("s:%v", s) }, "Sunday")
			EverySunday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Sunday")

			timeout.Reset(time.Duration(1500) * time.Millisecond)
			select {
			case <-timeout.C:
			}

		})
	}
}

func TestDelaySeries(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_DelaySeries_1"},
		{"Test_DelaySeries_2"},
		{"Test_DelaySeries_3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RightNow().Do(func(s string) { t.Logf("s:%v", s) }, "RightNow")
			Delay(5).MilliSeconds().Do(func(s string) { t.Logf("s:%v", s) }, "MilliSeconds")
			Delay(50).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds")
			Delay(50).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes")
			Delay(50).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours")
			Delay(50).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")

			unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
			Delay(10).MilliSeconds().Do(func(s string, pervUnixMillisecond int64) {
				unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
				diffTime := unixMillisecond - pervUnixMillisecond
				t.Logf("s:%v diff:%d", s, diffTime)
				if 10 < diffTime {
					t.Errorf("delay time got = %v, want >= 10", unixMillisecond)
				}
				Delay(20).MilliSeconds().Do(func(s string, pervUnixMillisecond int64) {
					unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
					diffTime := unixMillisecond - pervUnixMillisecond
					t.Logf("s:%v diff:%d", s, diffTime)
					if 20 < diffTime {
						t.Errorf("delay time got = %v, want >= 20", unixMillisecond)
					}
					Delay(30).MilliSeconds().Do(func(s string, pervUnixMillisecond int64) {
						unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
						diffTime := unixMillisecond - pervUnixMillisecond
						t.Logf("s:%v diff:%d", s, diffTime)
						if 30 < diffTime {
							t.Errorf("delay time got = %v, want >= 30", unixMillisecond)
						}
						Delay(40).MilliSeconds().Do(func(s string, pervUnixMillisecond int64) {
							unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
							diffTime := unixMillisecond - pervUnixMillisecond
							t.Logf("s:%v diff:%d", s, diffTime)
							if 40 < diffTime {
								t.Errorf("delay time got = %v, want >= 40", unixMillisecond)
							}
						}, "MilliSeconds 40", unixMillisecond)
					}, "MilliSeconds 30", unixMillisecond)
				}, "MilliSeconds 20", unixMillisecond)
			}, "MilliSeconds 10", unixMillisecond)

			timeout := time.NewTimer(time.Duration(120) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}
