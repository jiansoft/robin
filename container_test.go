package robin

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestContainer_Add(t *testing.T) {
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
		{"TestAdd", fields{Mutex: sync.Mutex{}}, args{disposable: newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &container{
				//locker: tt.fields.Mutex,
				//items: sync.Map{},
			}
			//newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
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

func TestContainer_Remove(t *testing.T) {
	g := NewGoroutineSingle()
	g.Start()
	tests := []struct {
		name string
	}{
		{"TestRemove"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &container{
				//locker: sync.Mutex{},
				//items: sync.Map{},
			}
			pending := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
			d.Add(pending)
			t.Logf("now items has count:%v", d.Count())
			d.Remove(pending)
			if d.Count() == 0 {
				t.Logf("Success count:%v", d.Count())
			} else {
				t.Fatalf("Fatal count:%v", d.Count())
			}
		})
	}
}

func TestContainer_Count(t *testing.T) {
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
			d := NewContainer()
			pending := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
			d.Add(pending)

			if got := d.Count(); got != tt.want {
				t.Errorf("items.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_Dispose(t *testing.T) {
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
			d := &container{
				//locker: tt.fields.Mutex,
				// items: sync.Map{},
			}
			pending1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
			pending2 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
			pending3 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
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

func TestContainer_Random(t *testing.T) {
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
			d := &container{}
			for i := 0; i < 10; i++ {
				ii := i
				RightNow().Do(func() {
					pending1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
					pending2 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
					pending3 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
					d.Add(pending1)
					d.Add(pending2)
					d.Add(pending3)
					d.Remove(pending2)
					d.Remove(pending3)
					d.Remove(pending3)

				})

				pending1 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				pending2 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				pending3 := newTimerTask(g.scheduler.(*scheduler), newTask(func() {}), 50, 50)
				d.Add(pending1)
				d.Add(pending2)
				d.Add(pending3)
				d.Remove(pending1)
				d.Remove(pending2)

				RightNow().Do(func() {
					//d.Dispose()
					ss := d.Items()
					t.Logf("%v add now has count:%v", ii, len(ss))
				})
			}
			/* Delay(1000).Do(func() {
			    d.Dispose()
			})*/
			timeout := time.NewTimer(time.Duration(1000) * time.Millisecond)

			select {
			case <-timeout.C:
				t.Logf("dispose now has count:%v", d.Count())
				d.Dispose()
				t.Logf("dispose now has count:%v", d.Count())
			}
		})
	}
}
