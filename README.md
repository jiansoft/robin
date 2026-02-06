# robin

**English** | [繁體中文](README_zh-TW.md)

[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jiansoft/robin)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiansoft/robin)](https://goreportcard.com/report/github.com/jiansoft/robin)
[![build-test](https://github.com/jiansoft/robin/actions/workflows/go.yml/badge.svg)](https://github.com/jiansoft/robin/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/jiansoft/robin/branch/master/graph/badge.svg)](https://codecov.io/gh/jiansoft/robin)
[![](https://img.shields.io/github/tag/jiansoft/robin.svg)](https://github.com/jiansoft/robin/releases)

A Go library providing **Fiber (actor model)**, **Cron (job scheduling)**, **Channel (pub/sub)**, and **Concurrent
Collections** — all in a single package with zero external dependencies.

Requires **Go 1.22+**.

## Features

### Fiber

Task execution fibers backed by goroutines.

- **GoroutineSingle** — dedicated goroutine, tasks execute serially in order.
- **GoroutineMulti** — each task batch spawns new goroutines for concurrent execution.

```go
gm := robin.NewGoroutineMulti()
defer gm.Dispose()

gm.Enqueue(func (msg string) {
fmt.Println(msg)
}, "hello")

gm.Schedule(1000, func () {
fmt.Println("executed after 1 second")
})

d := gm.ScheduleOnInterval(0, 500, func () {
fmt.Println("every 500ms")
})
// d.Dispose() to cancel
```

#### Type-safe Enqueue

Generic helpers with compile-time type checking (no reflection for zero-arg):

```go
robin.Enqueue0(fiber, func () { fmt.Println("no args") })
robin.Enqueue1(fiber, func (s string) { fmt.Println(s) }, "hello")
robin.Enqueue2(fiber, func (a int, b string) { fmt.Println(a, b) }, 42, "world")
robin.Enqueue3(fiber, func (a int, b string, c bool) { fmt.Println(a, b, c) }, 1, "x", true)
```

### Cron

Human-friendly job scheduling with a fluent builder API. Inspired by [schedule](https://github.com/dbader/schedule).

```go
// Execute immediately
robin.RightNow().Do(task, args...)

// Execute once after delay
robin.Delay(2000).Do(task, args...)

// Execute up to N times at interval
robin.Delay(1000).Times(3).Do(task, args...)

// Every N milliseconds / seconds / minutes / hours / days
robin.Every(500).Milliseconds().Do(task, args...)
robin.Every(10).Seconds().Do(task, args...)
robin.Every(5).Minutes().Do(task, args...)
robin.Every(2).Hours().At(0, 30, 0).Do(task, args...) // at mm:ss
robin.Every(1).Days().At(8, 0, 0).Do(task, args...)     // at HH:mm:ss
robin.Everyday().At(23, 59, 59).Do(task, args...)

// Weekly scheduling
robin.EveryMonday().At(9, 0, 0).Do(task, args...)
robin.EveryFriday().At(14, 30, 0).Do(task, args...)
// Also: EveryTuesday, EveryWednesday, EveryThursday, EverySaturday, EverySunday

// Execute once at a specific time
robin.Until(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)).Do(task, args...)

// Limit execution count
robin.Every(1).Seconds().Times(10).Do(task, args...)

// Restrict to a time range (HH:mm:ss)
from := time.Date(0, 0, 0, 9, 0, 0, 0, time.Local)
to := time.Date(0, 0, 0, 17, 0, 0, 0, time.Local)
robin.Every(30).Seconds().Between(from, to).Do(task, args...)

// Calculate next time after task completes (for long-running tasks)
robin.Every(10).Seconds().AfterExecuteTask().Do(task, args...)

// Cancel a scheduled job
job := robin.Every(1).Seconds().Do(task, args...)
job.Dispose()
```

#### CronScheduler (Instance-based)

For independent lifecycle management — each scheduler owns its own fiber:

```go
cs := robin.NewCronScheduler()
cs.Every(5).Seconds().Do(task, args...)
cs.Delay(1000).Do(task, args...)
cs.RightNow().Do(task, args...)
cs.Everyday().At(12, 0, 0).Do(task, args...)
cs.EveryMonday().At(9, 0, 0).Do(task, args...)
cs.Until(targetTime).Do(task, args...)
cs.Dispose() // stops all jobs on this scheduler
```

### Channel (Pub/Sub)

Thread-safe publish-subscribe messaging:

```go
channel := robin.NewChannel()

// Subscribe — returns a *Subscriber
sub1 := channel.Subscribe(func (msg string) {
fmt.Println("received:", msg)
})
sub2 := channel.Subscribe(handler)

// Publish to all subscribers
channel.Publish("hello")

// Unsubscribe
sub1.Unsubscribe() // via subscriber
channel.Unsubscribe(sub2)   // via channel

// Get subscriber count
fmt.Println(channel.Count())

// Remove all subscribers
channel.Clear()
```

### Concurrent Collections

Generic, thread-safe collections with zero external dependencies.

#### ConcurrentQueue[T] (FIFO)

Ring buffer implementation with automatic grow/shrink:

```go
q := robin.NewConcurrentQueue[string]()
q.Enqueue("a")
q.Enqueue("b")

val, ok := q.TryPeek() // peek without removing
val, ok = q.TryDequeue() // remove and return
fmt.Println(q.Len()) // element count (lock-free)
arr := q.ToArray() // copy to slice
q.Clear()          // remove all
```

#### ConcurrentStack[T] (LIFO)

```go
s := robin.NewConcurrentStack[int]()
s.Push(10)
s.Push(20)

val, ok := s.TryPeek() // peek at top
val, ok = s.TryPop() // pop from top
arr := s.ToArray() // LIFO order (top first)
s.Clear()
```

#### ConcurrentBag[T] (Unordered)

```go
b := robin.NewConcurrentBag[float64]()
b.Add(3.14)

val, ok := b.TryTake() // remove an element
arr := b.ToArray()
b.Clear()
```

### Utility

```go
robin.Abs(-42) // 42 (int)
robin.Abs(-3.14) // 3.14 (float64)
```

## Installation

```
go get github.com/jiansoft/robin
```

## Full Example

See [example/main.go](https://github.com/jiansoft/robin/blob/master/example/main.go) for a comprehensive, self-verifying
example that demonstrates and validates every public API.

## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjiansoft%2Frobin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjiansoft%2Frobin?ref=badge_large)
