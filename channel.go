package robin

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Channel is a pub/sub messaging channel with thread-safe subscriber management.
type Channel struct {
	fiber       Fiber
	subscribers sync.Map
	count       atomic.Int64
}

// NewChannel creates a new Channel instance using the default global fiber.
func NewChannel() *Channel {
	return &Channel{fiber: fiber}
}

// Subscribe registers a function to receive the Channel's messages.
func (c *Channel) Subscribe(taskFunc any) *Subscriber {
	s := &Subscriber{channel: c, funcValue: reflect.ValueOf(taskFunc)}
	c.subscribers.Store(s, s)
	c.count.Add(1)
	return s
}

// Publish sends a message to all subscribers.
func (c *Channel) Publish(msg ...any) {
	// Pre-convert message to reflect.Value once, shared across all subscribers.
	params := make([]reflect.Value, len(msg))
	for i, m := range msg {
		params[i] = reflect.ValueOf(m)
	}
	c.fiber.enqueueTask(task{fn: func() {
		c.subscribers.Range(func(k, v any) bool {
			if s, ok := v.(*Subscriber); ok {
				c.fiber.enqueueTask(task{
					funcCache:   s.funcValue,
					paramsCache: params,
				})
			}
			return true
		})
	}})
}

// Clear removes all subscribers from the channel.
func (c *Channel) Clear() {
	c.subscribers.Range(func(k, v any) bool {
		c.subscribers.Delete(k)
		c.count.Add(-1)
		return true
	})
}

// Count returns the number of subscribers in the Channel.
func (c *Channel) Count() int {
	return int(c.count.Load())
}

// Unsubscribe removes a subscriber from the channel.
func (c *Channel) Unsubscribe(subscriber any) {
	if _, loaded := c.subscribers.LoadAndDelete(subscriber); loaded {
		c.count.Add(-1)
	}
}

// Subscriber is a struct for registering to a channel.
type Subscriber struct {
	channel   *Channel
	funcValue reflect.Value
}

// Unsubscribe removes the subscriber from the channel.
func (s *Subscriber) Unsubscribe() {
	s.channel.Unsubscribe(s)
}
