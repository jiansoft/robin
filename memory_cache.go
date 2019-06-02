package robin

import (
	"math"
	"sync"
	"time"
)

const expirationInterval = 30

type MemoryCach interface {
	Remember(key string, value interface{}, ttl int64)
	Forget(key string)
	Read(string) (interface{}, bool)
}

type memoryCacheEntry struct {
	utcCreated int64
	utcAbsExp  int64
	key        string
	value      interface{}
}

type memoryCacheStore struct {
	usage   sync.Map
	expires *ConcurrentQueue
}

var (
	store *memoryCacheStore
	once  sync.Once
)

func Memory() MemoryCach {
	return memoryCache()
}

func memoryCache() *memoryCacheStore {
	once.Do(func() {
		if store == nil {
			store = new(memoryCacheStore)
			store.expires = NewConcurrentQueue()
			Every(expirationInterval).Seconds().AfterExecuteTask().Do(store.flushExpiredItems)
		}
	})

	return store
}

func (m *memoryCacheStore) Remember(key string, val interface{}, ttl int64) {
	e := &memoryCacheEntry{key: key, value: val}
	m.usage.Store(e.key, e)
	utc := time.Now().UTC().UnixNano()
	e.utcCreated = utc
	e.utcAbsExp = utc + ttl*int64(time.Second)
	//if e.utcAbsExp < e.utcCreated {
	//	atomic.SwapInt32(&e.inEexpire, 1)
	//	m.expires.Enqueue(e.key)
	//}
}

func (m *memoryCacheStore) Read(key string) (interface{}, bool) {
	t := time.Now().UTC().UnixNano()
	if val, ok := m.usage.Load(key); ok {
		//e := val.(*memoryCacheEntry)
		if val.(*memoryCacheEntry).utcAbsExp >= t {
			return val.(*memoryCacheEntry).value, true
		}
	}
	return nil, false
}

func (m *memoryCacheStore) Forget(key string) {
	m.usage.Delete(key)
}

// flushExpiredItems
func (m *memoryCacheStore) flushExpiredItems() {

	limit := time.Now().UTC().UnixNano()
	num := 0
	max := math.MaxInt32
	m.usage.Range(func(key, val interface{}) bool {
		e := val.(*memoryCacheEntry)
		if e.utcAbsExp < limit {
			//if atomic.CompareAndSwapInt32(&e.inEexpire, 0, 1) {
			num++
			m.expires.Enqueue(e.key)
			//}

			if num >= max {
				return false
			}
		}
		return true
	})

	for {
		val, ok := m.expires.TryDequeue()
		if !ok {
			break
		}
		m.usage.Delete(val)
	}

}
