package robin

import (
	"testing"
	"time"
)

func (p player) receiveEvent(someBossInfo string) {
}

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
	for i := 0; i < b.N; i++ {
		channel.Publish(time.Now().Format("15:04:05.000"))
	}
}
