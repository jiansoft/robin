package robin

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Channel is a pub/sub message bus with thread-safe subscriber management.
// It keeps a copy-on-write subscriber snapshot for lock-free publish reads,
// while subscription updates are guarded by a mutex.
// Channel 是一個發布/訂閱訊息匯流排，具備執行緒安全的訂閱者管理能力。
// 內部以 copy-on-write 快照提供發布讀取無鎖化，訂閱更新則由 mutex 保護。
type Channel struct {
	fiber       Fiber
	mu          sync.Mutex
	subscribers map[*Subscriber]struct{}
	snapshot    atomic.Value // stores []*Subscriber (copy-on-write snapshot)
	count       atomic.Int64
}

// NewChannel creates a Channel that publishes callbacks through the package-level default fiber.
// NewChannel 會建立一個 Channel，並透過套件層級的預設 fiber 來排程回呼執行。
func NewChannel() *Channel {
	c := &Channel{
		fiber:       fiber,
		subscribers: make(map[*Subscriber]struct{}),
	}
	c.snapshot.Store([]*Subscriber{})
	return c
}

// Subscribe registers a callback function and returns a Subscriber handle for later removal.
// The function is stored as reflect.Value so arbitrary signatures can be invoked at publish time.
// Subscribe 會註冊回呼函式，並回傳可後續取消訂閱的 Subscriber 控制代碼。
// 函式會以 reflect.Value 保存，讓 Publish 時可呼叫任意簽章的函式。
func (c *Channel) Subscribe(taskFunc any) *Subscriber {
	s := &Subscriber{channel: c, funcValue: reflect.ValueOf(taskFunc)}
	c.mu.Lock()
	c.ensureSubscribersLocked()
	c.subscribers[s] = struct{}{}
	c.rebuildSnapshotLocked()
	c.count.Add(1)
	c.mu.Unlock()
	return s
}

// Publish asynchronously dispatches one message batch to all subscribers captured
// at Publish call time (call-time snapshot semantics).
// The msg arguments are converted once and shared to reduce per-subscriber reflection overhead.
// Publish 會以非同步方式把同一批訊息分送給 Publish 呼叫當下快照中的訂閱者。
// msg 參數只會先轉換一次後共用，降低每位訂閱者的反射成本。
func (c *Channel) Publish(msg ...any) {
	subs := c.loadSnapshot()
	if len(subs) == 0 {
		return
	}

	// Pre-convert message to reflect.Value once, shared across all subscribers.
	params := make([]reflect.Value, len(msg))
	for i, m := range msg {
		params[i] = reflect.ValueOf(m)
	}

	for _, s := range subs {
		c.fiber.enqueueTask(task{
			funcCache:   s.funcValue,
			paramsCache: params,
		})
	}
}

// Clear removes all subscribers currently registered in the channel.
// Clear 會移除頻道上目前已註冊的所有訂閱者。
func (c *Channel) Clear() {
	c.mu.Lock()
	n := len(c.subscribers)
	if n > 0 {
		clear(c.subscribers)
		c.snapshot.Store([]*Subscriber{})
		c.count.Add(-int64(n))
	}
	c.mu.Unlock()
}

// Count returns the current number of active subscribers.
// Count 回傳目前有效訂閱者的數量。
func (c *Channel) Count() int {
	return int(c.count.Load())
}

// Unsubscribe removes a subscriber handle from the channel.
// Calling this with an already-removed subscriber is safe and is a no-op.
// Unsubscribe 會從頻道移除指定的訂閱者控制代碼。
// 若訂閱者已移除，再次呼叫是安全的 no-op。
func (c *Channel) Unsubscribe(subscriber any) {
	s, ok := subscriber.(*Subscriber)
	if !ok || s == nil {
		return
	}
	c.mu.Lock()
	if _, exists := c.subscribers[s]; exists {
		delete(c.subscribers, s)
		c.rebuildSnapshotLocked()
		c.count.Add(-1)
	}
	c.mu.Unlock()
}

func (c *Channel) loadSnapshot() []*Subscriber {
	if v := c.snapshot.Load(); v != nil {
		if subs, ok := v.([]*Subscriber); ok {
			return subs
		}
	}
	return nil
}

func (c *Channel) ensureSubscribersLocked() {
	if c.subscribers == nil {
		c.subscribers = make(map[*Subscriber]struct{})
	}
}

func (c *Channel) rebuildSnapshotLocked() {
	subs := make([]*Subscriber, 0, len(c.subscribers))
	for s := range c.subscribers {
		subs = append(subs, s)
	}
	c.snapshot.Store(subs)
}

// Subscriber is the subscription handle returned by Channel.Subscribe.
// Use Unsubscribe on this handle to stop receiving future messages.
// Subscriber 是 Channel.Subscribe 回傳的訂閱控制代碼。
// 可透過此代碼呼叫 Unsubscribe 停止接收後續訊息。
type Subscriber struct {
	channel   *Channel
	funcValue reflect.Value
}

// Unsubscribe removes this Subscriber from its parent Channel.
// Unsubscribe 會將此 Subscriber 自其所屬 Channel 移除。
func (s *Subscriber) Unsubscribe() {
	s.channel.Unsubscribe(s)
}
