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
			Every(1).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")
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
		})
	}
}

func TestDelaySeries(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_DelaySeries"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Delay(50).MilliSeconds().Do(func(s string) { t.Logf("s:%v", s) }, "MilliSeconds")
			Delay(50).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds")
			Delay(50).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes")
			Delay(50).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours")
			Delay(50).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")

			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}

		})
	}
}
