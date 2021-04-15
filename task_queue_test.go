package robin

import (
	"testing"
)

func TestDefaultQueue(t *testing.T) {
	type args struct {
		task Task
	}
	params := []args{
		{newTask(func(s string) { t.Logf("s:%v", s) }, "Enqueue 1")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "Enqueue 2")},
		{newTask(func(s string) { t.Logf("s:%v", s) }, "Enqueue 3")}}

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
				d.Enqueue(ttt.task)
			}
			if len(tt.args) != d.Count() {
				t.Fatal("they should be equal")
			}
			if got, ok := d.DequeueAll(); ok {
				if tt.want != len(got) {
					t.Fatal("they should be equal")
				}
			}
			if 0 != d.Count() {
				t.Fatal("they should be equal")
			}
			d.Dispose()
		})
	}
}
