package robin

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func equal(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[37m%s:%d:\n got: %#v\nwant: %#v\033[39m\n ", filepath.Base(file), line, got, want)
		t.FailNow()
	}
}

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

		for i := range total {
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
