# robin
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jiansoft/robin)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiansoft/robin)](https://goreportcard.com/report/github.com/jiansoft/robin)
[![build-test](https://github.com/jiansoft/robin/actions/workflows/go.yml/badge.svg)](https://github.com/jiansoft/robin/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/jiansoft/robin/branch/master/graph/badge.svg)](https://codecov.io/gh/jiansoft/robin)
[![](https://img.shields.io/github/tag/jiansoft/robin.svg)](https://github.com/jiansoft/robin/releases)
### Features

Fiber  
-------
* GoroutineSingle - a fiber backed by a dedicated goroutine. Every job is executed by a goroutine.
* GoroutineMulti - a fiber backed by more goroutine. Each job is executed by a new goroutine.

Channels
-------
* Channels callback is executed for each message received.

Cron
-------
Golang job scheduling for humans. It is inspired by [schedule](<https://github.com/dbader/schedule>).

Usage
================

### Quick Start

1.Install
~~~
  go get github.com/jiansoft/robin
~~~

2.Use examples
~~~ golang
import (
    "log"
    "time"
    
    "github.com/jiansoft/robin"
)

func main() {
    //The method is going to execute only once after 2000 ms.
    //指定的方法將在2000亳秒之後執行一次
    robin.Delay(2000).Do(runCron, "a Delay 2000 ms")
    
    minute := 11
    second := 50
    
    //It will execute once at 14:11:50 every Friday.
    //每星期五在14:11:50時將執行一次
    robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

    //It will execute once every N days at 14:11:50
    //每 N 天在14:11:50時將執行一次
    robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")
    
    //It will execute once at 11 minutes and 50 seconds past every N hours.
    //每過 N 小時在11分50秒時將執行一次
    robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Every 1 Hours")

    //It will execute once at the 50th second of every N minutes.
    //每過 N 分鐘在第50秒時將執行一次
    robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Every 1 Minutes")

    //Every N seconds is going to execute once
    //每過 N 秒將執行一次
    robin.Every(10).Seconds().Do(runCron, "Every 10 Seconds")
    
    p1 := player{Nickname: "Player 1"}
    p2 := player{Nickname: "Player 2"}
    p3 := player{Nickname: "Player 3"}
    p4 := player{Nickname: "Player 4"}
    
    //Create a channel
    //建立一個頻道
    channel := robin.NewChannel()
    
    //Four player subscribe the channel
    //建立四個玩家訂閱頻道
    channel.Subscribe(p1.eventFinalBossResurge)
    channel.Subscribe(p2.eventFinalBossResurge)
    p3Subscribe := channel.Subscribe(p3.eventFinalBossResurge)
    p4Subscribe := channel.Subscribe(p4.eventFinalBossResurge)
    
    //Publish a message to the channel and then the four subscribers of the channel will 
    //receives the message each that "The boss resurge first." .
    //發布一則訊息到頻道中然後頻道中的四個訂閱者將會各自收到一則"魔王首次復活"的訊息
    channel.Publish("The boss resurge first.")
    
    //Unsubscribe p3 and p4 from the channel.
    //取消第3位和第4位玩家在該頻道的訂閱
    channel.Unsubscribe(p3Subscribe)
    p4Subscribe.Unsubscribe()
    
    //This time just p1 and p2 receives the message that "The boss resurge second.".
    //這次發布的訊息將只有第1位與第2位玩家各自收到一則"魔王第二次復活"的訊息
    channel.Publish("The boss resurge second.")
    
    //Unsubscribe all subscribers from the channel
    //取消頻道內的所有訂閱者 
    channel.Clear()
    
    //The channel is empty so no one can receive the message
    //頻道內沒有任何訂閱者，所以沒有玩家會收到"魔王第三次復活"的訊息
    channel.Publish("The boss resurge third.")
}

func runCron(s string) {
    log.Printf("I am %s CronTest %v\n", s, time.Now())
}

type player struct {
	Nickname string
}
func (p player) eventFinalBossResurge(someBossInfo string) {
	log.Printf("%s receive a message : %s", p.Nickname, someBossInfo)
}
~~~
[More example](<https://github.com/jiansoft/robin/blob/master/example/main.go>)

## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_large)