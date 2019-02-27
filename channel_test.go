package robin

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewChannel(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{name: "Test_TestNewChannel_1", want: 6},
		{name: "Test_TestNewChannel_2", want: 3},
	}

	var recvCount int32

	channel := NewChannel()
	channel.Subscribe(func(s string, test *testing.T) {
		atomic.AddInt32(&recvCount, 1)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.SwapInt32(&recvCount, 0)
			subscribe1 := channel.Subscribe(func(s string, test *testing.T) {
				atomic.AddInt32(&recvCount, 1)
			})

			subscribe2 := channel.Subscribe(func(s string, test *testing.T) {
				atomic.AddInt32(&recvCount, 1)
			})

			channel.Publish(fmt.Sprintf("Publish message 1 channel:%v", channel.Count()), t)
			<-time.After(time.Duration(30) * time.Millisecond)

			subscribe1.Dispose()
			channel.Remove(subscribe2)

			channel.Publish(fmt.Sprintf("Publish message 2 channel:%v", channel.Count()), t)
			<-time.After(time.Duration(30) * time.Millisecond)

			channel.Subscribe(func(s string, test *testing.T) {
				atomic.AddInt32(&recvCount, 1)
			})

			channel.Publish(fmt.Sprintf("Publish message 3 channel:%v", channel.Count()), t)
			<-time.After(time.Duration(30) * time.Millisecond)

			channel.Clear()
			channel.Publish(fmt.Sprintf("Publish message 4 channel:%v", channel.Count()), t)
			<-time.After(time.Duration(30) * time.Millisecond)

			if got := atomic.LoadInt32(&recvCount); got != tt.want {
				t.Errorf("recvCount = %v, want %v", got, tt.want)
			}
		})
	}
}
