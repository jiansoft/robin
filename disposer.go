package robin

import (
	"sync"
)

type Disposable interface {
	Dispose()
	Identify() string
}

type Disposer struct {
	sync.Mutex
	sync.Map
}

func (d *Disposer) init() *Disposer {
	return d
}

func NewDisposer() *Disposer {
	return new(Disposer).init()
}

func (d *Disposer) Add(disposable Disposable) {
	d.Lock()
	defer d.Unlock()
	d.Store(disposable.Identify(), disposable)
}

func (d *Disposer) Remove(disposable Disposable) {
	d.Lock()
	defer d.Unlock()
	d.Delete(disposable.Identify())
}

func (d *Disposer) Count() int {
	d.Lock()
	defer d.Unlock()
	count:=0
	d.Range(func(k, v interface{}) bool {
		count++
		return true
	})

	return count
}

func (d *Disposer) Dispose() {
	d.Lock()
	defer d.Unlock()
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
