package robin

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// equal compares values with DeepEqual and reports caller location on mismatch.
// equal 以 DeepEqual 比較值，若不相等會帶出呼叫端位置方便定位。
func equal(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[37m%s:%d:\n got: %#v\nwant: %#v\033[39m\n ", filepath.Base(file), line, got, want)
		t.FailNow()
	}
}

// TestConcurrentQueueEmpty validates empty-queue semantics for peek/dequeue/len/toArray/clear.
// TestConcurrentQueueEmpty 驗證空佇列在 peek/dequeue/len/toArray/clear 的語義。
func TestConcurrentQueueEmpty(t *testing.T) {
	q := NewConcurrentQueue[string]()

	if v, ok := q.TryPeek(); ok || v != "" {
		t.Errorf("TryPeek() on empty queue = (%v, %v), want (\"\", false)", v, ok)
	}
	if v, ok := q.TryDequeue(); ok || v != "" {
		t.Errorf("TryDequeue() on empty queue = (%v, %v), want (\"\", false)", v, ok)
	}
	if got := q.Len(); got != 0 {
		t.Errorf("Len() on empty queue = %d, want 0", got)
	}
	if arr := q.ToArray(); len(arr) != 0 {
		t.Errorf("ToArray() on empty queue len = %d, want 0", len(arr))
	}
	q.Clear() // should not panic
}

// TestConcurrentStackEmpty validates empty-stack semantics for peek/pop/len/toArray/clear.
// TestConcurrentStackEmpty 驗證空堆疊在 peek/pop/len/toArray/clear 的語義。
func TestConcurrentStackEmpty(t *testing.T) {
	s := NewConcurrentStack[string]()

	if v, ok := s.TryPeek(); ok || v != "" {
		t.Errorf("TryPeek() on empty stack = (%v, %v), want (\"\", false)", v, ok)
	}
	if v, ok := s.TryPop(); ok || v != "" {
		t.Errorf("TryPop() on empty stack = (%v, %v), want (\"\", false)", v, ok)
	}
	if got := s.Len(); got != 0 {
		t.Errorf("Len() on empty stack = %d, want 0", got)
	}
	if arr := s.ToArray(); len(arr) != 0 {
		t.Errorf("ToArray() on empty stack len = %d, want 0", len(arr))
	}
	s.Clear() // should not panic
}

// TestConcurrentBagEmpty validates empty-bag semantics for take/len/toArray/clear.
// TestConcurrentBagEmpty 驗證空 bag 在 take/len/toArray/clear 的語義。
func TestConcurrentBagEmpty(t *testing.T) {
	b := NewConcurrentBag[string]()

	if v, ok := b.TryTake(); ok || v != "" {
		t.Errorf("TryTake() on empty bag = (%v, %v), want (\"\", false)", v, ok)
	}
	if got := b.Len(); got != 0 {
		t.Errorf("Len() on empty bag = %d, want 0", got)
	}
	if arr := b.ToArray(); len(arr) != 0 {
		t.Errorf("ToArray() on empty bag len = %d, want 0", len(arr))
	}
	b.Clear() // should not panic
}

// TestConcurrentQueueShrink verifies automatic shrink from 32 to 16 and data order preservation.
// TestConcurrentQueueShrink 驗證佇列會由 32 縮回 16，且資料順序保持正確。
func TestConcurrentQueueShrink(t *testing.T) {
	q := NewConcurrentQueue[int]()

	// Enqueue 17 items to trigger grow (cap 16 → 32)
	for i := range 17 {
		q.Enqueue(i)
	}

	// Dequeue 10 items: count 17 → 7, which is < cap(32)/4 = 8
	// The 10th dequeue triggers shrink (cap 32 → 16)
	for i := range 10 {
		v, ok := q.TryDequeue()
		if !ok {
			t.Fatalf("TryDequeue failed at i=%d", i)
		}
		equal(t, i, v)
	}

	equal(t, 7, q.Len())

	// Verify remaining elements are intact after shrink
	for i := range 7 {
		v, ok := q.TryDequeue()
		if !ok {
			t.Fatalf("TryDequeue failed at i=%d", i)
		}
		equal(t, i+10, v)
	}
	equal(t, 0, q.Len())
}

// TestConcurrentQueueResizeWrapAround verifies resize correctness when head/tail are wrapped.
// TestConcurrentQueueResizeWrapAround 驗證 head/tail 發生環繞時擴容仍能維持正確順序。
func TestConcurrentQueueResizeWrapAround(t *testing.T) {
	q := NewConcurrentQueue[int]()

	// Step 1: Enqueue 12, dequeue 10 → advance head to 10
	for i := range 12 {
		q.Enqueue(i)
	}
	for range 10 {
		q.TryDequeue()
	}
	// State: head=10, tail=12, count=2

	// Step 2: Enqueue 14 more → tail wraps around, count=16=capacity
	// head=10, tail=10 (wrapped)
	for i := 12; i < 26; i++ {
		q.Enqueue(i)
	}

	// Step 3: One more enqueue triggers grow with wrap-around state
	// resize is called with head >= tail → else branch
	q.Enqueue(26)
	equal(t, 17, q.Len())

	// Verify all elements in correct FIFO order
	for i := range 17 {
		v, ok := q.TryDequeue()
		if !ok {
			t.Fatalf("TryDequeue failed at i=%d", i)
		}
		equal(t, i+10, v)
	}
	equal(t, 0, q.Len())
}

