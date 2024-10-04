package pubsub

import (
	"context"
	"sync"
)

type PubSub struct {
	lock sync.Mutex

	subs map[string]map[context.Context]chan interface{}

	closed bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		subs: make(map[string]map[context.Context]chan interface{}),
	}
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return
	}

	if _, ok := ps.subs[topic]; !ok {
		ps.subs[topic] = make(map[context.Context]chan interface{})
	}

	for _, ch := range ps.subs[topic] {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (ps *PubSub) Subscribe(topic string, ctx context.Context, bufferSize int) <-chan interface{} {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return nil
	}

	if _, ok := ps.subs[topic]; !ok {
		ps.subs[topic] = make(map[context.Context]chan interface{})
	}

	ch := make(chan interface{}, bufferSize)
	ps.subs[topic][ctx] = ch

	return ch
}

func (ps *PubSub) Unsubscribe(topic string, ctx context.Context) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return
	}

	if _, ok := ps.subs[topic]; !ok {
		return
	}

	if _, ok := ps.subs[topic][ctx]; ok {
		close(ps.subs[topic][ctx])
	}

	delete(ps.subs[topic], ctx)
}

func (ps *PubSub) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.closed = true

	for _, m := range ps.subs {
		for _, ch := range m {
			close(ch)
		}
	}
}
