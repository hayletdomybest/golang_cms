package atomic

import (
	"sync/atomic"
)

type SafeObject struct {
	value atomic.Value
}

func NewSafeObject(obj interface{}) *SafeObject {
	sb := &SafeObject{}
	sb.value.Store(obj)
	return sb
}

func (sb *SafeObject) Set(value interface{}) {
	sb.value.Store(value)
}

func (sb *SafeObject) Get() interface{} {
	return sb.value.Load()
}
