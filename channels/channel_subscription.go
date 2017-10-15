package channels

import (
	"fmt"

	"github.com/jiansoft/robin/core"
	"github.com/jiansoft/robin/fiber"
)

type channelSubscription struct {
	identifyId string
	fiber      fiber.Fiber
	receiver   core.Task
}

func (c *channelSubscription) init(fiber fiber.Fiber, task core.Task) *channelSubscription {
	c.fiber = fiber
	c.receiver = task
	c.identifyId = fmt.Sprintf("%p-%p", &c, &task)
	return c
}

func NewChannelSubscription(fiber fiber.Fiber, task core.Task) *channelSubscription {
	return new(channelSubscription).init(fiber, task)
}

func (c *channelSubscription) OnMessageOnProducerThread(msg ...interface{}) {
	c.fiber.Enqueue(c.receiver.Func, msg...)
}

//實作 IProducerThreadSubscriber.Subscriptions
func (c *channelSubscription) Subscriptions() core.ISubscriptionRegistry {
	return c.fiber.(core.ISubscriptionRegistry)
}

//實作 IProducerThreadSubscriber.ReceiveOnProducerThread
func (c *channelSubscription) ReceiveOnProducerThread(msg ...interface{}) {
	c.OnMessageOnProducerThread(msg...)
}

func (c *channelSubscription) Dispose() {
	c.fiber.(core.ISubscriptionRegistry).DeregisterSubscription(c)
}

func (c *channelSubscription) Identify() string {
	return c.identifyId
}
