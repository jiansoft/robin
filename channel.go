package robin

// Channel is a struct that has a member variable to store subscribers
type Channel struct {
	subscribers *container
}

//NewChannel new a Channel instance
func NewChannel() Channel {
	c := Channel{subscribers: newContainer()}
	return c
}

// Subscribe to register a receiver to receive the Channel's message
func (c Channel) Subscribe(taskFun interface{}, params ...interface{}) *Subscriber {
	s := &Subscriber{channel: c, receiver: newTask(taskFun, params...)}
	c.subscribers.Add(s)
	return s
}

// Publish a message to all subscribers
func (c Channel) Publish(msg ...interface{}) {
	items := c.subscribers.Items()
	for _, val := range items {
		if s, ok := val.(*Subscriber); ok {
			fiber.Enqueue(s.receiver.doFunc, msg...)
		}
	}
}

// Clear empty the subscribers
func (c Channel) Clear() {
	items := c.subscribers.Items()
	for _, value := range items {
		c.subscribers.Remove(value)
	}
}

// Count returns a number that how many subscribers in the Channel.
func (c Channel) Count() int {
	return c.subscribers.Count()
}

// Unsubscribe remove the subscriber from the channel
func (c Channel) Unsubscribe(subscriber interface{}) {
	c.subscribers.Remove(subscriber)
}

// Subscriber is a struct for register to a channel
type Subscriber struct {
	channel  Channel
	receiver Task
}

// Unsubscribe remove the subscriber from the channel
func (c *Subscriber) Unsubscribe() {
	c.channel.Unsubscribe(c)
}
