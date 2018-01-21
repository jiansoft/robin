package robin

import "reflect"

type Task struct {
	Func        interface{}
	Params      []interface{}
	funcCache   reflect.Value
	paramsCache []reflect.Value
}

func NewTask(t interface{}, p ...interface{}) Task {
	task := Task{Func: t, Params: p}
	task.funcCache = reflect.ValueOf(task.Func)
	task.paramsCache = make([]reflect.Value, len(task.Params))
	for k, param := range task.Params {
		task.paramsCache[k] = reflect.ValueOf(param)
	}
	return task
}

func (t Task) Run() {
	if t.funcCache.IsNil() {
		t.funcCache = reflect.ValueOf(t.Func)
	}
	//execFunc := reflect.ValueOf(t.Func)
	if t.paramsCache == nil {
		t.paramsCache = make([]reflect.Value, len(t.Params))
		for k, param := range t.Params {
			t.paramsCache[k] = reflect.ValueOf(param)
		}
	}
	//params := make([]reflect.Value, len(t.Params))
	//for k, param := range t.Params {
	//    params[k] = reflect.ValueOf(param)
	//}
    t.funcCache.Call(t.paramsCache)
	//func(in []reflect.Value) { _ = t.funcCache.Call(in) }(t.paramsCache)
}
