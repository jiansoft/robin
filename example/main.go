package main

import (
    "fmt"
    "log"
    "time"

    "github.com/jiansoft/robin"
)

func main() {
    now := time.Now()
    log.Printf("Start at %v", now.Format(time.RFC3339))

    robin.Every(450).Milliseconds().Times(2).Do(runCron, "every 1 Milliseconds Times 2")
    robin.Every(1).Seconds().Times(3).Do(runCron, "every 1 Seconds Times 3")
    robin.Every(10).Minutes().Do(runCron, "every 10 Minutes")
    robin.Every(10).Minutes().AfterExecuteTask().Do(runCronAndSleep, "every 10 Minutes and sleep 4 Minutes", 4*60*1000)
    robin.Every(60).Seconds().AfterExecuteTask().Do(runCronAndSleep, "every 60 Seconds and sleep 4 Minutes", 4*60*1000)
    robin.Delay(4000).Times(4).Do(runCron, "Delay 4000 ms Times 4")
    robin.Delay(21).Seconds().Do(func() {
        p1 := player{NickName: "Player 1"}
        p2 := player{NickName: "Player 2"}
        p3 := player{NickName: "Player 3"}
        p4 := player{NickName: "Player 4"}
        channel := robin.NewChannel()

        channel.Subscribe(p1.eventFinalBossResurge)
        channel.Subscribe(p2.eventFinalBossResurge)
        channel.Subscribe(p3.eventFinalBossResurge)
        p4Unsubscribe := channel.Subscribe(p4.eventFinalBossResurge)

        log.Printf("the channel have %v subscribers.", channel.Count())
        channel.Publish("The boss resurge first.")
        <- time.After(time.Duration(10* time.Millisecond))

        p4Unsubscribe.Dispose()
        log.Printf("the channel have %v subscribers.", channel.Count())
        channel.Publish("The boss resurge second.")
        <- time.After(time.Duration(10* time.Millisecond))

        channel.Clear()
        log.Printf("the channel have %v subscribers.", channel.Count())
        channel.Publish("The boss resurge third.")
        <- time.After(time.Duration(10* time.Millisecond))

        channel.Subscribe(p1.eventFinalBossResurge)
        log.Printf("the channel have %v subscribers.", channel.Count())
        channel.Publish("The boss resurge fourth.")
        <- time.After(time.Duration(10* time.Millisecond))
    })

    now = now.Add(time.Duration(17*time.Second + 100))
    robin.EveryMonday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Monday")
    robin.EveryTuesday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Tuesday")
    robin.EveryWednesday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Wednesday")
    robin.EveryThursday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Thursday")
    robin.EveryFriday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Friday")
    robin.EverySaturday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Saturday")
    robin.EverySunday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Sunday")

    now = now.Add(time.Duration(1 * time.Second))
    robin.Every(1).Hours().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "every 1 Hours")

    now = now.Add(time.Duration(1 * time.Second))
    robin.Every(1).Days().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "every 1 Days")

    now = now.Add(time.Duration(1 * time.Second))
    robin.Everyday().At(now.Hour(), now.Minute(), now.Second()).Do(runCron, "Everyday")

    _, _ = fmt.Scanln()

    var runCronFiber = robin.NewGoroutineSingle()
    var runCronFiber2 = robin.NewGoroutineSingle()
    runCronFiber.Start()
    runCronFiber2.Start()

    runCronFiber.Schedule(0, func() {
        log.Printf("just test\n")
        //	a.Dispose()
    })

    robin.Delay(2000).Do(runCron, "a Delay 2000 ms")
    //robin.every(1000).Milliseconds().Do(runCron, "every(1000).Milliseconds()")
    //robin.every(2).Seconds().Do(runCron, "every(2).Seconds()")
    //every N seconds do once.
    //robin.every(10).Seconds().Do(runCron, "every 10 Seconds")
    runCronFiber.ScheduleOnInterval(10000, 10000, func() {
        log.Printf("runCronFiber 1\n")
    })
    runCronFiber2.ScheduleOnInterval(1000, 1000, func() {
        //log.Printf("runCronFiber 2\n")
    })
    runCronFiber2.ScheduleOnInterval(1000, 1000, func() {
        //log.Printf("runCronFiber 3\n")
    })

    //Make a cron that will be executed everyday at 15:30:04(HH:mm:ss)
    b := robin.Every(1).Days().At(15, 30, 4).Do(runCron, "Will be cancel")

    //After one second, a & b crons well be canceled by fiber
    runCronFiber.Schedule(1000, func() {
        log.Printf("a & b are dispose\n")
        b.Dispose()
        //	a.Dispose()
    })
    minute := 4
    second := 10
    robin.Every(int64(second)).Seconds().AfterExecuteTask().Do(runSleepCron, "AfterExecuteTask", second, minute)
    robin.Every(int64(second)).Seconds().Do(runSleepCron, "BeforeExecuteTask", second, 0)

    //every friday do once at 11:50:00(HH:mm:ss).
    robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

    //every N day do once at 11:50:00(HH:mm:ss)
    robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

    //every N hours do once at N:50:00(HH:mm:ss).
    robin.Every(1).Hours().At(0, minute, second).Do(runCron, "every 1 Hours")

    //every N minutes do once.
    robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "every 1 Minutes")

    //every N seconds do once.
    robin.Every(100).Seconds().Do(runCron, "every 10 Seconds")

    robin.Delay(2000).Times(3).AfterExecuteTask().Do(CronTestAndSleepASecond, " Delay 2000 ms Times 2 AfterExecuteTask")
    _, _ = fmt.Scanln()
}

func runSleepCron(s string, second int, sleep int) {
    log.Printf("%s every %d second and sleep %d  now %v\n", s, second, sleep, time.Now())
    time.Sleep(time.Duration(sleep) * time.Second)
}

func runCron(s string) {
    log.Printf("I am %s CronTest %v\n", s, time.Now())
}

func runCronAndSleep(s string, sleepInMs int) {
    log.Printf("I am %s CronTest %v\n", s, time.Now())
    timeout := time.NewTimer(time.Duration(sleepInMs) * time.Millisecond)
    select {
    case <-timeout.C:
    }
}

func CronTestAndSleepASecond(s string) {
    log.Printf("I am %s CronTest and sleep a second %v\n", s, time.Now())
    timeout := time.NewTimer(time.Duration(1000) * time.Millisecond)
    select {
    case <-timeout.C:
    }
}

type player struct {
    NickName string
}

func (p player) eventFinalBossResurge(someBossInfo string) {
    log.Printf("%s receive a message : %s", p.NickName, someBossInfo)
}
