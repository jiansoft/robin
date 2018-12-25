package robin

import (
	"testing"
	"time"
)

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
	testCount := 0
	tests := []struct {
		name  string
		fiber *GoroutineMulti
		args  Task
		want  int32
	}{
		{"Test_GoroutineMulti_EnqueueWithTask",
			g,
			newTask(func(s string) {
				t.Logf("s:%v", s)
				testCount++
			}, "Test 1"),
			2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.EnqueueWithTask(tt.args)
			tt.fiber.EnqueueWithTask(newTask(tt.args.doFunc, "Test 2"))
			tt.fiber.executionState = stopped
			tt.fiber.EnqueueWithTask(newTask(func(s string) { t.Logf("s:%v", s) }, "Test 3"))
			timeout := time.NewTimer(time.Duration(30) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			if tt.want != int32(testCount) {
				t.Errorf("%s test count %v, want %v", tt.name, testCount, tt.want)
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
		want      int32
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
			gotD1 := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 1")
			gotD2 := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, tt.args.doFunc, "Test 2")
			switch gotD2.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineMulti.ScheduleOnInterval() = %v, want Disposable", gotD2)
			}

			timeout := time.NewTimer(time.Duration(60) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD2.Dispose()
			}
			timeout.Stop()
			timeout.Reset(time.Duration(120) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD1.Dispose()
			}

			if tt.want != int32(test1Count) {
				t.Errorf("%s test 1 count %v, want %v", tt.name, test1Count, tt.want)
			}
		})
	}
}

func TestGoroutineMulti_Subscription(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name  string
		fiber *GoroutineMulti
		want  int
		want1 int
	}{
		{"Test_GoroutineMulti_Subscription", g, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fiber.Schedule(1000, func() {})
			tt.fiber.RegisterSubscription(d)
			if tt.want != tt.fiber.NumSubscriptions() {
				t.Errorf("%s test count %v, want %v", tt.name, tt.fiber.NumSubscriptions(), tt.want)
			}
			tt.fiber.DeregisterSubscription(d)
			if tt.want1 != tt.fiber.NumSubscriptions() {
				t.Errorf("%s test count %v, want %v", tt.name, tt.fiber.NumSubscriptions(), tt.want1)
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
			//if g.executionState != running {
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
			timeout.Reset(time.Duration(120) * time.Millisecond)
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

func TestGoroutineSingle_Subscription(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name  string
		fiber *GoroutineSingle
		want  int
		want1 int
	}{
		{"Test_GoroutineSingle_Subscription", g, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fiber.Schedule(1000, func() {})
			tt.fiber.RegisterSubscription(d)
			if tt.want != tt.fiber.NumSubscriptions() {
				t.Errorf("%s test count %v, want %v", tt.name, tt.fiber.NumSubscriptions(), tt.want)
			}
			tt.fiber.DeregisterSubscription(d)
			if tt.want1 != tt.fiber.NumSubscriptions() {
				t.Errorf("%s test count %v, want %v", tt.name, tt.fiber.NumSubscriptions(), tt.want1)
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
