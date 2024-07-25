package counter

import "sync/atomic"

type AtomicCounter struct {
	count int64
}

func (ac *AtomicCounter) Increment() {
	atomic.AddInt64(&ac.count, 1)
}

func (ac *AtomicCounter) Decrement() {
	atomic.AddInt64(&ac.count, -1)
}

func (ac *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&ac.count)
}
