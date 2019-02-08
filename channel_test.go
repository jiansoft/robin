package robin

import (
	"testing"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel.Subscribe(func(s string) {
				t.Logf("Subscribe 1 msg:%s", s)
			})

			subscribe2 := channel.Subscribe(func(s string) {
				t.Logf("Subscribe 2 msg:%s", s)
			})

			t.Logf("the channel have %v subscribers.", channel.NumSubscribers())
			channel.Publish("Publish message 1")

			subscribe2.Dispose()
			t.Logf("the channel have %v subscribers.", channel.NumSubscribers())
			channel.Publish("Publish message 2")

			channel.ClearSubscribers()
			t.Logf("the channel have %v subscribers.", channel.NumSubscribers())
			channel.Publish("Publish message 3")

			channel.Subscribe(func(s string) {
				t.Logf("Subscribe 3 msg:%s", s)
			})
			channel.Publish("Publish message 4")
		})
	}
}
