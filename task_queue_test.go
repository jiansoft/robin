package robin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultQueue_Dispose(t *testing.T) {
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
		{"TestDispose", params, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultQueue()
			for _, ttt := range tt.args {
				d.Enqueue(ttt.task)
			}
			assert.Equal(t, tt.want, d.Count(), "they should be equal")
			d.Dispose()
			assert.Equal(t, 0, len(d.paddingTasks), "they should be equal")
			assert.Equal(t, 0, len(d.toDoTasks), "they should be equal")
		})
	}
}

func TestDefaultQueue_Enqueue(t *testing.T) {
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
		{"TestEnqueue", params, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultQueue()
			for _, ttt := range tt.args {
				d.Enqueue(ttt.task)
			}
			assert.Equal(t, tt.want, d.Count(), "they should be equal")
			if got := d.Count(); got != tt.want {
				t.Errorf("DefaultQueue.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueue_DequeueAll(t *testing.T) {
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
		{"TestDequeueAll", params, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultQueue()
			for _, ttt := range tt.args {
				d.Enqueue(ttt.task)
			}
			got, ok := d.DequeueAll()
			if ok {
				assert.Equal(t, tt.want, len(got), "they should be equal")
			} else {
				t.Errorf("DequeueAll test err %+v", got)
			}
			assert.Equal(t, 0, d.Count(), "they should be equal")
		})
	}
}

func TestDefaultQueue_Count(t *testing.T) {
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
			d := NewDefaultQueue()
			for _, ttt := range tt.args {
				d.Enqueue(ttt.task)
			}
			assert.Equal(t, len(tt.args), d.Count(), "they should be equal")
			got, ok := d.DequeueAll()
			if ok {
				assert.Equal(t, tt.want, len(got), "they should be equal")
			}
			assert.Equal(t, 0, d.Count(), "they should be equal")
		})
	}
}
