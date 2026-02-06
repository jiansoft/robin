# robin

[English](README.md) | **繁體中文**

[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jiansoft/robin)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiansoft/robin)](https://goreportcard.com/report/github.com/jiansoft/robin)
[![build-test](https://github.com/jiansoft/robin/actions/workflows/go.yml/badge.svg)](https://github.com/jiansoft/robin/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/jiansoft/robin/branch/master/graph/badge.svg)](https://codecov.io/gh/jiansoft/robin)
[![](https://img.shields.io/github/tag/jiansoft/robin.svg)](https://github.com/jiansoft/robin/releases)

### 功能特色

Fiber（纖程）
-------

* GoroutineSingle - 由單一專屬 goroutine 驅動的 fiber，所有任務依序在同一個 goroutine 中執行。
* GoroutineMulti - 由多個 goroutine 驅動的 fiber，每個任務由一個新的 goroutine 執行。

Channels（頻道）
-------

* 頻道的回呼函式會在每次收到訊息時被執行。

Cron（排程）
-------
為 Golang 設計的人性化任務排程，靈感來自 [schedule](<https://github.com/dbader/schedule>)。

使用方式
================

### 快速開始

1.安裝

~~~
  go get github.com/jiansoft/robin
~~~

2.使用範例

~~~ golang
import (
    "log"
    "time"

    "github.com/jiansoft/robin"
)

func main() {
    // 指定的方法將在 2000 毫秒之後執行一次
    robin.Delay(2000).Do(runCron, "a Delay 2000 ms")

    minute := 11
    second := 50

    // 每星期五在 14:11:50 時執行一次
    robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

    // 每 N 天在 14:11:50 時執行一次
    robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

    // 每過 N 小時在 11 分 50 秒時執行一次
    robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Every 1 Hours")

    // 每過 N 分鐘在第 50 秒時執行一次
    robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Every 1 Minutes")

    // 每過 N 秒執行一次
    robin.Every(10).Seconds().Do(runCron, "Every 10 Seconds")

    p1 := player{Nickname: "Player 1"}
    p2 := player{Nickname: "Player 2"}
    p3 := player{Nickname: "Player 3"}
    p4 := player{Nickname: "Player 4"}

    // 建立一個頻道
    channel := robin.NewChannel()

    // 四個玩家訂閱頻道
    channel.Subscribe(p1.eventFinalBossResurge)
    channel.Subscribe(p2.eventFinalBossResurge)
    p3Subscribe := channel.Subscribe(p3.eventFinalBossResurge)
    p4Subscribe := channel.Subscribe(p4.eventFinalBossResurge)

    // 發布一則訊息到頻道，四個訂閱者將各自收到「魔王首次復活」的訊息
    channel.Publish("The boss resurge first.")

    // 取消第 3 位和第 4 位玩家的頻道訂閱
    channel.Unsubscribe(p3Subscribe)
    p4Subscribe.Unsubscribe()

    // 這次只有第 1 位與第 2 位玩家會各自收到「魔王第二次復活」的訊息
    channel.Publish("The boss resurge second.")

    // 取消頻道內的所有訂閱者
    channel.Clear()

    // 頻道內沒有任何訂閱者，所以沒有玩家會收到「魔王第三次復活」的訊息
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

[更多範例](<https://github.com/jiansoft/robin/blob/master/example/main.go>)

## 授權條款

Copyright (c) 2017

以 MIT 授權條款發布：

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_large)
