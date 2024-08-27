package atomic

import (
	"sync/atomic"
)

type SafeBool struct {
	value atomic.Value
}

func NewSafeBool(initial bool) *SafeBool {
	sb := &SafeBool{}
	sb.value.Store(initial)
	return sb
}

func (sb *SafeBool) Set(value bool) {
	sb.value.Store(value)
}

func (sb *SafeBool) Get() bool {
	return sb.value.Load().(bool)
}
