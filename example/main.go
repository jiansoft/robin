package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jiansoft/robin"
)

//var quitSemaphore chan bool

func main() {
	//RunChannelTest()

	var runCronFiber = robin.NewGoroutineSingle()
	runCronFiber.Start()
	var runCronFiber2 = robin.NewGoroutineSingle()
	runCronFiber2.Start()

	runCronFiber.Schedule(0, func() {
		log.Printf("just test\n")
		//	a.Dispose()
	})

	robin.Delay(2000).Do(runCron, "a Delay 2000 ms")
	//robin.Every(1000).MilliSeconds().Do(runCron, "Every(1000).MilliSeconds()")
	//robin.Every(2).Seconds().Do(runCron, "Every(2).Seconds()")
	//Every N seconds do once.
	//robin.Every(10).Seconds().Do(runCron, "Every 10 Seconds")
	runCronFiber.ScheduleOnInterval(10000, 10000, func() {
		log.Printf("runCronFiber 1\n")
	})
	runCronFiber2.ScheduleOnInterval(1000, 1000, func() {
		//log.Printf("runCronFiber 2\n")
	})
	runCronFiber2.ScheduleOnInterval(1000, 1000, func() {
		//log.Printf("runCronFiber 3\n")
	})

	//fmt.Scanln()
	//<-quitSemaphore
	//<-quitSemaphore
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

	//Every friday do once at 11:50:00(HH:mm:ss).
	robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

	//Every N day do once at 11:50:00(HH:mm:ss)
	robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

	//Every N hours do once at N:50:00(HH:mm:ss).
	robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Every 1 Hours")

	//Every N minutes do once.
	robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Every 1 Minutes")

	//Every N seconds do once.
	robin.Every(100).Seconds().Do(runCron, "Every 10 Seconds")

	robin.Every(1).Seconds().Times(3).Do(runCron, "Every 1 Seconds Times 3")
	robin.Delay(2000).Times(2).Do(runCron, " Delay 2000 ms Times 2")
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

func runCron1() {
	log.Printf("康中文測試 %s", time.Now())
}

func RunChannelTest() {
	var channelThreadFiberPool = robin.NewGoroutineMulti()
	channelThreadFiberPool.Start()
	var channel1 = robin.NewChannel()
	var channel2 = robin.NewChannel()
	//channelThreadFiber.Start()
	channelThreadFiberPool.Start()
	channel1.Subscribe(channelThreadFiberPool, func() {
		log.Printf("我是 channel1 的訂閱者 1")
	})

	channel1.Subscribe(channelThreadFiberPool, func() {
		log.Printf("我是 channel1 的訂閱者 2")
	})
	//定時2秒發布到  channel1
	dd := channelThreadFiberPool.ScheduleOnInterval(2000, 2000, channel1.Publish)
	//2.1 秒後取消定期的發布
	channelThreadFiberPool.Schedule(2100, func() { dd.Dispose() })

	Subscribe1 := channel2.Subscribe(channelThreadFiberPool, subscribe)
	channel2.Publish(2, channel2.NumSubscribers(), "第一次發布 => Subscribe1")

	Subscribe2 := channel2.Subscribe(channelThreadFiberPool, subscribe)
	channel2.Publish(2, channel2.NumSubscribers(), "第二次發布 => Subscribe2")
	//取消訂閱者 1 的訂閱
	Subscribe1.Dispose()
	//取消訂閱者 2 的訂閱
	Subscribe2.Dispose()
	//channel2 內目前沒有任何的訂閱者，不會有訊息輸出
	channel2.Publish(2, channel2.NumSubscribers(), "第三次發布")

	index := 0
	//重新加入一個訂閱者
	channel2.Subscribe(channelThreadFiberPool, func(msg string) {
		index++
		log.Printf("channe2 訂閱人數 %d %s 發布次數 %d %v ", channel2.NumSubscribers(), msg, index, time.Now().Format(time.RFC3339))
	})
	robin.Every(10).Seconds().Do(channel2.Publish, "定期發布")
	//channelThreadFiberPool.ScheduleOnInterval(0, 10000, channel2.Publish, "定期發布")

}

func subscribe(channel int, numSubscribers int, msg string) {
	log.Printf("通道 %d 訂閱人數 %d 參數: %s %v", channel, numSubscribers, msg, time.Now().Format(time.RFC3339))
}

func CronTestAndSleepASecond(s string){
	log.Printf("I am %s CronTest and sleep a second %v\n", s, time.Now())
	timeout := time.NewTimer(time.Duration(1000) * time.Millisecond)
	select {
	case <-timeout.C:
	}
}