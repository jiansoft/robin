package robin

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// gameEvent is a struct used to test TypedChannel with non-primitive types.
type gameEvent struct {
	Source string
	Target string
	Damage int
	Time   time.Time
}

func TestTypedChannelPublishNoSubscribers(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		ch := NewTypedChannel[string]()
		ch.Publish("hello")
		ch.Publish("world")
		<-time.After(100 * time.Millisecond)
		if ch.Count() != 0 {
			t.Errorf("TypedChannel count = %d, want 0", ch.Count())
		}
	})

	t.Run("struct", func(t *testing.T) {
		ch := NewTypedChannel[gameEvent]()
		ch.Publish(gameEvent{Source: "boss", Damage: 999})
		ch.Publish(gameEvent{})
		<-time.After(100 * time.Millisecond)
		if ch.Count() != 0 {
			t.Errorf("TypedChannel count = %d, want 0", ch.Count())
		}
	})
}

func TestTypedChannel(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		testTypedChannelFlow(t,
			NewTypedChannel[string](),
			func(i int) string {
				return fmt.Sprintf("Publish message %d", i)
			},
		)
	})

	t.Run("struct", func(t *testing.T) {
		testTypedChannelFlow(t,
			NewTypedChannel[gameEvent](),
			func(i int) gameEvent {
				return gameEvent{
					Source: "player",
					Target: "boss",
					Damage: i * 100,
					Time:   time.Now(),
				}
			},
		)
	})
}

// testTypedChannelFlow exercises the full subscribe/publish/unsubscribe/clear lifecycle
// for any message type T.
func testTypedChannelFlow[T any](t *testing.T, channel *TypedChannel[T], makeMsg func(int) T) {
	t.Helper()
	var recvCount atomic.Int32

	handler := func(_ T) {
		recvCount.Add(1)
	}

	channel.Subscribe(handler)
	recvCount.Store(0)

	subscriber1 := channel.Subscribe(handler)
	subscriber2 := channel.Subscribe(handler)

	// Publish 1: 3 subscribers → cumulative count = 3
	channel.Publish(makeMsg(1))
	waitForCount(t, &recvCount, 3)

	subscriber1.Unsubscribe()
	channel.Unsubscribe(subscriber2)

	// Publish 2: 1 subscriber → cumulative count = 4
	channel.Publish(makeMsg(2))
	waitForCount(t, &recvCount, 4)

	channel.Subscribe(handler)

	// Publish 3: 2 subscribers → cumulative count = 6
	channel.Publish(makeMsg(3))
	waitForCount(t, &recvCount, 6)

	channel.Clear()
	// Publish 4: 0 subscribers → count stays at 6
	channel.Publish(makeMsg(4))
	<-time.After(50 * time.Millisecond)

	if got := recvCount.Load(); got != 6 {
		t.Errorf("recvCount = %v, want 6", got)
	}

	subscriber1 = channel.Subscribe(handler)
	subscriber2 = channel.Subscribe(handler)
	channel.Subscribe(handler)
	if channel.Count() != 3 {
		t.Errorf("Channel count = %v, want 3", channel.Count())
	}

	subscriber1.Unsubscribe()
	if channel.Count() != 2 {
		t.Errorf("Channel count = %v, want 2", channel.Count())
	}

	subscriber2.Unsubscribe()
	if channel.Count() != 1 {
		t.Errorf("Channel count = %v, want 1", channel.Count())
	}

	channel.Unsubscribe(subscriber1)
	channel.Unsubscribe(subscriber2)
	if channel.Count() != 1 {
		t.Errorf("Channel count = %v, want 1", channel.Count())
	}

	channel.Clear()
	if channel.Count() != 0 {
		t.Errorf("Channel count = %v, want 0", channel.Count())
	}
}

func TestTypedChannelStructFieldVerify(t *testing.T) {
	ch := NewTypedChannel[gameEvent]()
	var received atomic.Value

	ch.Subscribe(func(e gameEvent) {
		received.Store(e)
	})

	sent := gameEvent{Source: "archer", Target: "dragon", Damage: 42, Time: time.Now()}
	ch.Publish(sent)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if v := received.Load(); v != nil {
			got := v.(gameEvent)
			if got.Source != sent.Source || got.Target != sent.Target || got.Damage != sent.Damage {
				t.Errorf("received = %+v, want %+v", got, sent)
			}
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("timed out waiting for struct event")
}

func TestTypedChannelConcurrency(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		channel := NewTypedChannel[string]()
		var mu sync.RWMutex
		var wg sync.WaitGroup
		count := 0
		loop := 8

		wg.Add(loop * loop)
		for range loop {
			RightNow().Do(func(c *TypedChannel[string]) {
				for range loop {
					mu.RLock()
					nickname := strconv.Itoa(count)
					mu.RUnlock()
					c.Subscribe(func(msg string) {
						_ = nickname
					})
					mu.Lock()
					count++
					mu.Unlock()
				}

				for range loop {
					c.Publish(fmt.Sprintf("Publish message Channel:%v", c.Count()))
					wg.Done()
				}
			}, channel)
		}

		wg.Wait()
	})

	t.Run("struct", func(t *testing.T) {
		channel := NewTypedChannel[gameEvent]()
		var mu sync.RWMutex
		var wg sync.WaitGroup
		count := 0
		loop := 8

		wg.Add(loop * loop)
		for range loop {
			RightNow().Do(func(c *TypedChannel[gameEvent]) {
				for range loop {
					mu.RLock()
					id := count
					mu.RUnlock()
					c.Subscribe(func(e gameEvent) {
						_ = id
					})
					mu.Lock()
					count++
					mu.Unlock()
				}

				for range loop {
					c.Publish(gameEvent{Source: "npc", Damage: c.Count()})
					wg.Done()
				}
			}, channel)
		}

		wg.Wait()
	})
}
