package robin

import (
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

            Every(2).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")
            Every(2).Days().At(0, 1, 2).Do(func(s string) { t.Logf("s:%v", s) }, "Days")

            hh := time.Now().Hour()
            mm := time.Now().Minute()
            ss := time.Now().Second()
            t.Logf("now:%v:%v:%v", hh, mm, ss)
            Every(200).Milliseconds().Times(3).Do(func(s string) { t.Logf("s:%v", s) }, "every 200 Milliseconds times 3")

            h := Every(1200).Milliseconds().Do(func(s string) { t.Logf("s:%v", s) }, "every 3 Seconds")
            everyOneSecondDispose := Every(1).Seconds().Do(func(s string) {
                t.Logf("s:%v", s)
                h.Dispose()
            }, "h.Dispose")

            everyOneSecond := Every(1).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Seconds")
            everySecondSecond := Every(2).Seconds().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Seconds")

            Every(1).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Minutes")
            Every(1).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Minutes")
            Every(2).Minutes().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Minutes")

            Every(1).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Hours")
            Every(1).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Hours")
            Every(2).Hours().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Hours")

            Everyday().Do(func(s string) { t.Logf("s:%v", s) }, "Everyday")
            Everyday().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "Everyday")
            Every(1).Days().Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Days")
            Every(1).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 1 Days")
            Every(2).Days().At(hh, mm, ss+1).Do(func(s string) { t.Logf("s:%v", s) }, "every 2 Days")

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

            timeout.Reset(time.Duration(2500) * time.Millisecond)
            select {
            case <-timeout.C:
            }

            everyOneSecondDispose.Dispose()
            everyOneSecond.Dispose()
            everySecondSecond.Dispose()
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
            RightNow().Times(3).Do(func(s string) { t.Logf("s:%v", s) }, "RightNow time 3")
            RightNow().Do(func(s string) { t.Logf("s:%v", s) }, "RightNow")
            Delay(5).Milliseconds().Do(func(s string) { t.Logf("s:%v", s) }, "Milliseconds")
            Delay(50).Seconds().Do(func(s string) { t.Logf("s:%v", s) }, "Seconds")
            Delay(50).Minutes().Do(func(s string) { t.Logf("s:%v", s) }, "Minutes")
            Delay(50).Hours().Do(func(s string) { t.Logf("s:%v", s) }, "Hours")
            Delay(50).Days().Do(func(s string) { t.Logf("s:%v", s) }, "Days")

            unixMillisecond := time.Now().UnixNano() / int64(time.Millisecond)
            Delay(10).Milliseconds().Do(func(s string, pervUnixMillisecond int64) {
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
            }, "Milliseconds 10", unixMillisecond)

            timeout := time.NewTimer(time.Duration(2500) * time.Millisecond)
            select {
            case <-timeout.C:
            }
        })
    }
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
