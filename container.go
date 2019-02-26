package robin

import (
    "sync"
)

type Disposable interface {
    Dispose()
}

type container struct {
    sync.Mutex
    container sync.Map
}

func NewContainer() *container {
    return new(container)
}

// Items return container`s all item
func (d *container) Items() []interface{} {
    d.Lock()
    var data []interface{}
    d.container.Range(func(k, v interface{}) bool {
        data = append(data, v)
        return true
    })
    d.Unlock()
    return data
}

// Add put a item into container
func (d *container) Add(item interface{}) {
    d.Lock()
    d.container.Store(item, item)
    d.Unlock()
}

// Remove a item from container
func (d *container) Remove(item interface{}) {
    d.Lock()
    d.container.Delete(item)
    d.Unlock()
}

// Get a item from container if it exist
/*func (d *container) Get(item interface{}) (value interface{}, ok bool) {
    d.Lock()
    defer d.Unlock()
    return d.container.Load(item)
}*/

// Count return container`s number of items
func (d *container) Count() int {
    d.Lock()
    count := 0
    d.container.Range(func(k, v interface{}) bool {
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
    d.container.Range(func(k, v interface{}) bool {
        data = append(data, k)
        v.(Disposable).Dispose()
        return true
    })
    for _, key := range data {
        d.container.Delete(key)
    }
    d.Unlock()
}
