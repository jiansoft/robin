package robin

import "fmt"

type channel struct {
	subscribers ConcurrentMap
}

//NewChannel new a Channel instance
func NewChannel() *channel {
	c := new(channel)
	c.subscribers = NewConcurrentMap()
	return c
}

// Subscribe to register a receiver to receive the channel's message
func (c *channel) Subscribe(taskFun interface{}, params ...interface{}) Disposable {
	job := newTask(taskFun, params...)
	subscription := NewChannelSubscription(c, job)
	c.subscribers.Set(subscription.Identify(), subscription)
	return subscription
}

// Publish a message to all subscribers
func (c *channel) Publish(msg ...interface{}) {
	for _, val := range c.subscribers.Items() {
		if subscriber, ok := val.(*channelSubscription); ok {
			fiber.Enqueue(subscriber.receiver.doFunc, msg...)
		}
	}
}

//ClearSubscribers empty the subscribers
func (c *channel) ClearSubscribers() {
	for key := range c.subscribers.Items() {
		c.subscribers.Remove(key)
	}
}

func (c *channel) NumSubscribers() int {
	return c.subscribers.Count()
}

// unsubscribe remove the subscriber
func (c *channel) unsubscribe(disposable Disposable) {
	if val, ok := c.subscribers.Get(disposable.Identify()); ok {
		if subscriber, ok := val.(*channelSubscription); ok {
			c.subscribers.Remove(subscriber.Identify())
		}
	}
}

type channelSubscription struct {
	channel    *channel
	receiver   Task
	identifyId string
}

func NewChannelSubscription(channel *channel, task Task) *channelSubscription {
	c := new(channelSubscription)
	c.channel = channel
	c.receiver = task
	c.identifyId = fmt.Sprintf("%p-%p", &c, &task)
	return c
}

func (c *channelSubscription) Dispose() {
	c.channel.unsubscribe(c)
}

func (c channelSubscription) Identify() string {
	return c.identifyId
}
