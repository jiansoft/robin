package robin

type channel struct {
    subscribers *container
}

//NewChannel new a Channel instance
func NewChannel() *channel {
    c := new(channel)
    c.subscribers = NewContainer()
    return c
}

// Subscribe to register a receiver to receive the channel's message
func (c *channel) Subscribe(taskFun interface{}, params ...interface{}) Disposable {
    job := newTask(taskFun, params...)
    subscription := NewChannelSubscription(c, job)
    c.subscribers.Add(subscription)
    return subscription
}

// Publish a message to all subscribers
func (c *channel) Publish(msg ...interface{}) {
    items := c.subscribers.Items()
    for _, val := range items {
        if subscriber, ok := val.(*channelSubscription); ok {
            subscriber.receiver.params(msg...)
            fiber.EnqueueWithTask(subscriber.receiver)
        }
    }
}

//Clear empty the subscribers
func (c *channel) Clear() {
    items := c.subscribers.Items()
    for _, value := range items {
        c.subscribers.Remove(value)
    }
}

func (c *channel) Count() int {
    return c.subscribers.Count()
}

// Remove remove the subscriber
func (c *channel) Remove(subscriber Disposable) {
    c.subscribers.Remove(subscriber)
}

type channelSubscription struct {
    channel  *channel
    receiver Task
}

func NewChannelSubscription(channel *channel, task Task) *channelSubscription {
    c := new(channelSubscription)
    c.channel = channel
    c.receiver = task
    return c
}

func (c *channelSubscription) Dispose() {
    c.channel.Remove(c)
}
