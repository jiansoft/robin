package robin

import "reflect"

type Task struct {
	Func   interface{}
	Params []interface{}
}

func NewTask(t interface{}, p ...interface{}) Task {
	return Task{Func: t, Params: p}
}

func (t Task) Run() {
	execFunc := reflect.ValueOf(t.Func)
	params := make([]reflect.Value, len(t.Params))
	for k, param := range t.Params {
		params[k] = reflect.ValueOf(param)
	}
	func(in []reflect.Value) { _ = execFunc.Call(in) }(params)
}
