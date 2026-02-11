package robin

import (
	"testing"
)

// TestDefaultQueue validates enqueue/count/dequeue/dispose behavior of the internal double-buffer queue.
// TestDefaultQueue 驗證內部雙緩衝佇列的 enqueue/count/dequeue/dispose 行為。
func TestDefaultQueue(t *testing.T) {
	type args struct {
		task task
	}
	params := []args{
		{newTask(func(s string) { t.Logf("s:%v", s) }, "enqueue 1")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "enqueue 2")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "enqueue 3")}}

	tests := []struct {
		name string
		args []args
		want int
	}{
		{"TestCount", params, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newDefaultQueue()
			for _, ttt := range tt.args {
				d.enqueue(ttt.task)
			}

			if len(tt.args) != d.count() {
				t.Fatal("they should be equal")
			}

			if got, ok := d.dequeueAll(); ok {
				if tt.want != len(got) {
					t.Fatal("they should be equal")
				}
			}

			if d.count() != 0 {
				t.Fatalf("expected 0, got %d", d.count())
			}

			for _, ttt := range tt.args {
				d.enqueue(ttt.task)
			}

			d.dispose()
		})
	}
}
