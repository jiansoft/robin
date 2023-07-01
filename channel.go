package robin

import "sync"

// Channel is a struct that has a member variable to store subscribers
type Channel struct {
	sync.Map
}

// NewChannel new a Channel instance
func NewChannel() *Channel {
	c := &Channel{}
	return c
}

// Subscribe to register a receiver to receive the Channel's message
func (c *Channel) Subscribe(taskFunc any, params ...any) *Subscriber {
	s := &Subscriber{channel: c, receiver: newTask(taskFunc, params...)}
	c.Store(s, s)
	return s
}

// Publish a message to all subscribers
func (c *Channel) Publish(msg ...any) {
	fiber.Enqueue(func(c *Channel, message []any) {
		c.Range(func(k, v any) bool {
			if s, ok := v.(*Subscriber); ok {
				s.locker.Lock()
				s.receiver.params(msg...)
				fiber.enqueueTask(s.receiver)
				s.locker.Unlock()
			}
			return true
		})
	}, c, msg)
}

// Clear empty the subscribers
func (c *Channel) Clear() {
	c.Range(func(k, v any) bool {
		c.Delete(k)
		return true
	})
}

// Count returns a number that how many subscribers in the Channel.
func (c *Channel) Count() int {
	count := 0
	c.Range(func(k, v any) bool {
		count++
		return true
	})
	return count
}

// Unsubscribe remove the subscriber from the channel
func (c *Channel) Unsubscribe(subscriber any) {
	c.Delete(subscriber)
}

// Subscriber is a struct for register to a channel
type Subscriber struct {
	channel  *Channel
	receiver task
	locker   sync.Mutex
}

// Unsubscribe remove the subscriber from the channel
func (c *Subscriber) Unsubscribe() {
	c.channel.Unsubscribe(c)
}
