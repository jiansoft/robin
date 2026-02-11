package robin

import (
	"sync"
	"sync/atomic"
)

// TypedChannel is a generic pub/sub messaging channel that avoids reflection.
// Handler type is func(T), and Publish calls handlers via type assertion instead of reflect.Value.Call().
// TypedChannel 是泛型發布/訂閱訊息頻道，可避免 reflection 呼叫成本。
// 處理器型別固定為 func(T)，Publish 透過型別斷言派送而非 reflect.Value.Call()。
type TypedChannel[T any] struct {
	fiber       Fiber
	subscribers sync.Map
	count       atomic.Int64
}

// TypedSubscriber is the subscription handle bound to a TypedChannel.
// TypedSubscriber 是綁定於 TypedChannel 的訂閱控制代碼。
type TypedSubscriber[T any] struct {
	channel *TypedChannel[T]
	handler func(T)
	invoke  func(any) // cached dispatch: func(a any) { handler(a.(T)) }
}

// NewTypedChannel creates a TypedChannel that uses the package-level default fiber.
// NewTypedChannel 會建立使用套件層級預設 fiber 的 TypedChannel。
func NewTypedChannel[T any]() *TypedChannel[T] {
	return &TypedChannel[T]{fiber: fiber}
}

// NewTypedChannelWithFiber creates a TypedChannel backed by a caller-provided fiber.
// NewTypedChannelWithFiber 會建立使用自訂 fiber 的 TypedChannel。
func NewTypedChannelWithFiber[T any](f Fiber) *TypedChannel[T] {
	return &TypedChannel[T]{fiber: f}
}

// Subscribe registers a handler to receive messages from this channel.
// The invoke function is cached at subscribe time to avoid per-Publish closure allocations.
// Subscribe 會註冊訊息處理器，並在訂閱時快取 invoke 閉包以避免每次 Publish 重複配置。
func (c *TypedChannel[T]) Subscribe(handler func(T)) *TypedSubscriber[T] {
	s := &TypedSubscriber[T]{
		channel: c,
		handler: handler,
		invoke: func(a any) {
			handler(a.(T))
		},
	}
	c.subscribers.Store(s, s)
	c.count.Add(1)
	return s
}

// Publish sends a message to all subscribers.
// Each subscriber's handler is independently enqueued as a task on the fiber,
// matching Channel's semantics (handlers are independently scheduled).
// Uses task.invoke/invokeArg path to avoid per-subscriber closure allocations.
// Publish 會把訊息分派給所有訂閱者，
// 每位訂閱者都會被獨立排程，語義上與 Channel 一致但走無反射路徑。
func (c *TypedChannel[T]) Publish(msg T) {
	var arg any = msg
	c.fiber.enqueueTask(task{fn: func() {
		c.subscribers.Range(func(k, v any) bool {
			if s, ok := v.(*TypedSubscriber[T]); ok {
				c.fiber.enqueueTask(task{invoke: s.invoke, invokeArg: arg})
			}
			return true
		})
	}})
}

// Clear removes all subscribers from the channel.
// Clear 會移除頻道上的全部訂閱者。
func (c *TypedChannel[T]) Clear() {
	c.subscribers.Range(func(k, v any) bool {
		c.subscribers.Delete(k)
		c.count.Add(-1)
		return true
	})
}

// Count returns the current subscriber count.
// Count 回傳目前的訂閱者數量。
func (c *TypedChannel[T]) Count() int {
	return int(c.count.Load())
}

// Unsubscribe removes a subscriber from the channel.
// Unsubscribe 會移除指定訂閱者；重複移除是安全的 no-op。
func (c *TypedChannel[T]) Unsubscribe(subscriber *TypedSubscriber[T]) {
	if _, loaded := c.subscribers.LoadAndDelete(subscriber); loaded {
		c.count.Add(-1)
	}
}

// Unsubscribe removes this subscriber from its parent typed channel.
// Unsubscribe 會將此訂閱者從其所屬的 typed channel 移除。
func (s *TypedSubscriber[T]) Unsubscribe() {
	s.channel.Unsubscribe(s)
}
