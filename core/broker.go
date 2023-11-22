package core

import (
	"sync"
	"sync/atomic"
)

type BrokerId uint32

type Broker struct {
	sync.RWMutex
	subscriptions subscriptions
	Pipe          chan Message
	check         chan Message
	id            BrokerId
}

func NewBroker() *Broker {
	broker := &Broker{
		Pipe:          make(chan Message),
		check:         make(chan Message),
		subscriptions: newNode(),
		id:            BrokerId(atomic.AddUint32(&uid, 1)),
	}
	go broker.connect()
	return broker
}

func (b *Broker) connect() {
	go b.handleMessages()
}

func (b *Broker) Close() {
	unsubscribe(b)
}

func (b *Broker) handleMessages() {
	defer func() {
		close(b.Pipe)
		close(b.check)
	}()

	for msg := range b.check {
		b.RLock()
		match := b.subscriptions.Match(msg.Subjects)
		b.RUnlock()
		if match {
			b.Pipe <- msg
		}
	}
}

func (b *Broker) Publish(tags []string, data string) error {
	return publish(b.id, tags, data)
}

func (b *Broker) Subscribe(subjects []string) error {
	b.Lock()
	defer b.Unlock()
	if len(subjects) == 0 {
		return nil
	}
	subscribe(b)
	b.subscriptions.Add(subjects)
	return nil
}

func (b *Broker) UnSubscribe(subjects []string) error {
	b.Lock()
	defer b.Unlock()
	b.subscriptions.Remove(subjects)
	return nil
}
