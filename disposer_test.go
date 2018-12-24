package robin

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDisposer_Add(t *testing.T) {
	type fields struct {
		Mutex *sync.Mutex
	}
	type args struct {
		disposable Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TestAdd", fields{Mutex: &sync.Mutex{}}, args{disposable: RightNow().Do(func() {})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
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
	type fields struct {
		Mutex *sync.Mutex
	}
	type args struct {
		disposable Disposable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TestRemove", fields{Mutex: &sync.Mutex{}}, args{disposable: RightNow().Do(func() {})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}

			d.Add(tt.args.disposable)
			t.Logf("now disposer has count:%v", d.Count())
			d.Remove(tt.args.disposable)
			if d.Count() == 0 {
				t.Logf("Success count:%v", d.Count())
			} else {
				t.Fatalf("Fatal count:%v", d.Count())
			}
		})
	}
}

func TestDisposer_Count(t *testing.T) {
	type fields struct {
		Mutex *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"TestCount", fields{Mutex: &sync.Mutex{}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
			d.Add(RightNow().Do(func() {}))

			if got := d.Count(); got != tt.want {
				t.Errorf("Disposer.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisposer_Dispose(t *testing.T) {
	type fields struct {
		Mutex *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"TestDispose", fields{Mutex: &sync.Mutex{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
			d.Add(RightNow().Do(func() {}))
			d.Add(RightNow().Do(func() {}))
			d.Add(RightNow().Do(func() {}))
			t.Logf("before dispose has count:%v", d.Count())
			d.Dispose()
			t.Logf("after dispose has count:%v", d.Count())
		})
	}
}

func TestDisposer_Random(t *testing.T) {
	type fields struct {
		Mutex *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"TestDispose", fields{Mutex: &sync.Mutex{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Disposer{
				lock: tt.fields.Mutex,
				Map:  sync.Map{},
			}
			for i := 0; i < 100000; i++ {
				RightNow().Do(func() {
					d.Add(RightNow().Do(func() {}))
					d.Add(RightNow().Do(func() {}))
					d.Add(RightNow().Do(func() {}))
					//t.Logf("%v add now has count:%v", ii,d.Count())
				})
				RightNow().Do(func() {
					d.Dispose()
					//t.Logf("dispose now has count:%v", d.Count())
				})
			}
			Delay(2000).Do(func() {
				d.Dispose()
			})
			timeout := time.NewTimer(time.Duration(2000) * time.Millisecond)

			select {
			case <-timeout.C:
				t.Logf("dispose now has count:%v", d.Count())
			}
		})
	}
}
