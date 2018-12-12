package robin

type channel struct {
	subscribers ConcurrentMap
}

func (c *channel) init() *channel {
	c.subscribers = NewConcurrentMap()
	return c
}

func NewChannel() *channel {
	return new(channel).init()
}

func (c *channel) Subscribe(fiber Fiber, taskFun interface{}, params ...interface{}) Disposable {
	job := newTask(taskFun, params...)
	subscription := NewChannelSubscription(fiber, job)
	return c.SubscribeOnProducerThreads(subscription)
}

func (c *channel) SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) Disposable {
	job := newTask(subscriber.ReceiveOnProducerThread)
	return c.subscribeOnProducerThreads(job, subscriber.Subscriptions())
}

func (c *channel) subscribeOnProducerThreads(subscriber Task, fiber SubscriptionRegistry) Disposable {
	unsubscriber := NewUnsubscriber(subscriber, c, fiber)
	//將訂閱者的方法註冊到 IFiber內 ，當 Fiber.Dispose()時，同步將訂閱的方法移除
	fiber.RegisterSubscription(unsubscriber)
	//放到Channel 內的貯列，當 Chanel.Publish 時發布給訂閱的方法
	c.subscribers.Set(unsubscriber.Identify(), unsubscriber)
	return unsubscriber
}

func (c *channel) Publish(msg ...interface{}) {
	for _, val := range c.subscribers.Items() {
		if subscriber, ok := val.(*unsubscriber); ok {
			subscriber.fiber.(Fiber).Enqueue(subscriber.receiver.doFunc, msg...)
		}
	}
}

// Remove the subscriber
func (c *channel) unsubscribe(disposable Disposable) {

	if val, ok := c.subscribers.Get(disposable.Identify()); ok {
		if subscriber, ok := val.(*unsubscriber); ok {
			subscriber.fiber.DeregisterSubscription(subscriber)
			c.subscribers.Remove(disposable.Identify())
		}
	}
}

func (c *channel) NumSubscribers() int {
	return c.subscribers.Count()
}
