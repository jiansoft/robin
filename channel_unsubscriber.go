package robin

import (
	"fmt"
)

type unsubscriber struct {
	identifyId string
	channel    *channel
	receiver   Task
	fiber      SubscriptionRegistry
}

func (u *unsubscriber) init(receiver Task, channel *channel, fiber SubscriptionRegistry) *unsubscriber {
	u.identifyId = fmt.Sprintf("%p-%p", &u, &u.channel)
	u.fiber = fiber
	u.receiver = receiver
	u.channel = channel
	return u
}

func NewUnsubscriber(receiver Task, channel *channel, fiber SubscriptionRegistry) *unsubscriber {
	return new(unsubscriber).init(receiver, channel, fiber)
}

func (u *unsubscriber) Dispose() {
	u.channel.unsubscribe(u)
}

func (u *unsubscriber) Identify() string {
	return u.identifyId
}
