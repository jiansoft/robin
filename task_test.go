package robin

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_timerTask_schedule(t *testing.T) {
	var runCount int32
	g := NewGoroutineMulti()
	g.Start()
	tests := []struct {
		name   string
		fields *timerTask
	}{
		{"Test_timerTask_schedule_1", newTimerTask(g.scheduler.(*scheduler), newTask(func() {
			atomic.AddInt32(&runCount, 1)
		}), 0, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			if tt.name == "Test_timerTask_schedule_4" {
				tt.fields.disposed = 1
			}
			tt.fields.schedule()
			for true {
				saveRunCount := atomic.LoadInt32(&runCount)
				if tt.name == "Test_timerTask_schedule_1" && saveRunCount >= 10 {
					tt.fields.Dispose()
					break
				}
			}
			tt.fields.Dispose()

			wg.Add(2)
			var runT1Count int32
			t1 := g.ScheduleOnInterval(0, 10, func() {
				atomic.AddInt32(&runT1Count, 1)
				wg.Done()
			})
			wg.Wait()
			t1.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			assert.Equal(t, int32(2), atomic.LoadInt32(&runT1Count), "they should be equal")

			var runT2Count int32
			t2 := g.ScheduleOnInterval(1000, 10, func() {
				atomic.AddInt32(&runT2Count, 1)
			})
			t2.Dispose()

			<-time.After(time.Duration(30) * time.Millisecond)
			assert.Equal(t, int32(0), atomic.LoadInt32(&runT2Count), "they should be equal")

			var runT3Count int32
			t3 := g.Schedule(1000, func() {
				atomic.AddInt32(&runT3Count, 1)
			})
			t3.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			assert.Equal(t, int32(0), atomic.LoadInt32(&runT3Count), "they should be equal")

			var runT4Count int32
			t4 := g.Schedule(0, func() {
				atomic.AddInt32(&runT4Count, 1)
			})
			t4.Dispose()
			<-time.After(time.Duration(30) * time.Millisecond)
			assert.Equal(t, int32(1), atomic.LoadInt32(&runT4Count), "they should be equal")
		})
	}
}
