package robin

import (
	"testing"
	"time"
)

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}

func BenchmarkTimingWheel_StartStop(b *testing.B) {

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000},
		{"N-5m", 5000},
		{"N-10m", 10000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]Disposable, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = Every(100).Seconds().Do(func() {})
				//base[i] = tw.AfterFunc(genD(i), func() {})
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				//tw.AfterFunc(time.Second, func() {}).Stop()
				Every(100).Seconds().Do(func() {}).Dispose()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Dispose()
			}
		})
	}
}
