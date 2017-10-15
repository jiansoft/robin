package channels

import (
	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/core"
	"github.com/jiansoft/robin/fiber"
)


type IPublisher interface {
	Publish(interface{})
}

//Channel subscription methods.
type ISubscriber interface {
	Subscribe(fiber fiber.Fiber, taskFun interface{}, params ...interface{}) robin.Disposable
	ClearSubscribers()
}

type IChannel interface {
	SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) robin.Disposable
}

type IRequest interface {
	Request() interface{}
	SendReply(replyMsg interface{}) bool
}

type IReply interface {
	Receive(timeoutInMs int, result *interface{}) robin.Disposable
}

type IReplySubscriber interface {
	//Action<IRequest<TR, TM>>
	Subscribe(fiber fiber.Fiber, onRequest *interface{}) robin.Disposable
}

type IQueueChannel interface {
	Subscribe(executionContext core.IExecutionContext, onMessage interface{}) robin.Disposable
	Publish(message interface{})
}
type IProducerThreadSubscriber interface {
	//Allows for the registration and deregistration of fiber. Fiber
	Subscriptions() core.ISubscriptionRegistry
	/*Method called from producer threads*/
	ReceiveOnProducerThread(msg ...interface{})
}
