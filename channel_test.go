package robin

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var recvCount atomic.Int32

func recvMsgAction(s string) {
	recvCount.Add(1)
}

// waitForCount polls until counter reaches target or timeout expires.
func waitForCount(t *testing.T, counter *atomic.Int32, target int32) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if counter.Load() >= target {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for count to reach %d, got %d", target, counter.Load())
}

func TestChannelPublishNoSubscribers(t *testing.T) {
	ch := NewChannel()
	// Publish with no subscribers should not panic
	ch.Publish("hello")
	ch.Publish("world")
	<-time.After(100 * time.Millisecond)
	if ch.Count() != 0 {
		t.Errorf("Channel count = %d, want 0", ch.Count())
	}
}

func TestChannel(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{name: "TestChannel", want: 6},
	}

	channel := NewChannel()
	channel.Subscribe(recvMsgAction)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recvCount.Store(0)

			subscriber1 := channel.Subscribe(recvMsgAction)
			subscriber2 := channel.Subscribe(recvMsgAction)

			// Publish 1: 3 subscribers → cumulative count = 3
			channel.Publish(fmt.Sprintf("Publish message 1 Channel:%v", channel.Count()))
			waitForCount(t, &recvCount, 3)

			subscriber1.Unsubscribe()
			channel.Unsubscribe(subscriber2)
			subscriber1 = nil
			subscriber2 = nil

			// Publish 2: 1 subscriber → cumulative count = 4
			channel.Publish(fmt.Sprintf("Publish message 2 Channel:%v", channel.Count()))
			waitForCount(t, &recvCount, 4)

			channel.Subscribe(recvMsgAction)

			// Publish 3: 2 subscribers → cumulative count = 6
			channel.Publish(fmt.Sprintf("Publish message 3 Channel:%v", channel.Count()))
			waitForCount(t, &recvCount, 6)

			channel.Clear()
			// Publish 4: 0 subscribers → count stays at 6
			channel.Publish(fmt.Sprintf("Publish message 4 Channel:%v", channel.Count()))
			<-time.After(50 * time.Millisecond)

			if got := recvCount.Load(); got != tt.want {
				t.Errorf("recvCount = %v, want %v", got, tt.want)
			}

			subscriber1 = channel.Subscribe(recvMsgAction)
			subscriber2 = channel.Subscribe(recvMsgAction)
			channel.Subscribe(recvMsgAction)
			if channel.Count() != 3 {
				t.Errorf("the Channel want %v but got %v", 3, channel.Count())
			}

			subscriber1.Unsubscribe()
			if channel.Count() != 2 {
				t.Errorf("the Channel want %v but got %v", 2, channel.Count())
			}

			subscriber2.Unsubscribe()
			if channel.Count() != 1 {
				t.Errorf("the Channel want %v but got %v", 1, channel.Count())
			}

			channel.Unsubscribe(subscriber1)
			channel.Unsubscribe(subscriber2)
			if channel.Count() != 1 {
				t.Errorf("the Channel want %v but got %v", 1, channel.Count())
			}

			channel.Clear()
			if channel.Count() != 0 {
				t.Errorf("the Channel should be empty %v, want %v", channel.Count(), 0)
			}
		})
	}
}

func TestChannelConcurrency(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test_TestConcurrency_1"},
	}
	lock := sync.RWMutex{}
	channel := NewChannel()
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock.Lock()
			count := 0
			lock.Unlock()
			loop := 8
			wg.Add(loop * loop)
			for range loop {
				RightNow().Do(func(c *Channel) {
					for range loop {
						lock.RLock()
						p := &player{Nickname: strconv.Itoa(count)}
						lock.RUnlock()
						c.Subscribe(p.recvMsgAction)
						lock.Lock()
						count++
						lock.Unlock()
					}

					for range loop {
						c.Publish(fmt.Sprintf("Publish message Channel:%v", c.Count()))
						wg.Done()
					}
				}, channel)
			}

			wg.Wait()
		})
	}
}

type player struct {
	Nickname string
}

func (p player) recvMsgAction(msg string) {
	recvCount.Add(1)
}
