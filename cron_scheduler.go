package robin

import (
	"fmt"
	"time"
)

var schedulerExecutor = NewSchedulerExecutor()

type jobSchedulerExecutor struct {
	fiber *GoroutineMulti
}

func NewSchedulerExecutor() *jobSchedulerExecutor {
	return new(jobSchedulerExecutor).init()
}

func (c *jobSchedulerExecutor) init() *jobSchedulerExecutor {
	c.fiber = NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func Every(interval int64) *Job {
	return schedulerExecutor.Every(interval)
}

func (c *jobSchedulerExecutor) Every(interval int64) *Job {
	return NewJob(interval, c.fiber, delayNone)
}

type Job struct {
	fiber        Fiber
	identifyId   string
	loc          *time.Location
	task         task
	taskDisposer Disposable
	weekday      time.Weekday
	hour         int
	minute       int
	second       int
	unit         unit
	delayUnit    delayUnit
	interval     int64
	nextRunTime  time.Time
}

func (c *Job) init(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	c.hour = -1
	c.minute = -1
	c.second = -1
	c.fiber = fiber
	c.loc = time.Local
	c.interval = intervel
	c.delayUnit = delayUnit
	c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
	return c
}

func NewJob(intervel int64, fiber Fiber, delayUnit delayUnit) *Job {
	return new(Job).init(intervel, fiber, delayUnit)
}

type unit int
type delayUnit int

const (
	delay unit = iota
	weeks
	days
	hours
	minutes
	seconds
)

const (
	delayNone delayUnit = iota
	delayWeeks
	delayDays
	delayHours
	delayMinutes
	delaySeconds
	delayMilliseconds
)

func (c *Job) Dispose() {
	c.taskDisposer.Dispose()
	c.fiber = nil
}

func (c Job) Identify() string {
	return c.identifyId
}

func (c *Job) Days() *Job {
	if c.delayUnit == delayNone {
		c.unit = days
	} else {
		c.delayUnit = delayDays
	}
	return c
}

func (c *Job) Hours() *Job {
	if c.delayUnit == delayNone {
		c.unit = hours
	} else {
		c.delayUnit = delayHours
	}
	return c
}

func (c *Job) Minutes() *Job {
	if c.delayUnit == delayNone {
		c.unit = minutes
	} else {
		c.delayUnit = delayMinutes
	}
	return c
}

func (c *Job) Seconds() *Job {
	if c.delayUnit == delayNone {
		c.unit = seconds
	} else {
		c.delayUnit = delaySeconds
	}
	return c
}

func (c *Job) At(hour int, minute int, second int) *Job {
	c.hour = Abs(c.hour)
	c.minute = Abs(c.minute)
	c.second = Abs(c.second)
	if c.unit != hours {
		c.hour = hour % 24
	}
	c.minute = minute % 60
	c.second = second % 60
	return c
}

func (c *Job) Do(fun interface{}, params ...interface{}) Disposable {
	c.task = newTask(fun, params...)
	now := time.Now()
	switch c.unit {
	case delay:
		switch c.delayUnit {
		case delayWeeks:
			c.nextRunTime = now.AddDate(0, 0, 7)
		case delayDays:
			c.nextRunTime = now.AddDate(0, 0, int(c.interval))
		case delayHours:
			c.nextRunTime = now.Add(time.Duration(c.interval) * time.Hour)
		case delayMinutes:
			c.nextRunTime = now.Add(time.Duration(c.interval) * time.Minute)
		case delaySeconds:
			c.nextRunTime = now.Add(time.Duration(c.interval) * time.Second)
		case delayMilliseconds:
			c.nextRunTime = now.Add(time.Duration(c.interval) * time.Millisecond)
		}
	case weeks:
		i := (7 - (int(now.Weekday() - c.weekday))) % 7
		c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day()+int(i), c.hour, c.minute, c.second, 0, c.loc)
		if c.nextRunTime.Before(now) {
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
		}
	case days:
		if c.second < 0 || c.minute < 0 || c.hour < 0 {
			c.nextRunTime = now.AddDate(0, 0, 1)
			c.second = c.nextRunTime.Second()
			c.minute = c.nextRunTime.Minute()
			c.hour = c.nextRunTime.Hour()
		} else {
			c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), c.hour, c.minute, c.second, 0, c.loc)
			if c.interval > 1 {
				c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval-1))
			}
			if c.nextRunTime.Before(now) {
				c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
			}
		}
	case hours:
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
	case minutes:
		if c.second < 0 {
			c.second = now.Second()
		}
		c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), c.second, 0, c.loc)
		c.nextRunTime.Add(time.Duration(c.interval-1) * time.Hour)
		if c.nextRunTime.Before(now) {
			c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Minute /*.Duration(60*60*1000000)*/)
		}
	case seconds:
		c.nextRunTime = now.Add(time.Duration(c.interval) * time.Second)
	}

	firstInMs := int64(c.nextRunTime.Sub(now) / time.Millisecond)
	c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
	return c
}

func (c *Job) canDo() {
	now := time.Now()
	if now.After(c.nextRunTime) {
		c.fiber.EnqueueWithTask(c.task)
		switch c.unit {
		case delay:
			return
		case weeks:
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
		case days:
			c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
		case hours:
			c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Hour)
		case minutes:
			c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Minute)
		case seconds:
			c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Second)
		}
	} else {
		c.taskDisposer.Dispose()
	}
	adjustTime := int64(c.nextRunTime.Sub(now)) / 1000000
	if adjustTime < 1 {
		adjustTime = 1
	}
	c.taskDisposer = c.fiber.Schedule(adjustTime, c.canDo)
}
