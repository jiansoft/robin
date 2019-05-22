package robin

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var recvCount int32

func recvMsgAction(s string) {
	atomic.AddInt32(&recvCount, 1)
	//log.Printf("recvMsgAction %s recvCount:%v", s, atomic.LoadInt32(&recvCount))
}

func TestChannel(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		{name: "TestChannel", want: 6},
		/*{name: "Test_TestNewChannel_2", want: 3},*/
	}

	channel := NewChannel()
	channel.Subscribe(recvMsgAction)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.SwapInt32(&recvCount, 0)

			subscriber1 := channel.Subscribe(recvMsgAction)
			subscriber2 := channel.Subscribe(recvMsgAction)

			channel.Publish(fmt.Sprintf("Publish message 1 Channel:%v", channel.Count()))
			<-time.After(time.Duration(30) * time.Millisecond)

			subscriber1.Unsubscribe()
			channel.Unsubscribe(subscriber2)
			subscriber1 = nil
			subscriber2 = nil

			channel.Publish(fmt.Sprintf("Publish message 2 Channel:%v", channel.Count()))
			<-time.After(time.Duration(30) * time.Millisecond)

			channel.Subscribe(recvMsgAction)

			channel.Publish(fmt.Sprintf("Publish message 3 Channel:%v", channel.Count()))
			<-time.After(time.Duration(30) * time.Millisecond)

			channel.Clear()
			channel.Publish(fmt.Sprintf("Publish message 4 Channel:%v", channel.Count()))
			<-time.After(time.Duration(30) * time.Millisecond)

			//log.Printf("recvCount %v", atomic.LoadInt32(&recvCount))
			if got := atomic.LoadInt32(&recvCount); got != tt.want {
				t.Errorf("recvCount = %v, want %v", got, tt.want)
			}

			subscriber1 = channel.Subscribe(recvMsgAction)
			subscriber2 = channel.Subscribe(recvMsgAction)
			channel.Subscribe(recvMsgAction)
			if channel.Count() != 3 {
				t.Errorf("the Channel want %v but got %v", 3, channel.Count())
			}

			subscriber1.Unsubscribe()
			if channel.Count() != 2 {
				t.Errorf("the Channel want %v but got %v", 2, channel.Count())
			}

			subscriber2.Unsubscribe()
			if channel.Count() != 1 {
				t.Errorf("the Channel want %v but got %v", 1, channel.Count())
			}

			channel.Unsubscribe(subscriber1)
			channel.Unsubscribe(subscriber2)
			if channel.Count() != 1 {
				t.Errorf("the Channel want %v but got %v", 1, channel.Count())
			}

			channel.Clear()
			if channel.Count() != 0 {
				t.Errorf("the Channel should be empty %v, want %v", channel.Count(), 0)
			}
		})
	}
}

func TestChannelConcurrency(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test_TestConcurrency_1"},
		//{name: "Test_TestConcurrency_2",},
	}
	lock := sync.RWMutex{}
	channel := NewChannel()
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		//atomic.SwapInt32(&recvCount, 0)
		t.Run(tt.name, func(t *testing.T) {
			lock.Lock()
			count := 0
			lock.Unlock()
			loop := 2
			wg.Add(loop * loop)
			for i := 0; i < loop; i++ {
				RightNow().Do(func(c *Channel) {
					for i := 0; i < loop; i++ {
						lock.RLock()
						p := &player{Nickname: strconv.Itoa(count)}
						lock.RUnlock()
						/*d :=*/ c.Subscribe(p.recvMsgAction)
						//c.Publish(fmt.Sprintf("Publish message Channel:%v", c.Count()))
						//d.Unsubscribe()
						//wg.Done()
						lock.Lock()
						count++
						lock.Unlock()
					}

					for i := 0; i < loop; i++ {
						c.Publish(fmt.Sprintf("Publish message Channel:%v", c.Count()))
						wg.Done()
					}
				}, channel)
			}

			wg.Wait()
			//<-time.After(time.Duration(5 * time.Second))
			//log.Printf("Finish %d", atomic.LoadInt32(&recvCount))
		})
	}
}

type player struct {
	Nickname string
}

func (p player) recvMsgAction(msg string) {
	atomic.AddInt32(&recvCount, 1)
	//wg.Add(1)
	//log.Printf("%s receive a message recvCount : %v", p.Nickname, atomic.LoadInt32(&recvCount))
	//wg.Done()
}
