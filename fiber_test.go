package robin

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestFiber validates high-volume enqueue/flush behavior on both multi and single fiber implementations.
// TestFiber 驗證多/單 fiber 在高量 enqueue 與 flush 情境下的行為正確性。
func TestFiber(t *testing.T) {
	tests := []struct {
		name          string
		count         int32
		intervalCount int32
	}{
		{"TestFiber", 0, 0},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := NewGoroutineMulti()
			gs := NewGoroutineSingle()
			loop := 100
			atomic.SwapInt32(&tt.count, 0)
			atomic.SwapInt32(&tt.intervalCount, 0)
			totalCount := loop * loop * 2 * 2
			wg.Add(totalCount)
			for ii := range loop {
				_ = ii
				gm.Schedule(10, func() {
					for i := range loop {
						gm.Enqueue(func(s int) {
							atomic.AddInt32(&tt.count, 1)
							wg.Done()
						}, i)
						gm.enqueueTask(newTask(func() {
							atomic.AddInt32(&tt.count, 1)
							wg.Done()
						}))
					}
				})
				gs.Schedule(10, func() {
					for i := range loop {
						gs.Enqueue(func(s int) {
							atomic.AddInt32(&tt.count, 1)
							wg.Done()
						}, i)
						gs.enqueueTask(newTask(func() {
							atomic.AddInt32(&tt.count, 1)
							wg.Done()
						}))
					}
				})
			}
			wg.Wait()

			if int32(totalCount) != atomic.LoadInt32(&tt.count) {
				t.Fatalf("they should be equal %v %v", totalCount, tt.count)
			}

			wg.Add(loop)
			gmd := gm.ScheduleOnInterval(10, 10, func() {
				atomic.AddInt32(&tt.intervalCount, 1)
				wg.Done()
			})
			wg.Wait()
			gmd.Dispose()
			if int32(loop) != atomic.LoadInt32(&tt.intervalCount) {
				t.Fatal("they should be equal")
			}

			atomic.SwapInt32(&tt.intervalCount, 0)

			wg.Add(loop)
			gsd := gs.ScheduleOnInterval(10, 10, func() {
				atomic.AddInt32(&tt.intervalCount, 1)
				wg.Done()
			})

			wg.Wait()
			gsd.Dispose()
			if int32(loop) != atomic.LoadInt32(&tt.intervalCount) {
				t.Fatal("they should be equal")
			}

			gm.enqueueTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)
			}))

			gm.flush()

			gs.enqueueTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)
			}))

			gm.Dispose()
			gs.Dispose()
		})
	}
}

// TestFiberSchedule validates schedule and interval timing paths for both fiber variants.
// TestFiberSchedule 驗證兩種 fiber 的一次性與週期性排程路徑。
func TestFiberSchedule(t *testing.T) {
	type fields struct {
		fiber Fiber
	}

	tests := []struct {
		fields fields
		name   string
		loop   int32
		want   int32
	}{
		{fields{NewGoroutineSingle()}, "GoroutineSingle", 0, 7},
		{fields{NewGoroutineMulti()}, "GoroutineMulti", 0, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.SwapInt32(&tt.loop, 0)
			tt.fields.fiber.ScheduleOnInterval(400, 400, func() {
				atomic.AddInt32(&tt.loop, 1)
			})

			tt.fields.fiber.ScheduleOnInterval(800, 800, func() {
				t.Logf("%s ScheduleOnInterval 800,800", tt.name)
			})

			tt.fields.fiber.Schedule(2500, func() {
				t.Logf("%s Schedule 2500", tt.name)
			})

			<-time.After(3 * time.Second)
			tt.fields.fiber.Dispose()
			loop := atomic.LoadInt32(&tt.loop)
			if tt.want != loop {
				t.Fatalf("%s ScheduleOnInterval error want %v got:%v", tt.name, tt.want, loop)
			}
		})
	}
}

// TestFiberEnqueueAfterDispose validates enqueue calls become no-op after disposal.
// TestFiberEnqueueAfterDispose 驗證 fiber Dispose 後 enqueue 會成為 no-op。
func TestFiberEnqueueAfterDispose(t *testing.T) {
	t.Run("GoroutineMulti", func(t *testing.T) {
		gm := NewGoroutineMulti()
		gm.Dispose()
		// enqueueTask after dispose should not panic, just return
		gm.Enqueue(func() {
			t.Error("should not execute after dispose")
		})
		gm.enqueueTask(newTask(func() {
			t.Error("should not execute after dispose")
		}))
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("GoroutineSingle", func(t *testing.T) {
		gs := NewGoroutineSingle()
		gs.Dispose()
		// enqueueTask after dispose should not panic, just return
		gs.Enqueue(func() {
			t.Error("should not execute after dispose")
		})
		gs.enqueueTask(newTask(func() {
			t.Error("should not execute after dispose")
		}))
		time.Sleep(50 * time.Millisecond)
	})
}

// TestFiber_GoroutineSingle_executeNextBatch validates dedicated worker loop can exit after Dispose.
// TestFiber_GoroutineSingle_executeNextBatch 驗證 GoroutineSingle 在 Dispose 後可正常退出 worker 迴圈。
func TestFiber_GoroutineSingle_executeNextBatch(t *testing.T) {
	tests := []struct {
		name string
		exit int32
	}{
		{"GoroutineSingle", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := new(GoroutineSingle)
			g.queue = newDefaultQueue()
			g.scheduler = newScheduler(g)
			g.executor = newDefaultExecutor()
			g.cond = sync.NewCond(&g.mu)

			go func() {
				for g.executeNextBatch() {
				}
				atomic.SwapInt32(&tt.exit, 1)
			}()

			g.Dispose()
			<-time.After(200 * time.Millisecond)

			if atomic.LoadInt32(&tt.exit) != 1 {
				t.Fatalf("%s executeNextBatch is not exit Goroutine", tt.name)
			}
		})
	}
}

// TestFiber_executeNextBatch_clearsUpcomingTaskReferences validates processed task references are cleared.
// TestFiber_executeNextBatch_clearsUpcomingTaskReferences 驗證已處理任務的引用會被清空。
func TestFiber_executeNextBatch_clearsUpcomingTaskReferences(t *testing.T) {
	g := new(GoroutineSingle)
	g.queue = newDefaultQueue()
	g.scheduler = newScheduler(g)
	g.executor = newDefaultExecutor()
	g.cond = sync.NewCond(&g.mu)

	dq, ok := g.queue.(*defaultQueue)
	if !ok {
		t.Fatal("queue type assertion failed")
	}

	g.enqueueTask(newTask(func() {}))
	if !g.executeNextBatch() {
		t.Fatal("executeNextBatch returned false unexpectedly")
	}

	if len(dq.upcoming) == 0 {
		t.Fatal("expected upcoming buffer to contain the processed slot")
	}

	cleared := dq.upcoming[0]
	if cleared.fn != nil || cleared.invoke != nil || cleared.invokeArg != nil || cleared.funcCache.IsValid() || len(cleared.paramsCache) != 0 {
		t.Fatal("processed task references were not cleared from upcoming buffer")
	}
}
