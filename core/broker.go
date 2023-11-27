package core

import (
	"sync"
	"sync/atomic"
)

type BrokerId uint32

type Broker interface {
	Publish([]string, string) error
	Subscribe([]string) error
}

type MemoryBroker struct {
	sync.RWMutex
	subscriptions subscriptions
	Pipe          chan Message
	check         chan Message
	id            BrokerId
}

func NewBroker() *MemoryBroker {
	broker := &MemoryBroker{
		Pipe:          make(chan Message),
		check:         make(chan Message),
		subscriptions: newSubjectNode(),
		id:            BrokerId(atomic.AddUint32(&uid, 1)),
	}
	go broker.connect()
	return broker
}

func (b *MemoryBroker) connect() {
	go b.handleMessages()
}

func (b *MemoryBroker) Close() {
	unsubscribe(b)
}

func (b *MemoryBroker) handleMessages() {
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

func (b *MemoryBroker) Publish(tags []string, data string) error {
	return publish(b.id, tags, data)
}

func (b *MemoryBroker) Subscribe(subjects []string) error {
	b.Lock()
	defer b.Unlock()
	if len(subjects) == 0 {
		return nil
	}
	subscribe(b)
	b.subscriptions.Add(subjects)
	return nil
}

func (b *MemoryBroker) UnSubscribe(subjects []string) error {
	b.Lock()
	defer b.Unlock()
	b.subscriptions.Remove(subjects)
	return nil
}
