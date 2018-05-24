package cmap

import (
	"sync"
)

type ConMap struct {
	mu  *sync.Mutex
	kvs map[interface{}]interface{}
}

func NewConMap() *ConMap {
	newMap := new(ConMap)
	newMap.kvs = make(map[interface{}]interface{})
	newMap.mu = new(sync.Mutex)
	return newMap
}

func (this *ConMap) Push(key, value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.kvs[key] = value
}

func (this *ConMap) Pop(key interface{}) interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()
	val, ok := this.kvs[key]
	if ok {
		delete(this.kvs, key)
		return val
	} else {
		return nil
	}
}

func (this *ConMap) Lock() *ConMap {
	this.mu.Lock()
	return this
}

func (this *ConMap) Unlock() *ConMap {
	this.mu.Unlock()
	return this
}

func (this *ConMap) KVs() map[interface{}]interface{} {
	return this.kvs
}
