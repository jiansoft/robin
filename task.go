package robin

import "reflect"

type task struct {
	doFunc      interface{}
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func newTask(t interface{}, p ...interface{}) task {
	task := task{doFunc: t}
	task.funcCache = reflect.ValueOf(t)
	task.paramsCache = make([]reflect.Value, len(p))
	for k, param := range p {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t task) run() {
	t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}
