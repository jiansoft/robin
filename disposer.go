package robin

import (
	"sync"

	"github.com/jiansoft/robin/collections"
)

type Disposable interface {
	Dispose()
	Identify() string
}

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

func (d *Disposer) Add(disposable Disposable) {
	d.disposables.Set(disposable.Identify(), disposable)
}

func (d *Disposer) Remove(disposable Disposable) {
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
			tmp.(Disposable).Dispose()
		}
	}
}
