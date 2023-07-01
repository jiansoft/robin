package robin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			assert.Equal(t, 0, d.count(), "they should be equal")

			for _, ttt := range tt.args {
				d.enqueue(ttt.task)
			}

			d.dispose()
		})
	}
}
