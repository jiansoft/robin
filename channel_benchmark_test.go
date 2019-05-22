package robin

import (
	"testing"
	"time"
)

func (p player) receiveEvent(someBossInfo string) {
}

func BenchmarkChannel_Publish(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1k", 1000},
		{"N-5k", 5000},
		{"N-10k", 10000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			p1 := player{Nickname: "Player 1"}
			p2 := player{Nickname: "Player 2"}
			p3 := player{Nickname: "Player 3"}
			p4 := player{Nickname: "Player 4"}
			channel := NewChannel()

			channel.Subscribe(p1.receiveEvent)
			channel.Subscribe(p2.receiveEvent)
			channel.Subscribe(p3.receiveEvent)
			channel.Subscribe(p4.receiveEvent)

			for i := 0; i < c.N; i++ {
				channel.Publish(time.Now().Format("15:04:05.000"))
			}
			//b.ResetTimer()

			/*for i := 0; i < b.N; i++ {
			      Every(100).Seconds().Do(func() {}).Dispose()
			  }

			  b.StopTimer()
			  for i := 0; i < len(base); i++ {
			      base[i].Dispose()
			  }*/
		})
	}
}
