package robin

import (
	"testing"
	"time"
)

// receiveEvent is a no-op subscriber used to isolate publish-path benchmark cost.
// receiveEvent 是空操作訂閱者，用來聚焦量測 publish 路徑本身的成本。
func (p player) receiveEvent(someBossInfo string) {
}

// BenchmarkChannel_Publish measures publish throughput of reflection-based Channel with 4 subscribers.
// BenchmarkChannel_Publish 量測反射版 Channel 在 4 個訂閱者下的發布吞吐量。
func BenchmarkChannel_Publish(b *testing.B) {
	p1 := player{Nickname: "Player 1"}
	p2 := player{Nickname: "Player 2"}
	p3 := player{Nickname: "Player 3"}
	p4 := player{Nickname: "Player 4"}
	channel := NewChannel()

	channel.Subscribe(p1.receiveEvent)
	channel.Subscribe(p2.receiveEvent)
	channel.Subscribe(p3.receiveEvent)
	channel.Subscribe(p4.receiveEvent)

	b.ResetTimer()
	for range b.N {
		channel.Publish(time.Now().Format("15:04:05.000"))
	}
}

// BenchmarkTypedChannel_Publish measures publish throughput of generic TypedChannel with 4 subscribers.
// BenchmarkTypedChannel_Publish 量測泛型 TypedChannel 在 4 個訂閱者下的發布吞吐量。
func BenchmarkTypedChannel_Publish(b *testing.B) {
	p1 := player{Nickname: "Player 1"}
	p2 := player{Nickname: "Player 2"}
	p3 := player{Nickname: "Player 3"}
	p4 := player{Nickname: "Player 4"}
	channel := NewTypedChannel[string]()

	channel.Subscribe(p1.receiveEvent)
	channel.Subscribe(p2.receiveEvent)
	channel.Subscribe(p3.receiveEvent)
	channel.Subscribe(p4.receiveEvent)

	b.ResetTimer()
	for range b.N {
		channel.Publish(time.Now().Format("15:04:05.000"))
	}
}
