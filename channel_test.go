package robin

import (
    "testing"
    "time"
)

func TestNewChannel(t *testing.T) {
    tests := []struct {
        name string
    }{
        {"Test_TestNewChannel_1"},
        {"Test_TestNewChannel_2"},
        {"Test_TestNewChannel_3"},
    }

    channel := NewChannel()
    channel.Subscribe(func(s string) { t.Logf("msg:%s", s) })

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            RightNow().Do(func() {
                t.Logf("%s channel have %v subscribe", tt.name, channel.Count())
                subscribe1 := channel.Subscribe(func(s string) {
                    t.Logf("%s Subscribe 1 msg:%s", tt.name, s)
                })

                subscribe2 := channel.Subscribe(func(s string) {
                    t.Logf("%s Subscribe 2 msg:%s", tt.name, s)
                })

                channel.Publish("Publish message 1")
                subscribe1.Dispose()
                channel.Remove(subscribe2)
                channel.Publish("Publish message 2")

                channel.Subscribe(func(s string) {
                    t.Logf("%s Subscribe 3 msg:%s", tt.name, s)
                })
                channel.Publish("Publish message 3")
                channel.Clear()
                channel.Publish("Publish message 4")

            })
            select {
            case <-time.After(time.Duration(500)):
            }
        })
    }
}
