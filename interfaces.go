package robin

type IPublisher interface {
	Publish(interface{})
}

//Channel subscription methods.
type ISubscriber interface {
	Subscribe(fiber Fiber, taskFun interface{}, params ...interface{}) Disposable
	ClearSubscribers()
}

type IChannel interface {
	SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) Disposable
}

type IRequest interface {
	Request() interface{}
	SendReply(replyMsg interface{}) bool
}

type IReply interface {
	Receive(timeoutInMs int, result *interface{}) Disposable
}

type IReplySubscriber interface {
	Subscribe(fiber Fiber, onRequest *interface{}) Disposable
}

type IQueueChannel interface {
	Subscribe(executionContext ExecutionContext, onMessage interface{}) Disposable
	Publish(message interface{})
}
type IProducerThreadSubscriber interface {
	//Allows for the registration and deregistration of fiber. Fiber
	Subscriptions() SubscriptionRegistry
	/*Method called from producer threads*/
	ReceiveOnProducerThread(msg ...interface{})
}
