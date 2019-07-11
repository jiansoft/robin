package robin

import (
	"errors"
	"math"
	"sync"
	"time"
)

const expirationInterval = 30

type MemoryCache interface {
	Keep(key interface{}, value interface{}, ttl time.Duration) (err error)
	Read(key interface{}) (value interface{}, ok bool)
	Have(key interface{}) (ok bool)
	Forget(key interface{})
}

type memoryCacheEntity struct {
	utcCreated int64
	utcAbsExp  int64
	key        interface{}
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
func (mce *memoryCacheEntity) isExpired() bool {
	mce.Lock()
	defer mce.Unlock()
	return time.Now().UTC().UnixNano() > mce.utcAbsExp
}

// expired to set the item has expired with the locker.
func (mce *memoryCacheEntity) expired() {
	mce.Lock()
	mce.utcAbsExp = 0
	mce.Unlock()
	mce.item.expired()
}

// loadMemoryCacheEntry returns memoryCacheEntity if it exists in the cache
func (mcs *memoryCacheStore) loadMemoryCacheEntry(key interface{}) (*memoryCacheEntity, bool) {
	val, yes := mcs.usage.Load(key)
	if !yes {
		return nil, false
	}
	return val.(*memoryCacheEntity), true
}

// Keep insert an item into the memory
func (mcs *memoryCacheStore) Keep(key interface{}, val interface{}, ttl time.Duration) error {
	nowUtc := time.Now().UTC().UnixNano()
	utcAbsExp := nowUtc + ttl.Nanoseconds()
	if utcAbsExp <= nowUtc {
		return errNewCacheHasExpired
	}

	if e, exist := mcs.loadMemoryCacheEntry(key); exist && !e.isExpired() {
		e.Lock()
		e.utcCreated = nowUtc
		e.utcAbsExp = utcAbsExp
		e.value = val
		e.Unlock()
		e.item.Lock()
		e.item.Priority = e.utcAbsExp
		e.item.Unlock()
		mcs.pq.Update(e.item)
	} else {
		cacheEntity := &memoryCacheEntity{key: key, value: val, utcCreated: nowUtc, utcAbsExp: utcAbsExp}
		item := &Item{Value: key, Priority: cacheEntity.utcAbsExp}
		cacheEntity.item = item
		mcs.usage.Store(cacheEntity.key, cacheEntity)
		mcs.pq.PushItem(item)
	}

	return nil
}

// Read returns the value if the key exists in the cache
func (mcs *memoryCacheStore) Read(key interface{}) (interface{}, bool) {
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

// Have returns true if the memory has the item and it's not expired.
func (mcs *memoryCacheStore) Have(key interface{}) bool {
	_, exist := mcs.Read(key)
	return exist
}

// Forget removes an item from the memory
func (mcs *memoryCacheStore) Forget(key interface{}) {
	if e, exist := mcs.loadMemoryCacheEntry(key); exist {
		mcs.Lock()
		e.expired()
		mcs.pq.Update(e.item)
		mcs.Unlock()
	}
}

// flushExpiredItems remove has expired item from the memory
func (mcs *memoryCacheStore) flushExpiredItems() {
	num, max, limit := 0, math.MaxInt32, time.Now().UTC().UnixNano()
	for num < max {
		item, yes := mcs.pq.TryDequeue(limit)
		if !yes {
			break
		}

		mcs.usage.Delete(item.Value)
		num++
	}
}