// TestConcurrentQueueClearRetainsCapacityForReuse verifies Clear keeps moderate
// grown capacity so subsequent reuse avoids regrowth allocations.
// TestConcurrentQueueClearRetainsCapacityForReuse 驗證 Clear 在一般放大後會保留容量，
// 以利後續重複使用時避免再次擴容配置。
func TestConcurrentQueueClearRetainsCapacityForReuse(t *testing.T) {
	q := NewConcurrentQueue[int]()

	for i := range 1000 {
		q.Enqueue(i)
	}

	capBefore := len(q.buf)
	if capBefore <= defaultRingCapacity {
		t.Fatalf("expected queue to grow beyond default capacity, got %d", capBefore)
	}

	q.Clear()
	equal(t, 0, q.Len())
	equal(t, capBefore, len(q.buf))

	for i := range 1000 {
		q.Enqueue(i)
	}
	equal(t, capBefore, len(q.buf))
}

// TestConcurrentQueueClearShrinksVeryLargeCapacity verifies Clear shrinks
// excessively large buffers back to default capacity.
// TestConcurrentQueueClearShrinksVeryLargeCapacity 驗證 Clear 會把過大容量回縮到預設值。
func TestConcurrentQueueClearShrinksVeryLargeCapacity(t *testing.T) {
	q := NewConcurrentQueue[int]()

	need := maxRetainedQueueCapacity + 1
	for i := range need {
		q.Enqueue(i)
	}

	if len(q.buf) <= maxRetainedQueueCapacity {
		t.Fatalf("expected capacity > %d, got %d", maxRetainedQueueCapacity, len(q.buf))
	}

	q.Clear()
	equal(t, 0, q.Len())
	equal(t, defaultRingCapacity, len(q.buf))
}

// TestConcurrent validates end-to-end behaviors of queue/stack/bag in normal non-empty flows.
// TestConcurrent 驗證 queue/stack/bag 在一般非空流程下的端到端行為。
func TestConcurrent(t *testing.T) {
	type args struct {
		item string
	}
	items := []args{
		{item: "a"},
		{item: "b"},
		{item: "c"},
		{item: "d"},
		{item: "e"},
	}

	t.Run("ConcurrentQueue", func(t *testing.T) {
		c := NewConcurrentQueue[string]()
		for _, item := range items {
			c.Enqueue(item.item)
		}

		total := c.ToArray()
		equal(t, len(items), len(total))
		for i := range total {
			equal(t, items[i].item, total[i])
		}

		if v, ok := c.TryPeek(); ok {
			equal(t, items[0].item, v)
		}

		for i := 0; ; i++ {
			if v, ok := c.TryDequeue(); !ok {
				break
			} else {
				equal(t, items[i].item, v)
			}
		}

		for i := range 200 {
			c.Enqueue(fmt.Sprintf("item%d", i))
		}

		c.Clear()
		for range 210 {
			c.TryPeek()
			c.TryDequeue()
		}
		equal(t, 0, c.Len())
	})

	t.Run("ConcurrentStack", func(t *testing.T) {
		c := NewConcurrentStack[string]()
		for _, item := range items {
			c.Push(item.item)
		}

		total := c.ToArray()
		equal(t, len(items), len(total))

		if v, ok := c.TryPeek(); ok {
			equal(t, items[len(items)-1].item, v)
		}

		for i := c.Len() - 1; ; i-- {
			if v, ok := c.TryPop(); !ok {
				break
			} else {
				equal(t, items[i].item, v)
			}
		}

		for i := range 200 {
			c.Push(fmt.Sprintf("item%d", i))
		}

		c.Clear()
		for range 210 {
			c.TryPeek()
			c.TryPop()
		}
		equal(t, 0, c.Len())
	})

	t.Run("ConcurrentBag", func(t *testing.T) {
		c := NewConcurrentBag[string]()
		for _, item := range items {
			c.Add(item.item)
		}

		total := c.ToArray()
		equal(t, len(items), len(total))

		for i := len(total) - 1; i >= 0; i-- {
			if v, ok := c.TryTake(); !ok {
				break
			} else {
				equal(t, items[i].item, v)
			}
		}

		for i := range 200 {
			c.Add(fmt.Sprintf("item%d", i))
		}

		for range 210 {
			c.TryTake()
		}

		c.Clear()
		equal(t, 0, c.Len())
	})
}
