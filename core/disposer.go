package core

import (
	"sync"

	"github.com/jiansoft/robin"
	"github.com/jiansoft/robin/collections"
)

type Disposer struct {
	disposables collections.ConcurrentMap
	lock        *sync.Mutex
}

func (d *Disposer) Init() *Disposer {
	d.disposables = collections.NewConcurrentMap()
	d.lock = new(sync.Mutex)
	return d
}

func NewDisposer() *Disposer {
	return new(Disposer).Init()
}

func (d *Disposer) Add(disposable robin.Disposable) {
	d.disposables.Set(disposable.Identify(), disposable)
}

func (d *Disposer) Remove(disposable robin.Disposable) {
	d.disposables.Remove(disposable.Identify())
}

func (d *Disposer) Count() int {
	return d.disposables.Count()
}

func (d *Disposer) Dispose() {
	d.lock.Lock()
	defer d.lock.Unlock()
	for _, key := range d.disposables.Keys() {
		if tmp, ok := d.disposables.Pop(key); ok {
			tmp.(robin.Disposable).Dispose()
		}
	}
}
