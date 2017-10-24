package main

import (
	"log"
	"time"

	"github.com/jiansoft/robin/cron"
	"github.com/jiansoft/robin/fiber"
)

var quitSemaphore chan bool

func main() {
	log.Printf("Start\n")
	var runCronFiber = fiber.NewGoroutineMulti()
	runCronFiber.Start()
	_ = cron.Delay(2000).Do(runCron, "a Delay 2000 ms")

	//Make a cron that will be executed everyday at 15:30:04(HH:mm:ss)
	b := cron.Every(1).Days().At(15, 30, 4).Do(runCron, "Will be cancel")

	//After one second, a & b crons well be canceled by fiber
	runCronFiber.Schedule(1000, func() {
		log.Printf("a & b are dispose\n")
		b.Dispose()
		//	a.Dispose()
	})
	minute := 41
	second := 40
	//Every friday do once at 11:50:00(HH:mm:ss).
	cron.EveryTuesday().At(14, minute, second).Do(runCron, "Friday")

	//Every N day do once at 11:50:00(HH:mm:ss)
	cron.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

	//Every N hours do once at N:50:00(HH:mm:ss).
	cron.Every(1).Hours().At(0, minute, second).Do(runCron, "Hours")

	//Every N minutes do once.
	cron.Every(1).Minutes().At(0, 0, second).Do(runCron, "Minutes")

	//Every N seconds do once.
	cron.Every(10).Seconds().Do(runCron, "Seconds")

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
