package robin

import (
	"sync"
	"sync/atomic"
)

// TypedChannel is a generic pub/sub messaging channel that avoids reflection.
// Handler type is func(T), and Publish calls handlers via type assertion instead of reflect.Value.Call().
type TypedChannel[T any] struct {
	fiber       Fiber
	subscribers sync.Map
	count       atomic.Int64
}

// TypedSubscriber is a subscriber bound to a TypedChannel.
type TypedSubscriber[T any] struct {
	channel *TypedChannel[T]
	handler func(T)
	invoke  func(any) // cached dispatch: func(a any) { handler(a.(T)) }
}

// NewTypedChannel creates a new TypedChannel using the default global fiber.
func NewTypedChannel[T any]() *TypedChannel[T] {
	return &TypedChannel[T]{fiber: fiber}
}

// NewTypedChannelWithFiber creates a new TypedChannel using a custom fiber.
func NewTypedChannelWithFiber[T any](f Fiber) *TypedChannel[T] {
	return &TypedChannel[T]{fiber: f}
}

// Subscribe registers a handler to receive messages from this channel.
// The invoke function is cached at subscribe time to avoid per-Publish closure allocations.
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
func (c *TypedChannel[T]) Clear() {
	c.subscribers.Range(func(k, v any) bool {
		c.subscribers.Delete(k)
		c.count.Add(-1)
		return true
	})
}

// Count returns the number of subscribers in the channel.
func (c *TypedChannel[T]) Count() int {
	return int(c.count.Load())
}

// Unsubscribe removes a subscriber from the channel.
func (c *TypedChannel[T]) Unsubscribe(subscriber *TypedSubscriber[T]) {
	if _, loaded := c.subscribers.LoadAndDelete(subscriber); loaded {
		c.count.Add(-1)
	}
}

// Unsubscribe removes the subscriber from the channel it belongs to.
func (s *TypedSubscriber[T]) Unsubscribe() {
	s.channel.Unsubscribe(s)
}
