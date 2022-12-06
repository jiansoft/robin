package robin

import (
	"container/heap"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func equal(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[37m%s:%d:\n got: %#v\nwant: %#v\033[39m\n ", filepath.Base(file), line, got, want)
		t.FailNow()
	}
}

func lessThan(t *testing.T, a, b int64) {
	if a > b {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n a: %#v\b: %#v\033[39m\n ", filepath.Base(file), line, a, b)
		t.FailNow()
	}
}

func TestRemove(t *testing.T) {
	c := 1024
	pq := NewPriorityQueue(c)

	for i := 0; i < c; i++ {
		v := rand.Int()
		heap.Push(pq, &Item{Value: "qq", Priority: int64(v)})
	}

	for i := 0; i < 10; i++ {
		v := rand.Intn(pq.Len() - 1)
		heap.Remove(pq, v)
	}

	lastItem := heap.Pop(pq).(*Item)
	count := len(*pq) - 1
	for i := 0; i < count; i++ {
		item := heap.Pop(pq).(*Item)
		lastItemPriority := lastItem.getPriority()
		itemPriority := item.getPriority()
		equal(t, lastItemPriority < itemPriority, true)
		lastItem = item
	}
}

func TestUpdate(t *testing.T) {
	c := 10
	pq := NewPriorityQueue(c)

	var items []*Item
	for i := 0; i < c; i++ {
		item := &Item{Value: fmt.Sprintf("qq-%v", i), Priority: int64(i)}
		items = append(items, item)
		pq.PushItem(item)
	}

	for i, j := 0, len(items); i < len(items); i, j = i+1, j-1 {
		items[i].setPriority(int64(j))
		pq.Update(items[i])
	}

	lastItem := heap.Pop(pq).(*Item)
	count := len(*pq) - 1
	for i := 0; i < count; i++ {
		item := heap.Pop(pq).(*Item)
		lastItemPriority := lastItem.getPriority()
		itemPriority := item.getPriority()
		equal(t, lastItemPriority < itemPriority, true)
		lastItem = item
	}
}

func TestPush(t *testing.T) {
	c := 10
	pq := NewPriorityQueue(c)

	for i := 0; i < c; i++ {
		v := rand.Int63n(time.Now().UnixNano() + int64(time.Second))
		pq.PushItem(&Item{Value: i, Priority: v})
	}
	equal(t, pq.Len(), c)
	equal(t, cap(*pq), c)
}

func TestPriorityQueue(t *testing.T) {
	pg := NewPriorityQueue(1024)
	running := int32(1)
	tests := []struct {
		pq   *PriorityQueue
		name string
	}{
		{pg, "1"},
		{pg, "2"},
	}
	var mu sync.Mutex
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RightNow().Do(func(pq *PriorityQueue) {
				for atomic.LoadInt32(&running) == 1 {
					item := &Item{Value: time.Now().UnixNano(), Priority: time.Now().UnixNano()}
					mu.Lock()
					pq.PushItem(item)
					mu.Unlock()
				}
			}, tt.pq)

			RightNow().Do(func(pq *PriorityQueue) {
				for atomic.LoadInt32(&running) == 1 {
					limit := time.Now().UnixNano()
					mu.Lock()
					item, ok := pq.PopItem(limit)
					mu.Unlock()
					if !ok {
						continue
					}
					priority := item.getPriority()
					lessThan(t, priority, limit)
				}
			}, tt.pq)

			timeout := time.NewTimer(time.Duration(1500) * time.Millisecond)
			select {
			case <-timeout.C:
			}
			atomic.StoreInt32(&running, 0)
		})
	}
}

func TestArray(t *testing.T) {
	var result = make([]int, 0, 10)
	result = append(result, 1)
	t.Logf("%v", result)
}
