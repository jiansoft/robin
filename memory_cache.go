package robin

import (
	"math"
	"sync"
	"time"
)

const expirationInterval = 30

type MemoryCache interface {
	Keep(key any, value any, ttl time.Duration)
	Read(key any) (value any, ok bool)
	Have(key any) (ok bool)
	Forget(key any)
	ForgetAll()
}

type memoryCacheEntity struct {
	key        any
	value      any
	item       *Item
	utcCreated int64
	utcAbsExp  int64
	sync.Mutex
}

type memoryCacheStore struct {
	pq    *PriorityQueue
	usage sync.Map
	sync.Mutex
}

var (
	store *memoryCacheStore
	once  sync.Once
)

// Memory returns memoryCacheStore instance.
func Memory() MemoryCache {
	return memoryCache()
}

// memoryCache returns memoryCacheStore singleton instance
func memoryCache() *memoryCacheStore {
	once.Do(func() {
		store = new(memoryCacheStore)
		store.pq = NewPriorityQueue(1024)
		Every(expirationInterval).Seconds().AfterExecuteTask().Do(store.flushExpiredItems)
	})

	return store
}

// isExpired returns true if the item has expired with the mu.
func (mce *memoryCacheEntity) isExpired() bool {
	return time.Now().UTC().UnixNano() > mce.utcAbsExp
}

// expired to set the item has expired with the mu.
func (mce *memoryCacheEntity) expired() {
	mce.Lock()
	mce.utcAbsExp = 0
	mce.Unlock()
	mce.item.expired()
}

// loadMemoryCacheEntry returns memoryCacheEntity if it exists in the cache
func (mcs *memoryCacheStore) loadMemoryCacheEntry(key any) (*memoryCacheEntity, bool) {
	val, yes := mcs.usage.Load(key)
	if !yes {
		return nil, false
	}
	return val.(*memoryCacheEntity), true
}

// Keep insert an item into the memory
func (mcs *memoryCacheStore) Keep(key any, val any, ttl time.Duration) {
	nowUtc := time.Now().UTC().UnixNano()
	utcAbsExp := nowUtc + ttl.Nanoseconds()
	if utcAbsExp <= nowUtc {
		return
	}

	if e, exist := mcs.loadMemoryCacheEntry(key); exist && !e.isExpired() {
		e.Lock()
		e.utcCreated = nowUtc
		e.utcAbsExp = utcAbsExp
		e.value = val
		e.Unlock()
		e.item.setPriority(e.utcAbsExp)
		mcs.Lock()
		mcs.pq.Update(e.item)
		mcs.Unlock()
	} else {
		e = &memoryCacheEntity{key: key, value: val, utcCreated: nowUtc, utcAbsExp: utcAbsExp}
		item := &Item{Value: key, Priority: e.utcAbsExp}
		e.item = item
		mcs.usage.Store(e.key, e)
		mcs.Lock()
		mcs.pq.PushItem(item)
		mcs.Unlock()
	}
}

// Read returns the value if the key exists in the cache
func (mcs *memoryCacheStore) Read(key any) (any, bool) {
	mcs.Lock()
	defer mcs.Unlock()

	cacheEntity, exist := mcs.loadMemoryCacheEntry(key)
	if !exist {
		return nil, false
	}

	if cacheEntity.isExpired() {
		return nil, false
	}

	return cacheEntity.value, true
}

// Have returns true if the memory has the item, and it's not expired.
func (mcs *memoryCacheStore) Have(key any) bool {
	_, exist := mcs.Read(key)
	return exist
}

// Forget removes an item from the memory
func (mcs *memoryCacheStore) Forget(key any) {
	e, exist := mcs.loadMemoryCacheEntry(key)
	if !exist {
		return
	}

	e.expired()
	mcs.Lock()
	mcs.pq.Update(e.item)
	mcs.Unlock()
}

// ForgetAll removes all items from the memory
func (mcs *memoryCacheStore) ForgetAll() {
	mcs.usage.Range(func(k, v any) bool {
		mcs.Forget(k)
		return true
	})
}

// flushExpiredItems remove has expired item from the memory
func (mcs *memoryCacheStore) flushExpiredItems() {
	num, max, limit := 0, math.MaxInt32, time.Now().UTC().UnixNano()
	for num < max {
		mcs.Lock()
		item, yes := mcs.pq.PopItem(limit)
		mcs.Unlock()

		if !yes {
			break
		}

		mcs.usage.Delete(item.Value)
		num++
	}

}
