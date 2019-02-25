package robin

import (
	"sync"
)

type Disposable interface {
	Dispose()
	Identify() string
}

type disposer struct {
	lock sync.Mutex
	sync.Map
}

func NewDisposer() *disposer {
	return new(disposer)
}

func (d *disposer) Add(disposable Disposable) {
	d.lock.Lock()
	d.Store(disposable, disposable)
	d.lock.Unlock()
}

func (d *disposer) Remove(disposable Disposable) {
	d.lock.Lock()
	d.Delete(disposable)
	d.lock.Unlock()
}

func (d *disposer) Count() int {
	d.lock.Lock()
	count := 0
	d.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	d.lock.Unlock()
	return count
}

func (d *disposer) Dispose() {
	d.lock.Lock()
	var data []interface{}
	d.Range(func(k, v interface{}) bool {
		data = append(data, k)
		v.(Disposable).Dispose()
		return true
	})

	for _, key := range data {
		d.Delete(key)
	}
	d.lock.Unlock()
}
