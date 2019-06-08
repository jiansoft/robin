package robin

import (
	"fmt"
	"testing"
	"time"
)

func Test_memoryCacheStore_flushExpiredItems(t *testing.T) {
	m := memoryCache()
	tests := []struct {
		name        string
		memoryCache *memoryCacheStore
	}{
		{"1", m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1024; i++ {
				tt.memoryCache.Remember(fmt.Sprintf("QQ-%v", i), i, -time.Duration(1*time.Second))
			}
			tt.memoryCache.flushExpiredItems()

			equal(t, tt.memoryCache.pq.Len(), 0)
		})
	}
}

func Test_memoryCacheStore_Remember(t *testing.T) {
	//m := memoryCache()
	tests := []struct {
		name        string
		memoryCache *memoryCacheStore
		want        int
	}{
		{"1", memoryCache(), 1024},
		{"2", memoryCache(), 10240},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.want; i++ {
				tt.memoryCache.Remember(fmt.Sprintf("QQ-%s-%v", tt.name, i), i, time.Duration(-1*time.Hour))
			}
			equal(t, tt.memoryCache.pq.Len(), tt.want)
			tt.memoryCache.flushExpiredItems()
		})
	}
}

func Test_memoryCacheStore(t *testing.T) {
	tests := []struct {
		name        string
		memoryCache *memoryCacheStore
		want        int
	}{
		{"1", memoryCache(), 1024},
		{"2", memoryCache(), 10240},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var keyMap []string
			for i := 0; i < tt.want; i++ {
				key := fmt.Sprintf("QQ-%s-%v", tt.name, i)
				tt.memoryCache.Remember(key, key, 1*time.Hour)
				keyMap = append(keyMap, key)
			}

			for _, v := range keyMap {
				ok := tt.memoryCache.Have(v)
				equal(t, ok, true)
				val, ok := tt.memoryCache.Read(v)
				equal(t, ok, true)
				equal(t, val, v)

				tt.memoryCache.Forget(v)
				ok = tt.memoryCache.Have(v)
				equal(t, ok, false)
				val, ok = tt.memoryCache.Read(v)
				equal(t, ok, false)
			}
		})
	}
}
