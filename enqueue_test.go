package robin

import (
	"sync"
	"testing"
)

func TestEnqueue0(t *testing.T) {
	gm := NewGoroutineMulti()
	defer gm.Dispose()

	wg := sync.WaitGroup{}
	wg.Add(1)
	Enqueue0(gm, func() {
		wg.Done()
	})
	wg.Wait()
}

func TestEnqueue1(t *testing.T) {
	gm := NewGoroutineMulti()
	defer gm.Dispose()

	wg := sync.WaitGroup{}
	wg.Add(1)
	Enqueue1(gm, func(msg string) {
		if msg != "hello" {
			t.Errorf("Enqueue1 got %q, want %q", msg, "hello")
		}
		wg.Done()
	}, "hello")
	wg.Wait()
}

func TestEnqueue2(t *testing.T) {
	gm := NewGoroutineMulti()
	defer gm.Dispose()

	wg := sync.WaitGroup{}
	wg.Add(1)
	Enqueue2(gm, func(a int, b string) {
		if a != 42 || b != "world" {
			t.Errorf("Enqueue2 got (%d, %q), want (42, %q)", a, b, "world")
		}
		wg.Done()
	}, 42, "world")
	wg.Wait()
}

func TestEnqueue3(t *testing.T) {
	gm := NewGoroutineMulti()
	defer gm.Dispose()

	wg := sync.WaitGroup{}
	wg.Add(1)
	Enqueue3(gm, func(a int, b string, c bool) {
		if a != 1 || b != "test" || c != true {
			t.Errorf("Enqueue3 got (%d, %q, %v), want (1, %q, true)", a, b, c, "test")
		}
		wg.Done()
	}, 1, "test", true)
	wg.Wait()
}
