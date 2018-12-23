package robin

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestGoroutineMulti_init(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		subscriptions  *Disposer
		flushPending   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *GoroutineMulti
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineMulti{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				subscriptions:  tt.fields.subscriptions,
				flushPending:   tt.fields.flushPending,
			}
			if got := g.init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoroutineMulti.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGoroutineMulti(t *testing.T) {
	tests := []struct {
		name string
		want *GoroutineMulti
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGoroutineMulti(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoroutineMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoroutineMulti_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineMulti_Start"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineMulti()
			g.Start()
			g.Start()
		})
	}
}

func TestGoroutineMulti_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineMulti_Stop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineMulti()
			g.Start()
			g.Stop()
		})
	}
}

func TestGoroutineMulti_Dispose(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineMulti_Dispose"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineMulti()
			g.Dispose()
		})
	}
}

func TestGoroutineMulti_Enqueue(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name  string
		fiber Fiber
		args  interface{}
	}{
		{"Test_GoroutineMulti_Enqueue", g, func(s string) { t.Logf("s:%v", s) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.Enqueue(tt.args, "Test 1")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 2")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 3")
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineMulti_EnqueueWithTask(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name  string
		fiber Fiber
		args  Task
	}{
		{"Test_GoroutineMulti_EnqueueWithTask", g, newTask(func(s string) { t.Logf("s:%v", s) }, "Test 1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.EnqueueWithTask(tt.args)
			tt.fiber.EnqueueWithTask(newTask(func(s string) { t.Logf("s:%v", s) }, "Test 2"))
			tt.fiber.EnqueueWithTask(newTask(func(s string) { t.Logf("s:%v", s) }, "Test 3"))
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineMulti_Schedule(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name    string
		fiber   Fiber
		args    Task
		firstMs int64
	}{
		{"Test_GoroutineMulti_Schedule", g, newTask(func(s string) { t.Logf("s:%v", s) }), 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 1")
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 2")
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 3")
			gotD := tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 4")
			switch gotD.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineMulti.Schedule() = %v, want Disposable", gotD)
			}
			gotD.Dispose()

			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineMulti_ScheduleOnInterval(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	test1Count := 0
	tests := []struct {
		name      string
		fiber     Fiber
		args      Task
		firstMs   int64
		regularMs int64
		want1     int32
	}{
		{"Test_GoroutineMulti_ScheduleOnInterval",
			g,
			newTask(func(s string) {
				t.Logf("s:%v", s)
				test1Count++
			}), 50, 50, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 1")
			gotD := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 2")
			switch gotD.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineMulti.ScheduleOnInterval() = %v, want Disposable", gotD)
			}

			timeout := time.NewTimer(time.Duration(60) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD.Dispose()
			}
			timeout.Stop()
			timeout.Reset(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD.Dispose()
			}

			if tt.want1 != int32(test1Count) {
				t.Errorf("test 1 count %v, want %v", test1Count, tt.want1)
			}
		})
	}
}

func TestGoroutineMulti_RegisterSubscription(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		subscriptions  *Disposer
		flushPending   bool
	}
	type args struct {
		toAdd Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineMulti{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				subscriptions:  tt.fields.subscriptions,
				flushPending:   tt.fields.flushPending,
			}
			g.RegisterSubscription(tt.args.toAdd)
		})
	}
}

func TestGoroutineMulti_DeregisterSubscription(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		subscriptions  *Disposer
		flushPending   bool
	}
	type args struct {
		toRemove Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineMulti{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				subscriptions:  tt.fields.subscriptions,
				flushPending:   tt.fields.flushPending,
			}
			g.DeregisterSubscription(tt.args.toRemove)
		})
	}
}

func TestGoroutineMulti_NumSubscriptions(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		subscriptions  *Disposer
		flushPending   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GoroutineMulti{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				subscriptions:  tt.fields.subscriptions,
				flushPending:   tt.fields.flushPending,
			}
			if got := g.NumSubscriptions(); got != tt.want {
				t.Errorf("GoroutineMulti.NumSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoroutineMulti_flush(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		subscriptions  *Disposer
		flushPending   bool
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineMulti{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				subscriptions:  tt.fields.subscriptions,
				flushPending:   tt.fields.flushPending,
			}
			g.flush()
		})
	}
}

func TestGoroutineSingle_init(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	tests := []struct {
		name   string
		fields fields
		want   *GoroutineSingle
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			if got := g.init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoroutineSingle.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGoroutineSingle(t *testing.T) {
	tests := []struct {
		name string
		want *GoroutineSingle
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGoroutineSingle(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoroutineSingle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoroutineSingle_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineSingle_Start"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineSingle()
			g.Start()
			g.Start()
		})
	}
}

func TestGoroutineSingle_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineSingle_Stop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineSingle()
			g.Start()
			g.Stop()
		})
	}
}

