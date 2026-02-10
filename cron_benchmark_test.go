package robin

import (
	"testing"
)

func BenchmarkCron_StartStop(b *testing.B) {
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
			base := make([]Disposable, c.N)
			for i := range base {
				base[i] = Every(100).Seconds().Do(func() {})
			}
			b.ResetTimer()

			for range b.N {
				Every(100).Seconds().Do(func() {}).Dispose()
			}

			b.StopTimer()
			for i := range base {
				base[i].Dispose()
			}
		})
	}
}
