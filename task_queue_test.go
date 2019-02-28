package robin

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, len(tt.args), d.Count(), "they should be equal")
			if got, ok := d.DequeueAll(); ok {
				assert.Equal(t, tt.want, len(got), "they should be equal")
			}
			assert.Equal(t, 0, d.Count(), "they should be equal")
			d.Dispose()
		})
	}
}
