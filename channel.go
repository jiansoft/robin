package robin

import "sync"

// Channel is a struct that has a member variable to store subscribers
type Channel struct {
	sync.Map
}

//NewChannel new a Channel instance
func NewChannel() *Channel {
	c := &Channel{}
	return c
}

// Subscribe to register a receiver to receive the Channel's message
func (c *Channel) Subscribe(taskFun interface{}, params ...interface{}) *Subscriber {
	s := &Subscriber{channel: c, receiver: newTask(taskFun, params...)}
	c.Store(s, s)
	return s
}

// Publish a message to all subscribers
func (c *Channel) Publish(msg ...interface{}) {
	fiber.Enqueue(func(c *Channel) {
		c.Range(func(k, v interface{}) bool {
			if s, ok := v.(*Subscriber); ok {
				fiber.Enqueue(s.receiver.doFunc, msg...)
			}
			return true
		})
	}, c)
}

// Clear empty the subscribers
func (c *Channel) Clear() {
	c.Range(func(k, v interface{}) bool {
		c.Delete(k)
		return true
	})
}

// Count returns a number that how many subscribers in the Channel.
func (c *Channel) Count() int {
	count := 0
	c.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

// Unsubscribe remove the subscriber from the channel
func (c *Channel) Unsubscribe(subscriber interface{}) {
	c.Delete(subscriber)
}

// Subscriber is a struct for register to a channel
type Subscriber struct {
	channel  *Channel
	receiver Task
}

// Unsubscribe remove the subscriber from the channel
func (c *Subscriber) Unsubscribe() {
	c.channel.Unsubscribe(c)
}
