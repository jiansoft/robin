package core

import (
	"github.com/jiansoft/robin"
)

//Allows for the registration and deregistration of subscriptions /*The IFiber has implemented*/
type ISubscriptionRegistry interface {
	//Register subscription to be unsubcribed from when the scheduler is disposed.
	RegisterSubscription(robin.Disposable)
	//Deregister a subscription.
	DeregisterSubscription(robin.Disposable)
}

type IExecutionContext interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
}
