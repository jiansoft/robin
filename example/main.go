package main

import (
	"log"
	"time"

	"github.com/jiansoft/robin"
)

var quitSemaphore chan bool

func main() {
	log.Printf("Start\n")
    RunChannelTest()
	var runCronFiber = robin.NewGoroutineMulti()
	runCronFiber.Start()
	_ = robin.Delay(2000).Do(runCron, "a Delay 2000 ms")

	//Make a cron that will be executed everyday at 15:30:04(HH:mm:ss)
	b := robin.Every(1).Days().At(15, 30, 4).Do(runCron, "Will be cancel")

	//After one second, a & b crons well be canceled by fiber
	runCronFiber.Schedule(1000, func() {
		log.Printf("a & b are dispose\n")
		b.Dispose()
		//	a.Dispose()
	})
	minute := 41
	second := 40
	//Every friday do once at 11:50:00(HH:mm:ss).
	robin.EveryTuesday().At(14, minute, second).Do(runCron, "Friday")

	//Every N day do once at 11:50:00(HH:mm:ss)
	robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

	//Every N hours do once at N:50:00(HH:mm:ss).
	robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Hours")

	//Every N minutes do once.
	robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Minutes")

	//Every N seconds do once.
	robin.Every(10).Seconds().Do(runCron, "Seconds")

	// Use a new cron executor
	//newCronDelay := cron.NewCronDelayExecutor()
	//newCronDelay.Delay(2000).Do(runCron, "newDelayCron Delay 2000 in Ms")

	// Use a new cron executor
	//newCronScheduler := cron.NewCronSchedulerExecutor()
	//newCronScheduler.Every(1).Hours().Do(runCron, "newCronScheduler Hours")
	//newCronScheduler.Every(1).Minutes().Do(runCron, "newCronScheduler Minutes")
	<-quitSemaphore
}

func runCron(s string) {
	log.Printf("I am %s CronTest %v\n", s, time.Now())
}





func RunChannelTest() {
    var channelThreadFiberPool = robin.NewGoroutineMulti()
    channelThreadFiberPool.Start()
    var channel1 = robin.NewChannel()
    var channel2 = robin.NewChannel()
    //channelThreadFiber.Start()
    channelThreadFiberPool.Start()
    channel1.Subscribe(channelThreadFiberPool, func() {
        log.Printf("我是訂閱者 1")
    })

    channel1.Subscribe(channelThreadFiberPool, func() {
        log.Printf("我是訂閱者 2")
    })
    //定時發布到  channel1
    dd := channelThreadFiberPool.ScheduleOnInterval(2000, 2000, channel1.Publish)

    channelThreadFiberPool.Schedule(2100, func() { dd.Dispose() })

    Subscribe1 := channel2.Subscribe(channelThreadFiberPool, subscribe)
    channel2.Publish(2, channel2.NumSubscribers(), "第一次發布 => Subscribe1")

    Subscribe2 := channel2.Subscribe(channelThreadFiberPool, subscribe)
    channel2.Publish(2, channel2.NumSubscribers(), "第二次發布 => Subscribe2")
    Subscribe1.Dispose()
    Subscribe2.Dispose()

    channel2.Publish(2, channel2.NumSubscribers(), "第三次發布")

    index := 0
    channel2.Subscribe(channelThreadFiberPool, func(msg string) {
        index++
        log.Printf("通道 2 訂閱人數 %d %v %s 發布次數 %d", channel2.NumSubscribers(), time.Now(), msg, index)
    })

    channelThreadFiberPool.ScheduleOnInterval(0, 10000, channel2.Publish, "定期發布")

}

func subscribe(channel int, numSubscribers int, msg string) {
    log.Printf("通道 %d 訂閱人數 %d 參數: %s %v", channel, numSubscribers, msg, time.Now())
}
