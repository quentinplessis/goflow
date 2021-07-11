package goflow

import "sync"

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type stringStateMap struct {
	sync.RWMutex
	internal map[string]state
}

func newStringStateMap() *stringStateMap {
	return &stringStateMap{
		internal: make(map[string]state),
	}
}

func (ssm *stringStateMap) Load(key string) (value state, ok bool) {
	ssm.RLock()
	result, ok := ssm.internal[key]
	ssm.RUnlock()
	return result, ok
}

func (ssm *stringStateMap) Store(key string, value state) {
	ssm.Lock()
	ssm.internal[key] = value
	ssm.Unlock()
}

func (ssm *stringStateMap) Range(f func(key string, value state) bool) {
	ssm.Lock()
	for k, v := range ssm.internal {
		if !f(k, v) {
			break
		}
	}
	ssm.Unlock()
}
