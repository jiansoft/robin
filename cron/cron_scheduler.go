package cron

import (
	"fmt"
	"time"

	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/core"
	"github.com/jiansoft/robin/fiber"
	"github.com/jiansoft/robin/math"
)

var schedulerExecutor = NewSchedulerExecutor()

type jobSchedulerExecutor struct {
	fiber *fiber.GoroutineMulti
}

func NewSchedulerExecutor() *jobSchedulerExecutor {
	return new(jobSchedulerExecutor).init()
}

func (c *jobSchedulerExecutor) init() *jobSchedulerExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func Every(interval int64) *Job {
	return schedulerExecutor.Every(interval)
}

func (c *jobSchedulerExecutor) Every(interval int64) *Job {
	return NewJob(interval, c.fiber)
}

type Job struct {
	fiber        fiber.Fiber
	identifyId     string
	loc          *time.Location
	task         core.Task
	taskDisposer robin.Disposable
	weekday      time.Weekday
	hour         int
	minute       int
	second       int
	unit         unit
	interval     int64
	nextRunTime  time.Time
}

func (c *Job) init(intervel int64, fiber fiber.Fiber) *Job {
	c.hour = -1
	c.minute = -1
	c.second = -1
	c.fiber = fiber
	c.loc = time.Local
	c.interval = intervel
	c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
	return c
}

func NewJob(intervel int64, fiber fiber.Fiber) *Job {
	return new(Job).init(intervel, fiber)
}

type unit int

// iota 初始化後會自動遞增
const (
	delay unit = iota
	weeks
	days
	hours
	minutes
	seconds
)

func (c *Job) Dispose() {
	c.taskDisposer.Dispose()
	c.fiber = nil
}

func (c Job) Identify() string {
	return c.identifyId
}

func (c *Job) Days() *Job {
	c.unit = days
	return c
}

func (c *Job) Hours() *Job {
	c.unit = hours
	return c
}

func (c *Job) Minutes() *Job {
	c.unit = minutes
	return c
}

func (c *Job) Seconds() *Job {
	c.unit = seconds
	return c
}

func (c *Job) At(hour int, minute int, second int) *Job {
	c.hour = math.Abs(c.hour)
	c.minute = math.Abs(c.minute)
	c.second = math.Abs(c.second)
	if c.unit != hours {
		c.hour = hour % 24
	}
	c.minute = minute % 60
	c.second = second % 60
	return c
}

func (c *Job) Do(fun interface{}, params ...interface{}) robin.Disposable {
	c.task = core.NewTask(fun, params...)
	firstInMs := int64(0)
	now := time.Now()
	switch c.unit {
	case weeks:
		i := (7 - (int(now.Weekday() - c.weekday))) % 7
		c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day()+int(i), c.hour, c.minute, c.second, 0, c.loc)
		if c.nextRunTime.Before(now) {
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
		}
		firstInMs = int64(c.nextRunTime.Sub(now)) / 1000000
		c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
	case days:
		if c.second < 0 || c.minute < 0 || c.hour < 0 {
			c.nextRunTime = now.AddDate(0, 0, 1)
			c.second = c.nextRunTime.Second()
			c.minute = c.nextRunTime.Minute()
			c.hour = c.nextRunTime.Hour()
			firstInMs = int64(c.nextRunTime.Sub(now)) / 1000000
		} else {
			c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), c.hour, c.minute, c.second, 0, c.loc)
			if c.interval > 1 {
				c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval-1))
			}
			if now.After(c.nextRunTime) {
				c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
			}
			firstInMs = int64(c.nextRunTime.Sub(now)) / 1000000
		}
		c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
	case hours:
		firstInMs = c.interval * 60 * 60 * 1000
		if c.minute < 0 {
			c.minute = now.Minute()
		}
		if c.second < 0 {
			c.second = now.Second()
		}
		c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), c.minute, c.second, 0, c.loc)
		c.nextRunTime.Add(time.Duration(c.interval-1) * time.Hour)
		if c.nextRunTime.Before(now) {
			c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Hour /*.Duration(60*60*1000000)*/)
		}
		firstInMs = int64(c.nextRunTime.Sub(now)) / 1000000
		c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
	case minutes:
		firstInMs = c.interval * 60 * 1000
		c.taskDisposer = c.fiber.ScheduleOnInterval(firstInMs, firstInMs, c.task.Run)
	case seconds:
		firstInMs = c.interval * 1000
		c.taskDisposer = c.fiber.ScheduleOnInterval(firstInMs, firstInMs, c.task.Run)
	case delay:
		c.taskDisposer = c.fiber.Schedule(c.interval, c.task.Run)
	}
	return c
}

func (c *Job) canDo() {
	now := time.Now()
	//var adjustTime int64
	if now.After(c.nextRunTime) {
		c.fiber.EnqueueWithTask(c.task)
		switch c.unit {
		case weeks:
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
			//adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
		case days:
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
			// adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
		case hours:
			c.nextRunTime = c.nextRunTime.Add(time.Hour)
			// adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
		}
	} else {
		c.taskDisposer.Dispose()
		//adjustTime = int64(c.nextRunTime.Sub(now)) / 1000000
		//if adjustTime < 1 {
		//    adjustTime = 1
		//}
	}
	adjustTime := int64(c.nextRunTime.Sub(now)) / 1000000
	if adjustTime < 1 {
		adjustTime = 1
	}
	c.taskDisposer = c.fiber.Schedule(adjustTime, c.canDo)
}
