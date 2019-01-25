package robin

import (
	"sync"
)

type Disposable interface {
	Dispose()
	Identify() string
}

type Disposer struct {
	lock sync.Mutex
	sync.Map
}

func NewDisposer() *Disposer {
	return new(Disposer)
}

func (d *Disposer) Add(disposable Disposable) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.Store(disposable.Identify(), disposable)
}

func (d *Disposer) Remove(disposable Disposable) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.Delete(disposable.Identify())
}

func (d *Disposer) Count() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	count := 0
	d.Range(func(k, v interface{}) bool {
		count++
		return true
	})

	return count
}

func (d *Disposer) Dispose() {
	d.lock.Lock()
	defer d.lock.Unlock()
	var data []string
	d.Range(func(k, v interface{}) bool {
		data = append(data, k.(string))
		v.(Disposable).Dispose()
		return true
	})

	for _, key := range data {
		d.Delete(key)
	}
	data = data[:0]
}
