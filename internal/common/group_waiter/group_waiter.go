package group_waiter

import (
	"sync"
)

type Waitable interface {
	Wait()
}

type GroupWaiter struct {
	waitGroup *sync.WaitGroup
}

func New() *GroupWaiter {
	return &GroupWaiter{
		waitGroup: new(sync.WaitGroup),
	}
}

func (gw *GroupWaiter) Do(f func()) {
	gw.waitGroup.Add(1)
	go func() {
		defer gw.waitGroup.Done()
		f()
	}()
}

func (gw *GroupWaiter) Wait() {
	gw.waitGroup.Wait()
}
