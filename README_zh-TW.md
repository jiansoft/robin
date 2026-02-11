# robin

[English](README.md) | **繁體中文**

[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jiansoft/robin)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiansoft/robin)](https://goreportcard.com/report/github.com/jiansoft/robin)
[![build-test](https://github.com/jiansoft/robin/actions/workflows/go.yml/badge.svg)](https://github.com/jiansoft/robin/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/jiansoft/robin/branch/master/graph/badge.svg)](https://codecov.io/gh/jiansoft/robin)
[![](https://img.shields.io/github/tag/jiansoft/robin.svg)](https://github.com/jiansoft/robin/releases)

Go 函式庫，提供 **Fiber（Actor 模型）**、**Cron（任務排程）**、**Channel（發布/訂閱）**、**Concurrent Collections（執行緒安全集合）
**— 單一套件、零外部依賴。

需要 **Go 1.22+**。

## 功能特色

### Fiber

以 goroutine 驅動的任務執行纖程。

- **GoroutineSingle** — 單一專屬 goroutine，所有任務依序串行執行。
- **GoroutineMulti** — 每個任務批次由新的 goroutine 併發執行。

```go
gm := robin.NewGoroutineMulti()
defer gm.Dispose()

gm.Enqueue(func (msg string) {
fmt.Println(msg)
}, "hello")

gm.Schedule(1000, func () {
fmt.Println("1 秒後執行")
})

d := gm.ScheduleOnInterval(0, 500, func () {
fmt.Println("每 500ms 執行一次")
})
// d.Dispose() 取消排程
```

### Cron

人性化的任務排程，使用流式建構器 API。靈感來自 [schedule](https://github.com/dbader/schedule)。

```go
// 立即執行
robin.RightNow().Do(task, args...)

// 延遲 N 毫秒後執行一次
robin.Delay(2000).Do(task, args...)

// 延遲後重複執行，最多 N 次
robin.Delay(1000).Times(3).Do(task, args...)

// 每 N 毫秒 / 秒 / 分 / 時 / 天
robin.Every(500).Milliseconds().Do(task, args...)
robin.Every(10).Seconds().Do(task, args...)
robin.Every(5).Minutes().Do(task, args...)
robin.Every(2).Hours().At(0, 30, 0).Do(task, args...) // 在 mm:ss 時
robin.Every(1).Days().At(8, 0, 0).Do(task, args...)     // 在 HH:mm:ss 時
robin.Everyday().At(23, 59, 59).Do(task, args...)

// 每週排程
robin.EveryMonday().At(9, 0, 0).Do(task, args...)
robin.EveryFriday().At(14, 30, 0).Do(task, args...)
// 另有：EveryTuesday, EveryWednesday, EveryThursday, EverySaturday, EverySunday

// 在指定時間點執行一次
robin.Until(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)).Do(task, args...)

// 限制執行次數
robin.Every(1).Seconds().Times(10).Do(task, args...)

// 限制在指定時間區間內執行（HH:mm:ss）
from := time.Date(0, 0, 0, 9, 0, 0, 0, time.Local)
to := time.Date(0, 0, 0, 17, 0, 0, 0, time.Local)
robin.Every(30).Seconds().Between(from, to).Do(task, args...)

// 等待任務完成後才計算下次執行時間（適用於耗時任務）
robin.Every(10).Seconds().AfterExecuteTask().Do(task, args...)

// 取消已排程的工作
job := robin.Every(1).Seconds().Do(task, args...)
job.Dispose()
```

#### CronScheduler

獨立生命週期管理 — 每個排程器擁有自己的 Fiber：

```go
cs := robin.NewCronScheduler()
cs.Every(5).Seconds().Do(task, args...)
cs.Delay(1000).Do(task, args...)
cs.RightNow().Do(task, args...)
cs.Everyday().At(12, 0, 0).Do(task, args...)
cs.EveryMonday().At(9, 0, 0).Do(task, args...)
cs.Until(targetTime).Do(task, args...)
cs.Dispose() // 停止此排程器上所有工作
```

### Channel（發布/訂閱）

執行緒安全的發布/訂閱訊息頻道：

```go
channel := robin.NewChannel()

// Subscribe 訂閱 — 回傳 *Subscriber
sub1 := channel.Subscribe(func (msg string) {
fmt.Println("收到:", msg)
})
sub2 := channel.Subscribe(handler)

// Publish 發布訊息給所有訂閱者
channel.Publish("hello")

// 取消訂閱
sub1.Unsubscribe() // 透過訂閱者物件
channel.Unsubscribe(sub2)   // 透過頻道方法

// 取得訂閱者數量
fmt.Println(channel.Count())

// 移除所有訂閱者
channel.Clear()
```

#### TypedChannel（泛型，無反射）

型別安全的泛型頻道 — 不經反射，效能更好：

```go
ch := robin.NewTypedChannel[string]()

sub := ch.Subscribe(func (msg string) {
fmt.Println("收到:", msg)
})

ch.Publish("hello")

sub.Unsubscribe()
fmt.Println(ch.Count())
ch.Clear()
```

支援任意型別，包含結構體：

```go
type Event struct {
Name string
Code int
}

ch := robin.NewTypedChannel[Event]()
ch.Subscribe(func (e Event) {
fmt.Printf("事件: %s (%d)\n", e.Name, e.Code)
})
ch.Publish(Event{Name: "click", Code: 42})
```

使用 `NewTypedChannelWithFiber` 指定自訂 Fiber：

```go
f := robin.NewGoroutineSingle()
defer f.Dispose()
ch := robin.NewTypedChannelWithFiber[string](f)
```

### Concurrent Collections（執行緒安全集合）

泛型、執行緒安全的集合，零外部依賴。

#### ConcurrentQueue[T]（FIFO 佇列）

環形緩衝區實作，自動擴容/縮容：

```go
q := robin.NewConcurrentQueue[string]()
q.Enqueue("a")
q.Enqueue("b")

val, ok := q.TryPeek() // 查看前端元素但不移除
val, ok = q.TryDequeue() // 取出並移除前端元素
fmt.Println(q.Len()) // 元素數量（無鎖，使用 atomic）
arr := q.ToArray() // 複製到 slice
q.Clear()          // 清除所有元素
```

#### ConcurrentStack[T]（LIFO 堆疊）

```go
s := robin.NewConcurrentStack[int]()
s.Push(10)
s.Push(20)

val, ok := s.TryPeek() // 查看頂端
val, ok = s.TryPop() // 彈出頂端
arr := s.ToArray() // LIFO 順序（頂端在前）
s.Clear()
```

#### ConcurrentBag[T]（無序集合）

```go
b := robin.NewConcurrentBag[float64]()
b.Add(3.14)

val, ok := b.TryTake() // 取出一個元素
arr := b.ToArray()
b.Clear()
```

### PanicHandler（Panic 處理）

預設情況下，任務 panic 時 robin 會攔截並印到 stderr — Fiber 繼續處理後續任務。可自訂或停用此行為：

```go
// 自訂處理：導向自己的 logger
robin.SetPanicHandler(func (r any, stack []byte) {
log.Printf("task panic: %v\n%s", r, stack)
})

// 停用攔截：讓 panic 直接終止程式（Go 預設行為）
robin.SetPanicHandler(nil)
```

### 工具函式

```go
robin.Abs(-42) // 42 (int)
robin.Abs(-3.14) // 3.14 (float64)
```

## 安裝

```
go get github.com/jiansoft/robin/v2
```

## 完整範例

參閱 [example/main.go](https://github.com/jiansoft/robin/blob/master/example/main.go) — 包含所有公開 API
的自我驗證範例，可直接執行驗證所有功能是否正確運作。

## 授權條款

Copyright (c) 2017

以 MIT 授權條款發布：

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_large)
