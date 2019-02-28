package robin

import (
	"sync"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Test_Goroutine"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := NewGoroutineMulti()
			gm.Start()
			gm.Stop()
			gm.Dispose()
			gs := NewGoroutineSingle()
			gs.Start()
			gs.Stop()
			gs.Dispose()
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
			for ii := 0; ii < 10; ii++ {
				tt.fiber.Schedule(10, func() {
					for i := 0; i < 10000; i++ {
						tt.fiber.Enqueue(
							func(s int) {
								if s == 5000 {
									//                                    t.Logf("For:%v", s)
								}
							}, i)
					}
				})
			}

			tt.fiber.Enqueue(tt.args, "Test 1")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 2")
			tt.fiber.Enqueue(func(s string) { t.Logf("s:%v", s) }, "Test 3")
			timeout := time.NewTimer(time.Duration(2000) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}

func TestGoroutineMulti_EnqueueWithTask(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	lock := sync.Mutex{}
	lock.Lock()
	testCount := 0
	lock.Unlock()
	tests := []struct {
		name  string
		fiber *GoroutineMulti
		args  Task
		want  int32
	}{
		{"Test_GoroutineMulti_EnqueueWithTask",
			g,
			newTask(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				testCount++
				lock.Unlock()
			}, "Test 1"),
			2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.EnqueueWithTask(tt.args)
			//tt.fiber.EnqueueWithTask(newTask(tt.args.doFunc, "Test 2"))
			tt.fiber.EnqueueWithTask(newTask(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				testCount++
				lock.Unlock()
			}, "Test 2"))
			tt.fiber.executionState = stopped
			tt.fiber.EnqueueWithTask(newTask(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				testCount++
				lock.Unlock()
			}, "Test 3"))
			timeout := time.NewTimer(time.Duration(30) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			lock.Lock()
			defer lock.Unlock()
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
		{"Test_GoroutineMulti_Schedule", g, newTask(func(s string) { t.Logf("s:%v", s) }), 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for ii := 0; ii <= 10; ii++ {
				tt.fiber.Enqueue(func(ii int) {
					for i := 0; i <= 1000; i++ {
						tt.fiber.Schedule(
							tt.firstMs,
							func(i int, ii int) {
								if i == 1000 {
									//t.Logf("For i:%v ii:%v", i, ii)
								}
							}, i, ii)
					}
				}, ii)
			}
			test1 := tt.fiber.Schedule(tt.firstMs, func(s string) { t.Logf("s:%v", s) }, "Test 1")
			test2 := tt.fiber.Schedule(tt.firstMs, func(s string) { t.Logf("s:%v", s) }, "Test 2")
			test3 := tt.fiber.Schedule(tt.firstMs, func(s string) { t.Logf("s:%v", s) }, "Test 3")
			test4 := tt.fiber.Schedule(tt.firstMs, func(s string) { t.Logf("s:%v", s) }, "Test 4")
			switch test4.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineMulti.Schedule() = %v, want Disposable", test4)
			}

			timeout := time.NewTimer(time.Duration(1500) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			test1.Dispose()
			test2.Dispose()
			test3.Dispose()
			test4.Dispose()
		})
	}
}

func TestGoroutineMulti_ScheduleOnInterval(t *testing.T) {
	g := NewGoroutineMulti()
	g.Start()
	lock := sync.Mutex{}
	lock.Lock()
	test1Count := 0
	lock.Unlock()
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
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}), 50, 100, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD1 := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, func(s string) {
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}, "Test 1")
			gotD2 := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, func(s string) {
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}, "Test 2")
			switch gotD2.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineMulti.ScheduleOnInterval() = %v, want Disposable", gotD2)
			}

			timeout := time.NewTimer(time.Duration(70) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD2.Dispose()
			}
			timeout.Stop()
			timeout.Reset(time.Duration(125) * time.Millisecond)
			select {
			case <-timeout.C:
				gotD1.Dispose()
			}
			lock.Lock()
			if tt.want > int32(test1Count) {
				t.Errorf("%s test 1 count %v, want %v", tt.name, test1Count, tt.want)
			}
			lock.Unlock()
		})
	}
}

func TestGoroutineSingle_Enqueue(t *testing.T) {
	lock := sync.Mutex{}
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
			tt.fiber.Enqueue(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
			}, "Test 2")
			tt.fiber.Enqueue(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
			}, "Test 3")
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
	lock := sync.Mutex{}
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name    string
		fiber   Fiber
		firstMs int64
	}{
		{"Test_GoroutineSingle_Schedule", g, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(3)
			tt.fiber.Schedule(tt.firstMs, func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
				wg.Done()
			}, "Test 1")
			tt.fiber.Schedule(tt.firstMs, func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
				wg.Done()
			}, "Test 2")
			tt.fiber.Schedule(tt.firstMs, func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
				wg.Done()
			}, "Test 3")
			gotD := tt.fiber.Schedule(tt.firstMs, newTask(func(s string) { t.Logf("s:%v", s) }), "Test 4")
			switch gotD.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineSingle.Schedule() = %v, want Disposable", gotD)
			}
			gotD.Dispose()
			wg.Wait()
		})
	}
}

func TestGoroutineSingle_ScheduleOnInterval(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	lock := sync.Mutex{}
	lock.Lock()
	test1Count := 0
	lock.Unlock()
	tests := []struct {
		name      string
		fiber     Fiber
		args      Task
		firstMs   int64
		regularMs int64
		want      int32
	}{
		{"Test_GoroutineSingle_ScheduleOnInterval",
			g,
			newTask(func(s string) {
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}), 50, 100, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, func(s string) {
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}, "Test 1")
			gotD2 := tt.fiber.ScheduleOnInterval(tt.firstMs, tt.regularMs, func(s string) {
				//t.Logf("s:%v", s)
				lock.Lock()
				test1Count++
				lock.Unlock()
			}, "Test 2")
			switch gotD2.(type) {
			case Disposable:
			default:
				t.Errorf("GoroutineSingle.ScheduleOnInterval() = %v, want Disposable", gotD2)
			}

			timeout := time.NewTimer(time.Duration(65) * time.Millisecond)
			select {
			case <-timeout.C:
				//gotD2.Dispose()
			}
			timeout.Stop()
			timeout.Reset(time.Duration(130) * time.Millisecond)
			select {
			case <-timeout.C:
				//	gotD1.Dispose()
			}
			lock.Lock()
			defer lock.Unlock()
			if tt.want > int32(test1Count) {
				t.Errorf("%s count %v, want %v", tt.name, test1Count, tt.want)
			}
		})
	}
}

func TestGoroutineSingle_dequeueAll(t *testing.T) {
	lock := sync.Mutex{}
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
			g.Enqueue(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
			}, "Test 1")
			g.Stop()
			g.Enqueue(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
			}, "Test 2")
			g.Start()
			g.Enqueue(func(s string) {
				lock.Lock()
				t.Logf("s:%v", s)
				lock.Unlock()
			}, "Test 3")

			timeout := time.NewTimer(time.Duration(100) * time.Millisecond)
			select {
			case <-timeout.C:
			}
		})
	}
}
