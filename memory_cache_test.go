package robin

import (
	"fmt"
	"sync"
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
				_ = tt.memoryCache.Keep(fmt.Sprintf("QQ-%v", i), i, -1*time.Second)
			}
			tt.memoryCache.flushExpiredItems()
			equal(t, tt.memoryCache.pq.Len(), 0)
		})
	}
}

func Test_memoryCacheStore_Keep(t *testing.T) {
	tests := []struct {
		name        string
		memoryCache *memoryCacheStore
	}{
		{"1", memoryCache()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.memoryCache.Keep(fmt.Sprintf("QQ-%s-%v", tt.name, 1), 1, -1*time.Hour)
			equal(t, tt.memoryCache.pq.Len(), 0)
			_ = tt.memoryCache.Keep("QQ-1-1", 1, 1*time.Hour)
			equal(t, tt.memoryCache.pq.Len(), 1)
			v1, _ := tt.memoryCache.Read("QQ-1-1")
			equal(t, v1, 1)
			_ = tt.memoryCache.Keep("QQ-1-1", 11, 1*time.Minute)
			v2, _ := tt.memoryCache.Read("QQ-1-1")
			equal(t, v2, 11)
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
		/* {"2", memoryCache(), 10240},
		   {"3", memoryCache(), 102400},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.want; i++ {
				key := fmt.Sprintf("QQ-%s-%v", tt.name, i)
				err := tt.memoryCache.Keep(key, key, 1*time.Hour)
				if err != nil {
					t.Logf("err:%v", err)
					continue
				}

				yes := tt.memoryCache.Have(key)
				equal(t, yes, true)
				val, ok := tt.memoryCache.Read(key)
				equal(t, ok, true)
				equal(t, val, key)

				tt.memoryCache.Forget(key)
				yes = tt.memoryCache.Have(key)
				equal(t, yes, false)
				val, ok = tt.memoryCache.Read(key)
				equal(t, ok, false)
			}
			tt.memoryCache.Forget("noKey")
			_, ok := tt.memoryCache.Read("noKey")
			equal(t, ok, false)
			_, ok = tt.memoryCache.loadMemoryCacheEntry("noKey")
			equal(t, ok, false)
		})
	}
}

func Test_DataRace(t *testing.T) {
	tests := []struct {
		name        string
		memoryCache *memoryCacheStore
		want        int
	}{
		{"1", memoryCache(), 1024},
		/* {"2", memoryCache(), 10240},
		   {"3", memoryCache(), 102400},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(4)
			RightNow().Do(func(want int, m *memoryCacheStore, swg *sync.WaitGroup) {
				for i := 0; i < want; i++ {
					key := fmt.Sprintf("RightNow-1-%v", i)
					_ = Memory().Keep(key, key, 1*time.Hour)
				}
				m.flushExpiredItems()
				swg.Done()
			}, tt.want, tt.memoryCache, &wg)

			RightNow().Do(func(want int, m *memoryCacheStore, swg *sync.WaitGroup) {
				for i := 0; i < want; i++ {
					key := fmt.Sprintf("RightNow-1-%v", i)
					m.Forget(key)
				}
				m.flushExpiredItems()
				swg.Done()
			}, tt.want, tt.memoryCache, &wg)

			RightNow().Do(func(want int, m *memoryCacheStore, swg *sync.WaitGroup) {
				for i := 0; i < want; i++ {
					key := fmt.Sprintf("RightNow-1-%v", i)
					_, _ = m.Read(key)
				}
				m.flushExpiredItems()
				swg.Done()
			}, tt.want, Memory(), &wg)

			RightNow().Do(func(want int, m *memoryCacheStore, swg *sync.WaitGroup) {
				for i := 0; i < want; i++ {
					key := fmt.Sprintf("QQ-%v", i)
					err := m.Keep(key, key, 1*time.Hour)
					if err != nil {
						t.Logf("err:%v", err)
						continue
					}

					_ = m.Have(key)
					_, _ = m.Read(key)
					m.Forget(key)
					_ = m.Have(key)
					_, _ = m.Read(key)
				}
				swg.Done()
			}, tt.want, tt.memoryCache, &wg)

			Every(10).Milliseconds().Times(10000).Do(tt.memoryCache.flushExpiredItems)
			wg.Add(1)
			RightNow().Do(keep, tt.want, tt.memoryCache, &wg, 1)
			wg.Add(1)
			RightNow().Do(read, tt.want, tt.memoryCache, &wg, 1)
			wg.Add(1)
			RightNow().Do(keep, tt.want, tt.memoryCache, &wg, 2)
			wg.Add(1)
			RightNow().Do(read, tt.want, tt.memoryCache, &wg, 2)
			wg.Add(1)
			RightNow().Do(keep, tt.want, tt.memoryCache, &wg, 3)
			wg.Wait()
		})
	}
}

func keep(want int, m *memoryCacheStore, swg *sync.WaitGroup, index int) {
	for i := 0; i < want; i++ {
		key := fmt.Sprintf("QQ-%v-%v", i, index)
		_ = m.Keep(key, key, time.Duration(int64(10+i)*int64(time.Millisecond)))
	}
	swg.Done()
}

func read(want int, m *memoryCacheStore, swg *sync.WaitGroup, index int) {
	for i := 0; i < want; i++ {
		key := fmt.Sprintf("QQ-%v-%v", i, index)

		m.Forget(key)
		_ = m.Have(key)
		_, _ = m.Read(key)

		m.Forget(key)
		_ = m.Have(key)
		_, _ = m.Read(key)
	}
	swg.Done()
}
