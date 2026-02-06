package robin

import (
	"sync"
	"sync/atomic"
)

// Channel is a pub/sub messaging channel with thread-safe subscriber management.
type Channel struct {
	subscribers sync.Map
	count       atomic.Int64
}

// NewChannel creates a new Channel instance.
func NewChannel() *Channel {
	return &Channel{}
}

// Subscribe registers a receiver to receive the Channel's messages.
func (c *Channel) Subscribe(taskFunc any, params ...any) *Subscriber {
	s := &Subscriber{channel: c, receiver: newTask(taskFunc, params...)}
	c.subscribers.Store(s, s)
	c.count.Add(1)
	return s
}

// Publish sends a message to all subscribers.
func (c *Channel) Publish(msg ...any) {
	fiber.Enqueue(func(ch *Channel, message []any) {
		ch.subscribers.Range(func(k, v any) bool {
			if s, ok := v.(*Subscriber); ok {
				s.mu.Lock()
				s.receiver.params(message...)
				fiber.enqueueTask(s.receiver)
				s.mu.Unlock()
			}
			return true
		})
	}, c, msg)
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
	channel  *Channel
	receiver task
	mu       sync.Mutex
}

// Unsubscribe removes the subscriber from the channel.
func (s *Subscriber) Unsubscribe() {
	s.channel.Unsubscribe(s)
}
