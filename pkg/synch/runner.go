package synch

import (
	"reflect"
	"sync"
)

type Executor struct {
	mutex sync.Mutex
	count int
	fn    reflect.Value
}

func (o *Executor) Run(argv ...reflect.Value) []reflect.Value {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.count <= 0 {
		o.count -= 1
		return nil
	} else {
		o.count -= 1
		return o.fn.Call(argv)
	}
}

func Runner(repeat int, fn any) *Executor {
	return &Executor{
		mutex: sync.Mutex{},
		count: repeat,
		fn:    reflect.ValueOf(fn),
	}
}
