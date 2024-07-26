package queue

import (
	"encoding/json"
	"sync"
)

type AtomicQueue[T comparable] struct {
	queue  []T
	locker sync.Locker
}

func (aq *AtomicQueue[T]) Push(value T) {
	aq.locker.Lock()
	defer aq.locker.Unlock()

	aq.queue = append(aq.queue, value)
}

func (aq *AtomicQueue[T]) Remove(value T) bool {
	aq.locker.Lock()
	defer aq.locker.Unlock()

	for i, v := range aq.queue {
		if v == value {
			aq.queue = append(aq.queue[:i], aq.queue[i+1:]...)
			return true
		}
	}

	return false
}

func (aq *AtomicQueue[T]) JSON() (string, error) {
	aq.locker.Lock()
	defer aq.locker.Unlock()

	bs, err := json.Marshal(aq.queue)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func NewAtomicQueue[T comparable]() *AtomicQueue[T] {
	return &AtomicQueue[T]{
		queue:  nil,
		locker: &sync.Mutex{},
	}
}
