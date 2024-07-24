package rx

import "sync"

type SubscriberID = uint64

type Broadcast[T any] struct {
	subscriberMaxCount uint64

	SubscriberNextID SubscriberID
	Subscribers      []*Subscriber[T]
	Locker           sync.Mutex
}

func (b *Broadcast[T]) Subscribe() *Subscriber[T] {
	b.Locker.Lock()
	defer b.Locker.Unlock()

	s := &Subscriber[T]{
		broadcast: b,

		ID:      b.SubscriberNextID,
		Channel: make(chan T, b.subscriberMaxCount),
	}

	b.SubscriberNextID++
	if b.SubscriberNextID == 0 {
		b.SubscriberNextID = 1
	}
	b.Subscribers = append(b.Subscribers, s)

	return s
}

func (b *Broadcast[T]) Unsubscribe(sub *Subscriber[T]) SubscriberID {
	b.Locker.Lock()
	defer b.Locker.Unlock()

	for i, s := range b.Subscribers {
		if s == sub {
			b.Subscribers = append(b.Subscribers[:i], b.Subscribers[i+1:]...)
			close(s.Channel)
			return s.ID
		}
	}

	return 0
}

func (b *Broadcast[T]) Send(data T) {
	b.Locker.Lock()
	defer b.Locker.Unlock()

	for _, s := range b.Subscribers {
		select {
		case s.Channel <- data:
			break
		default:
			// channel is full, kick out the subscriber
			go func() {
				b.Unsubscribe(s)
			}()
		}
	}
}

type Subscriber[T any] struct {
	broadcast *Broadcast[T]
	Channel   chan T
	ID        SubscriberID
}

func (sub *Subscriber[T]) Unsubscribe() SubscriberID {
	return sub.broadcast.Unsubscribe(sub)
}

func New[T any](subscriberMaxCount uint64) *Broadcast[T] {
	return &Broadcast[T]{
		subscriberMaxCount: subscriberMaxCount,
		Locker:             sync.Mutex{},
		SubscriberNextID:   1,
	}
}
