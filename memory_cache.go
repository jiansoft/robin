package robin

import (
	"errors"
	"math"
	"sync"
	"time"
)

const expirationInterval = 30

type MemoryCache interface {
	Keep(key string, value interface{}, ttl time.Duration) error
	Forget(key string)
	Read(string) (interface{}, bool)
	Have(string) bool
}

type memoryCacheEntity struct {
	utcCreated int64
	utcAbsExp  int64
	key        string
	value      interface{}
	item       *Item
	sync.Mutex
}

type memoryCacheStore struct {
	usage sync.Map
	pq    *PriorityQueue
	sync.Mutex
}

var (
	store *memoryCacheStore
	once  sync.Once
)

var (
	errNewCacheHasExpired = errors.New("memory cache: the TTL(Time To Live) is less than the current time")
)

// Memory returns memoryCacheStore instance.
func Memory() MemoryCache {
	return memoryCache()
}

// memoryCache returns memoryCacheStore dingleton instance
func memoryCache() *memoryCacheStore {
	once.Do(func() {
		store = new(memoryCacheStore)
		store.pq = NewPriorityQueue(1024)
		Every(expirationInterval).Seconds().AfterExecuteTask().Do(store.flushExpiredItems)
	})
	return store
}

// isExpired returns true if the item has expired with the locker..
func (m *memoryCacheEntity) isExpired() bool {
	m.Lock()
	defer m.Unlock()
	return time.Now().UTC().UnixNano() > m.utcAbsExp
}

// expired to set the item has expired with the locker.
func (m *memoryCacheEntity) expired() {
	m.Lock()
	m.utcAbsExp = 0
	m.Unlock()
	m.item.expired()
}

// loadMemoryCacheEntry returns memoryCacheEntity if it exists in the cache
func (m *memoryCacheStore) loadMemoryCacheEntry(key string) (*memoryCacheEntity, bool) {
	val, yes := m.usage.Load(key)
	if !yes {
		return nil, false
	}
	return val.(*memoryCacheEntity), true
}

// Keep insert an item into the memory
func (m *memoryCacheStore) Keep(key string, val interface{}, ttl time.Duration) error {
	nowUtc := time.Now().UTC().UnixNano()
	utcAbsExp := nowUtc + ttl.Nanoseconds()
	if utcAbsExp <= nowUtc {
		return errNewCacheHasExpired
	}

	if e, exist := m.loadMemoryCacheEntry(key); exist && !e.isExpired() {
		e.Lock()
		e.utcCreated = nowUtc
		e.utcAbsExp = utcAbsExp
		e.value = val
		e.Unlock()
		e.item.Lock()
		e.item.Priority = e.utcAbsExp
		e.item.Unlock()
		m.pq.Update(e.item)
	} else {
		cacheEntity := &memoryCacheEntity{key: key, value: val, utcCreated: nowUtc, utcAbsExp: utcAbsExp}
		item := &Item{Value: key, Priority: cacheEntity.utcAbsExp}
		cacheEntity.item = item
		m.usage.Store(cacheEntity.key, cacheEntity)
		m.pq.PushItem(item)
	}

	return nil
}

// Read returns the value if the key exists in the cache
func (m *memoryCacheStore) Read(key string) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()

	cacheEntity, exist := m.loadMemoryCacheEntry(key)
	if !exist {
		return nil, false
	}

	if cacheEntity.isExpired() {
		return nil, false
	}

	return cacheEntity.value, true
}

// Have returns true if the memory has the item and it's not expired.
func (m *memoryCacheStore) Have(key string) bool {
	_, exist := m.Read(key)
	return exist
}

// Forget removes an item from the memory
func (m *memoryCacheStore) Forget(key string) {
	if e, exist := m.loadMemoryCacheEntry(key); exist {
		m.Lock()
		e.expired()
		m.pq.Update(e.item)
		m.Unlock()
	}
}

// flushExpiredItems remove has expired item from the memory
func (m *memoryCacheStore) flushExpiredItems() {
	num, max, limit := 0, math.MaxInt32, time.Now().UTC().UnixNano()
	for num < max {
		item, yes := m.pq.TryDequeue(limit)
		if !yes {
			break
		}

		if key, ok := item.Value.(string); ok {
			m.usage.Delete(key)
			num++
		}
	}
}
