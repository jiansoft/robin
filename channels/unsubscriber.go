package channels

import (
	"fmt"

	"github.com/jiansoft/robin/core"
)

type unsubscriber struct {
	identifyId string
	channel    *channel
	receiver   core.Task
	fiber      core.SubscriptionRegistry
}

func (u *unsubscriber) init(receiver core.Task, channel *channel, fiber core.SubscriptionRegistry) *unsubscriber {
	u.identifyId = fmt.Sprintf("%p-%p", &u, &u.channel)
	u.fiber = fiber
	u.receiver = receiver
	u.channel = channel
	return u
}

func NewUnsubscriber(receiver core.Task, channel *channel, fiber core.SubscriptionRegistry) *unsubscriber {
	return new(unsubscriber).init(receiver, channel, fiber)
}

func (u *unsubscriber) Dispose() {
	u.channel.unsubscribe(u)
}

func (u *unsubscriber) Identify() string {
	return u.identifyId
}
