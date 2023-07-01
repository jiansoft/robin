package robin

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

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
			for ii := 0; ii < loop; ii++ {
				gm.Schedule(10, func() {
					for i := 0; i < loop; i++ {
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
					for i := 0; i < loop; i++ {
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
			//assert.Equal(t, int32(loop), atomic.LoadInt32(&tt.intervalCount), "they should be equal")
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
			//assert.Equal(t, int32(loop), atomic.LoadInt32(&tt.intervalCount), "they should be equal")
			if int32(loop) != atomic.LoadInt32(&tt.intervalCount) {
				t.Fatal("they should be equal")
			}
			/*gm.Stop()
			gs.Stop()*/

			gm.enqueueTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)

			}))

			gm.flush()

			gs.enqueueTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)

			}))

			gm.Dispose()
			gs.Dispose()
			//_, got := gs.dequeueAll()
			//assert.Equal(t, false, got, "GoroutineSingle.dequeueAll() got = %v, want %v", got, false)
			//if false != got {
			//	t.Fatalf("GoroutineSingle.dequeueAll() got = %v, want %v", got, false)
			//}

		})
	}
}

func TestFiberSchedule(t *testing.T) {
	type fields struct {
		fiber IFiber
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

			<-time.After(time.Duration(3) * time.Second)
			tt.fields.fiber.Dispose()
			loop := atomic.LoadInt32(&tt.loop)
			if tt.want != loop {
				t.Fatalf("%s ScheduleOnInterval error want %v got:%v", tt.name, tt.want, loop)
			}
		})
	}
}

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
			<-time.After(time.Duration(200) * time.Millisecond)

			if atomic.LoadInt32(&tt.exit) != 1 {
				t.Fatalf("%s executeNextBatch is not exit Goroutine", tt.name)
			}
		})
	}
}
