package robin

import (
	"fmt"
)

type channelSubscription struct {
	identifyId string
	fiber      Fiber
	receiver   Task
}

func (c *channelSubscription) init(fiber Fiber, task Task) *channelSubscription {
	c.fiber = fiber
	c.receiver = task
	c.identifyId = fmt.Sprintf("%p-%p", &c, &task)
	return c
}

func NewChannelSubscription(fiber Fiber, task Task) *channelSubscription {
	return new(channelSubscription).init(fiber, task)
}

func (c *channelSubscription) OnMessageOnProducerThread(msg ...interface{}) {
	c.fiber.Enqueue(c.receiver.Func, msg...)
}

//實作 IProducerThreadSubscriber.Subscriptions
func (c *channelSubscription) Subscriptions() SubscriptionRegistry {
	return c.fiber.(SubscriptionRegistry)
}

//實作 IProducerThreadSubscriber.ReceiveOnProducerThread
func (c *channelSubscription) ReceiveOnProducerThread(msg ...interface{}) {
	c.OnMessageOnProducerThread(msg...)
}

func (c *channelSubscription) Dispose() {
	c.fiber.(SubscriptionRegistry).DeregisterSubscription(c)
}

func (c *channelSubscription) Identify() string {
	return c.identifyId
}