func TestGoroutineSingle_Dispose(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_GoroutineSingle_Dispose"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineSingle()
			g.Dispose()
		})
	}
}

func TestGoroutineSingle_Enqueue(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name  string
		fiber Fiber
		args  interface{}
	}{
		{"Test_GoroutineSingle_Enqueue", g, func(s string) { t.Logf("s:%v", s) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.Enqueue(tt.args, "Test 1")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 2")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 3")
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineSingle_EnqueueWithTask(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name  string
		fiber Fiber
		args  Task
	}{
		{"Test_GoroutineSingle_EnqueueWithTask", g, newTask(func(s string) { t.Logf("s:%v", s) }, "Test 1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.EnqueueWithTask(tt.args)
			tt.fiber.EnqueueWithTask(newTask(func(s string) { t.Logf("s:%v", s) }, "Test 2"))
			tt.fiber.EnqueueWithTask(newTask(func(s string) { t.Logf("s:%v", s) }, "Test 3"))
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineSingle_Schedule(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name    string
		fiber   Fiber
		args    Task
		firstMs int64
	}{
		{"Test_GoroutineSingle_Schedule", g, newTask(func(s string) { t.Logf("s:%v", s) }), 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 1")
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 2")
			tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 3")
			gotD := tt.fiber.Schedule(tt.firstMs, tt.args.doFunc, "Test 4")
			switch gotD.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineSingle.Schedule() = %v, want Disposable", gotD)
			}
			gotD.Dispose()
			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineSingle_ScheduleOnInterval(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	test1Count := 0
	tests := []struct {
		name      string
		fiber     Fiber
		args      Task
		firstMs   int64
		regularMs int64
		want1     int32
	}{
		{"Test_GoroutineSingle_ScheduleOnInterval",
			g,
			newTask(func(s string) {
				t.Logf("s:%v", s)
				test1Count++
			}), 50, 50, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 1")
			gotD := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 2")
			switch gotD.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineSingle.ScheduleOnInterval() = %v, want Disposable", gotD)
			}

			timeout := time.NewTimer(time.Duration(60) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD.Dispose()
			}
			timeout.Stop()
			timeout.Reset(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD.Dispose()
			}

			if tt.want1 != int32(test1Count) {
				t.Errorf("test 1 count %v, want %v", test1Count, tt.want1)
			}
		})
	}
}

func TestGoroutineSingle_RegisterSubscription(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	type args struct {
		toAdd Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			g.RegisterSubscription(tt.args.toAdd)
		})
	}
}

func TestGoroutineSingle_DeregisterSubscription(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	type args struct {
		toRemove Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			g.DeregisterSubscription(tt.args.toRemove)
		})
	}
}

func TestGoroutineSingle_NumSubscriptions(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			if got := g.NumSubscriptions(); got != tt.want {
				t.Errorf("GoroutineSingle.NumSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoroutineSingle_executeNextBatch(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			if got := g.executeNextBatch(); got != tt.want {
				t.Errorf("GoroutineSingle.executeNextBatch() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestGoroutineSingle_dequeueAll(t *testing.T) {
	tests := []struct {
		name  string
		want  bool
		want1 bool
	}{
		{"Test_GoroutineSingle_dequeueAll", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoroutineSingle()
			got, got1 := g.dequeueAll()
			if got1 != tt.want {
				t.Errorf("GoroutineSingle.dequeueAll() got = %v, want %v", got, tt.want)
			}
			g.Start()
			g.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 1")
			g.Stop()
			g.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 2")
			g.Start()
			g.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 3")

			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineSingle_readyToDequeue(t *testing.T) {
	type fields struct {
		queue          taskQueue
		scheduler      IScheduler
		executor       executor
		executionState executionState
		lock           sync.Mutex
		cond           *sync.Cond
		subscriptions  *Disposer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GoroutineSingle{
				queue:          tt.fields.queue,
				scheduler:      tt.fields.scheduler,
				executor:       tt.fields.executor,
				executionState: tt.fields.executionState,
				lock:           tt.fields.lock,
				cond:           tt.fields.cond,
				subscriptions:  tt.fields.subscriptions,
			}
			if got := g.readyToDequeue(); got != tt.want {
				t.Errorf("GoroutineSingle.readyToDequeue() = %v, want %v", got, tt.want)
			}
		})
	}
}
