package robin

import (
	"math"
	"sync"
	"time"
)

const expirationInterval = 30

type MemoryCach interface {
	Remember(key string, value interface{}, ttl time.Duration)
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
}

type memoryCacheStore struct {
	usage sync.Map
	pq    *PriorityQueue
}

var (
	store *memoryCacheStore
	once  sync.Once
)

// Memory returns memoryCacheStore instance.
func Memory() MemoryCach {
	return memoryCache()
}

// memoryCache
func memoryCache() *memoryCacheStore {
	once.Do(func() {
		if store == nil {
			store = new(memoryCacheStore)
			store.pq = NewPriorityQueue(1024)
			Every(expirationInterval).Seconds().AfterExecuteTask().Do(store.flushExpiredItems)
		}
	})

	return store
}

// hasExpired returns true if the item has expired.
func (m memoryCacheEntity) hasExpired() bool {
	return time.Now().UTC().UnixNano() > m.utcAbsExp
}

// loadMemoryCacheEntry
func (m *memoryCacheStore) loadMemoryCacheEntry(key string) (*memoryCacheEntity, bool) {
	val, ok := m.usage.Load(key)
	if !ok {
		return nil, false
	}
	return val.(*memoryCacheEntity), true
}

// Remember keeps an item into memory
func (m *memoryCacheStore) Remember(key string, val interface{}, ttl time.Duration) {
	utc := time.Now().UTC().UnixNano()
	e := &memoryCacheEntity{key: key, value: val, utcCreated: utc, utcAbsExp: utc + ttl.Nanoseconds()}
	item := &Item{Value: e, Priority: e.utcAbsExp}
	e.item = item
	m.usage.Store(e.key, e)
	m.pq.PushItem(item)
}

// Read
func (m *memoryCacheStore) Read(key string) (interface{}, bool) {
	entity, ok := m.loadMemoryCacheEntry(key)
	if !ok {
		return nil, false
	}

	if entity.hasExpired() {
		return nil, false
	}

	return entity.value, true
}

// Have returns true if memory has the item and it's not expired.
func (m *memoryCacheStore) Have(key string) bool {
	_, ok := m.Read(key)
	return ok
}

// Forget removes an item from memory
func (m *memoryCacheStore) Forget(key string) {
	entity, ok := m.loadMemoryCacheEntry(key)
	if !ok {
		return
	}

	entity.utcAbsExp = 0
	entity.item.Priority = 0
	m.pq.Update(entity.item)
}

// flushExpiredItems remove has expired item from memory
func (m *memoryCacheStore) flushExpiredItems() {
	limit := time.Now().UTC().UnixNano()
	num, max := 0, math.MaxInt32
	for {
		item, ok := m.pq.TryDequeue(limit)
		if !ok || num >= max {
			break
		}
		entity := item.Value.(*memoryCacheEntity)
		m.usage.Delete(entity.key)
		entity.item = nil
		item.Value = nil
		num++
	}
}
