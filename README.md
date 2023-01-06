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
    robin.Delay(2000).Do(runCron, "a Delay 2000 ms")
    
    minute := 11
    second := 50
    
    //Every Friday is going to execute once at 14:11:50 (HH:mm:ss).
    robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

    //Every N day  is going to execute once at 14:11:50(HH:mm:ss)
    robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

    //Every N hours is going to execute once at 11:50:00(HH:mm:ss).
    robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Every 1 Hours")

    //Every N minutes is going to execute once at 50(ss).
    robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Every 1 Minutes")

    //Every N seconds is going to execute once
    robin.Every(10).Seconds().Do(runCron, "Every 10 Seconds")
    
    p1 := player{Nickname: "Player 1"}
    p2 := player{Nickname: "Player 2"}
    p3 := player{Nickname: "Player 3"}
    p4 := player{Nickname: "Player 4"}
    
    //Create a channel
    channel := robin.NewChannel()
    
    //Four player subscribe the channel
    channel.Subscribe(p1.eventFinalBossResurge)
    channel.Subscribe(p2.eventFinalBossResurge)
    p3Subscribe := channel.Subscribe(p3.eventFinalBossResurge)
    p4Subscribe := channel.Subscribe(p4.eventFinalBossResurge)
    
    //Publish a message to the channel and then the four subscribers of the channel will 
    //receives the message each that "The boss resurge first." .
    channel.Publish("The boss resurge first.")
    
    //Unsubscribe p3 and p4 from the channel.
    channel.Unsubscribe(p3Subscribe)
    p4Subscribe.Unsubscribe()
    
    //This time just p1 and p2 receives the message that "The boss resurge second.".
    channel.Publish("The boss resurge second.")
    
    //Unsubscribe all subscribers from the channel
    channel.Clear()
    
    //The channel is empty so no one can receive the message
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