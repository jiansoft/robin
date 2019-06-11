package robin

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkMemoryCache_Set(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-10k", 10000},
		{"N-100k", 100000},
		{"N-1000k", 1000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]string, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = strconv.FormatInt(int64(i), 10)
			}
			//cache := memoryCache()
			b.ResetTimer()
			//for i := 0; i < b.N; i++ {
			for i := 0; i < c.N; i++ {
				Memory().Keep(base[i], base[i], 10)
			}
			//}

			b.StopTimer()

		})
	}
}

func BenchmarkMemoryCache_flushExpiredItems(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-10k", 10000},
		{"N-100k", 100000},
		{"N-1000k", 1000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			cache := memoryCache()
			base := make([]string, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = strconv.FormatInt(int64(i), 10)
				Memory().Keep(base[i], base[i], 1)
			}

			b.ResetTimer()
			timeout := time.NewTimer(time.Duration(1010) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			cache.flushExpiredItems()

			b.StopTimer()

		})
	}
}
