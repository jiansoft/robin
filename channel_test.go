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

// recvMsgAction increments the shared receive counter for message-delivery assertions.
// recvMsgAction 會遞增共用接收計數器，用於驗證訊息是否被成功投遞。
func recvMsgAction(s string) {
	recvCount.Add(1)
}

// waitForCount polls until counter reaches target or timeout expires.
// waitForCount 會輪詢等待計數器達到目標值，逾時則使測試失敗。
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

// TestChannelPublishNoSubscribers verifies publishing on an empty channel is safe and count remains zero.
// TestChannelPublishNoSubscribers 驗證在無訂閱者時發布訊息是安全的，且訂閱數維持為 0。
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

// TestChannel validates subscribe/publish/unsubscribe/clear lifecycle and subscriber count consistency.
// TestChannel 驗證訂閱/發布/取消訂閱/清空的完整生命週期，以及訂閱者數量一致性。
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

// TestChannelConcurrency exercises concurrent subscribe/publish paths to ensure no deadlock or missed completion.
// TestChannelConcurrency 驗證並發訂閱/發布路徑，確保不死鎖且可完整完成。
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

// player is a lightweight test receiver model used in concurrency scenarios.
// player 是並發情境測試中使用的輕量接收者模型。
type player struct {
	Nickname string
}

// recvMsgAction records one delivered message for this player.
// recvMsgAction 會記錄此玩家收到的一次訊息投遞。
func (p player) recvMsgAction(msg string) {
	recvCount.Add(1)
}
