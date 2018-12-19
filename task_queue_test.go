package robin

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

func TestDefaultQueue_init(t *testing.T) {
	type fields struct {
		paddingTasks []Task
		toDoTasks    []Task
		lock         sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   *DefaultQueue
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultQueue{
				paddingTasks: tt.fields.paddingTasks,
				toDoTasks:    tt.fields.toDoTasks,
				lock:         tt.fields.lock,
			}
			if got := d.init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultQueue.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultQueue(t *testing.T) {
	tests := []struct {
		name string
		want *DefaultQueue
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultQueue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultQueue_Dispose(t *testing.T) {
	type fields struct {
		paddingTasks []Task
		toDoTasks    []Task
		lock         sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultQueue{
				paddingTasks: tt.fields.paddingTasks,
				toDoTasks:    tt.fields.toDoTasks,
				lock:         tt.fields.lock,
			}
			d.Dispose()
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
