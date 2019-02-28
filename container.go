package robin

import (
	"sync"
)

type Disposable interface {
	Dispose()
}

type container struct {
	sync.Mutex
	sync.Map
}

func NewContainer() *container {
	return new(container)
}

// Items return items`s all item
func (d *container) Items() []interface{} {
	d.Lock()
	var data []interface{}
	d.Range(func(k, v interface{}) bool {
		data = append(data, v)
		return true
	})
	d.Unlock()
	return data
}

// Add put a item into items
func (d *container) Add(item interface{}) {
	d.Lock()
	d.Store(item, item)
	d.Unlock()
}

// Remove a item from items
func (d *container) Remove(item interface{}) {
	d.Lock()
	d.Delete(item)
	d.Unlock()
}

// Get a item from items if it exist
/*func (d *items) Get(item interface{}) (value interface{}, ok bool) {
    d.Lock()
    defer d.Unlock()
    return d.items.Load(item)
}*/

// Count return items`s number of items
func (d *container) Count() int {
	d.Lock()
	count := 0
	d.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	d.Unlock()
	return count
}

// Dispose
func (d *container) Dispose() {
	d.Lock()
	var data []interface{}
	d.Range(func(k, v interface{}) bool {
		d.Delete(k)
		data = append(data, v)
		return true
	})
	d.Unlock()
	for _, v := range data {
		if d, ok := v.(Disposable); ok {
			d.Dispose()
		}
	}
}
