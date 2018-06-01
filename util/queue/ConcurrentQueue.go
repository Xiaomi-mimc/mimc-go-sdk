package que

import (
	"container/list"
	"sync"
)

type ConQueue struct {
	mu   *sync.Mutex
	eles *list.List
	size uint32
}

func NewConQueue() *ConQueue {
	newDue := new(ConQueue)
	newDue.eles = list.New()
	newDue.mu = new(sync.Mutex)
	newDue.size = 0
	return newDue
}

func (this *ConQueue) Push(ele interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.eles.PushBack(ele)
	this.size += 1
}

func (this *ConQueue) Pop() interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.size == 0 {
		return nil
	}
	head := this.eles.Front()
	this.eles.Remove(head)
	this.size -= 1
	return head.Value
}

func (this *ConQueue) Size() uint32 {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.size
}
