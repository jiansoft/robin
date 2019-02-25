package robin

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDisposer_Add(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	type fields struct {
		Mutex sync.Mutex
	}
	type args struct {
		disposable Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TestAdd", fields{Mutex: sync.Mutex{}}, args{disposable: newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &disposer{
				//lock: tt.fields.Mutex,
				Map: sync.Map{},
			}
			//newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			d.Add(tt.args.disposable)
			d.Add(tt.args.disposable)
			d.Add(tt.args.disposable)

			assert.Equal(t, 1, d.Count(), "they should be equal")

			if d.Count() == 1 {
				t.Logf("Success count:%v", d.Count())
			} else {
				t.Fatalf("Fatal count:%v", d.Count())
			}
		})
	}
}

func TestDisposer_Remove(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name string
	}{
		{"TestRemove"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &disposer{
				lock: sync.Mutex{},
				Map:  sync.Map{},
			}
			pending := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			d.Add(pending)
			t.Logf("now disposer has count:%v", d.Count())
			d.Remove(pending)
			if d.Count() == 0 {
				t.Logf("Success count:%v", d.Count())
			} else {
				t.Fatalf("Fatal count:%v", d.Count())
			}
		})
	}
}

func TestDisposer_Count(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name string
		want int
	}{
		{"TestCount", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDisposer()
			pending := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			d.Add(pending)

			if got := d.Count(); got != tt.want {
				t.Errorf("disposer.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisposer_Dispose(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	type fields struct {
		Mutex sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"TestDispose", fields{Mutex: sync.Mutex{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
			pending1 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			pending2 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			pending3 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
			d.Add(pending1)
			d.Add(pending2)
			d.Add(pending3)
			d.Add(pending3)
			t.Logf("before dispose has count:%v", d.Count())
			d.Dispose()
			t.Logf("after dispose has count:%v", d.Count())
		})
	}
}

func TestDisposer_Random(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	type fields struct {
		Mutex sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"TestDispose", fields{Mutex: sync.Mutex{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
			for i := 0; i < 10; i++ {
				ii := i
				RightNow().Do(func() {
					pending1 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
					pending2 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
					pending3 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
					d.Add(pending1)
					d.Add(pending2)
					d.Add(pending3)
					d.Remove(pending1)
					d.Remove(pending2)
					d.Remove(pending3)
					d.Remove(pending3)
					t.Logf("%v add now has count:%v", ii, d.Count())
				})

				pending1 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
				pending2 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
				pending3 := newTimerTask(g.scheduler.(*Scheduler), newTask(func() {}), 50, 50)
				d.Add(pending1)
				d.Add(pending2)
				d.Add(pending3)
				d.Remove(pending1)
				d.Remove(pending2)

				RightNow().Do(func() {
					//d.Dispose()
					//t.Logf("dispose now has count:%v", d.Count())
				})
			}
			/* Delay(1000).Do(func() {
			    d.Dispose()
			})*/
			timeout := time.NewTimer(time.Duration(2000) * time.Millisecond)

			select {
			case <-timeout.C:
				t.Logf("dispose now has count:%v", d.Count())
				d.Dispose()
				t.Logf("dispose now has count:%v", d.Count())
			}
		})
	}
}
