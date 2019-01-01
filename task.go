package robin

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	taskPool = sync.Pool{
		New: func() interface{} { return Task{} },
	}

	timerTaskPool = sync.Pool{
		New: func() interface{} {
			t := new(timerTask)
			t.lock = &sync.Mutex{}
			return t // new(timerTask)
		},
	}
)

//Task a struct
type Task struct {
	doFunc      interface{}
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func newTask(t interface{}, p ...interface{}) Task {
	tmp := taskPool.Get()
	task := tmp.(Task)
	task.doFunc = t
	//task := Task{doFunc: t}
	task.funcCache = reflect.ValueOf(t)
	task.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t Task) run() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}

func (t Task) release() {
	taskPool.Put(t)
}

type timerTask struct {
	identifyID   string
	scheduler    SchedulerRegistry
	firstInMs    int64
	intervalInMs int64
	//first        *time.Timer
	//interval     *time.Ticker
	task      Task
	cancelled bool
	lock      *sync.Mutex
}

func newTimerTask(fiber SchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	timerTask := timerTaskPool.Get().(*timerTask)
	return timerTask.init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) init(scheduler SchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	//t.lock.Lock()
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	//t.lock.Unlock()
	t.identifyID = fmt.Sprintf("%p-%p", &t, &task)
	//t.setIdentifyID(fmt.Sprintf("%p-%p", &t, &task))
	t.setCancelled(false)
	return t
}

func (t *timerTask) Dispose() {
	if t.cancelled {
		return
	}
	t.setCancelled(true)
	//t.lock.Lock()
	//defer t.lock.Unlock()
	//t.cancelled = true

	if nil != t.scheduler {
		t.scheduler.Remove(t)
	}

	t.task.release()

	timerTaskPool.Put(t)
}

func (t *timerTask) Identify() string {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.identifyID
}

func (t *timerTask) schedule() {
	if t.firstInMs <= 0 {
		t.doFirstSchedule()
		return
	}
	//t.lock.Lock()
	first := time.NewTimer(time.Duration(t.firstInMs) * time.Millisecond)
	//t.lock.Unlock()
	go func() {
		select {
		case <-first.C:
			t.doFirstSchedule()
		}
	}()
}

func (t *timerTask) doFirstSchedule() {
	t.executeOnFiber()
	t.doIntervalSchedule()
}

func (t *timerTask) doIntervalSchedule() {
	if t.intervalInMs <= 0 {
		t.Dispose()
		return
	}
	//t.lock.Lock()
	interval := time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	//t.lock.Unlock()
	go func() {
		for !t.cancelled {
			/*select {
			case <-t.interval.C:
				t.executeOnFiber()
			}*/
			<-interval.C
			t.executeOnFiber()
		}
		interval.Stop()
	}()
}

func (t timerTask) executeOnFiber() {
	if t.cancelled {
		return
	}
	t.scheduler.EnqueueWithTask(t.task)
}

/*func (t timerTask) isCancelled() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.cancelled
}*/

func (t *timerTask) setCancelled(r bool) {
	t.lock.Lock()
	t.cancelled = r
	t.lock.Unlock()
}

/*func (t *timerTask) setIdentifyID(id string ) {
	t.lock.Lock()
	t.identifyID = id
	t.lock.Unlock()
}
*/
