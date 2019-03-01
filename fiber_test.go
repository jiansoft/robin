package robin

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiber(t *testing.T) {
	tests := []struct {
		name          string
		count         int32
		intervalCount int32
		wg            sync.WaitGroup
	}{
		{"Test_Fiber", 0, 0, sync.WaitGroup{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := NewGoroutineMulti()
			gs := NewGoroutineSingle()

			gm.Start()
			gm.Start()
			gs.Start()
			gs.Start()

			loop := 100
			atomic.SwapInt32(&tt.count, 0)
			atomic.SwapInt32(&tt.intervalCount, 0)
			tt.wg.Add(loop * loop * 2 * 2)
			for ii := 0; ii < loop; ii++ {
				gm.Schedule(10, func() {
					for i := 0; i < loop; i++ {
						gm.Enqueue(func(s int) {
							atomic.AddInt32(&tt.count, 1)
							tt.wg.Done()
						}, i)
						gm.EnqueueWithTask(newTask(func() {
							atomic.AddInt32(&tt.count, 1)
							tt.wg.Done()
						}))
					}
				})
				gs.Schedule(10, func() {
					for i := 0; i < loop; i++ {
						gs.Enqueue(func(s int) {
							atomic.AddInt32(&tt.count, 1)
							tt.wg.Done()
						}, i)

						gs.EnqueueWithTask(newTask(func() {
							atomic.AddInt32(&tt.count, 1)
							tt.wg.Done()
						}))
					}
				})
			}
			tt.wg.Wait()

			assert.Equal(t, int32(loop*loop*2*2), atomic.LoadInt32(&tt.count), "they should be equal")

			tt.wg.Add(loop)
			gmd := gm.ScheduleOnInterval(10, 10, func() {
				atomic.AddInt32(&tt.intervalCount, 1)
				tt.wg.Done()
			})
			tt.wg.Wait()
			gmd.Dispose()
			assert.Equal(t, int32(loop), atomic.LoadInt32(&tt.intervalCount), "they should be equal")
			atomic.SwapInt32(&tt.intervalCount, 0)

			tt.wg.Add(loop)
			gsd := gs.ScheduleOnInterval(10, 10, func() {
				atomic.AddInt32(&tt.intervalCount, 1)
				tt.wg.Done()
			})

			tt.wg.Wait()
			gsd.Dispose()
			assert.Equal(t, int32(loop), atomic.LoadInt32(&tt.intervalCount), "they should be equal")

			gm.Stop()

			gm.EnqueueWithTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)

			}))

			gs.Stop()
			gs.EnqueueWithTask(newTask(func() {
				atomic.AddInt32(&tt.count, 1)

			}))

			_, got1 := gs.dequeueAll()
			assert.Equal(t, false, got1, "GoroutineSingle.dequeueAll() got = %v, want %v", got1, false)
			/*if got1 != false {
				t.Errorf("GoroutineSingle.dequeueAll() got = %v, want %v", got, false)
			}*/

			gm.Dispose()
			gs.Dispose()
		})
	}
}

//
//func TestGoroutineSingle_dequeueAll(t *testing.T) {
//	//lock := sync.Mutex{}
//	tests := []struct {
//		name  string
//		want  bool
//		want1 bool
//	}{
//		{"Test_GoroutineSingle_dequeueAll", false, true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			g := NewGoroutineSingle()
//			got, got1 := g.dequeueAll()
//			if got1 != tt.want {
//				t.Errorf("GoroutineSingle.dequeueAll() got = %v, want %v", got, tt.want)
//			}
//
//
//
//		})
//	}
//}
