package core

import (
	"fmt"
	"sync"
)

const (
	CmdSubscribe = "subscribe"
	CmdPublish   = "publish"
)

var (
	mutex       = sync.RWMutex{}
	uid         uint32
	subscribers = NewSubscribers()
)

type Subscribers struct {
	sync.RWMutex
	brokers map[BrokerId]*MemoryBroker
}

func NewSubscribers() *Subscribers {
	return &Subscribers{
		brokers: make(map[BrokerId]*MemoryBroker),
	}
}

func (s *Subscribers) Add(broker *MemoryBroker) {
	s.Lock()
	defer s.Unlock()
	s.brokers[broker.id] = broker
}

func (s *Subscribers) Unsubscribe(brokerId BrokerId) {
	s.Lock()
	defer s.Unlock()
	delete(s.brokers, brokerId)
}

func (s *Subscribers) Brokers(selfSubscribeMode bool, brokerId BrokerId) []*MemoryBroker {
	lens := len(s.brokers)
	if !selfSubscribeMode {
		lens = lens - 1
	}
	brokers := make([]*MemoryBroker, 0, lens)
	for _, broker := range s.brokers {
		if !selfSubscribeMode && brokerId == broker.id {
			continue
		}
		brokers = append(brokers, broker)
	}
	return brokers
}

type Message struct {
	Command  string   `json:"commands"`
	Subjects []string `json:"subjects,omitempty"`
	Payload  string   `json:"payload,omitempty"`
	Error    string   `json:"error,omitempty"`
}

type HandlerFn func(*MemoryBroker, Message) error
type ServerStartHandler func(port string, errCh chan<- error)

type subscriptions interface {
	Add([]string)
	Remove([]string)
	Match([]string) bool
}

func subscribe(broker *MemoryBroker) {
	subscribers.Add(broker)
}

func unsubscribe(broker *MemoryBroker) {
	subscribers.Unsubscribe(broker.id)
}

func publish(bId BrokerId, subjects []string, payload string) error {
	if len(subjects) == 0 {
		return fmt.Errorf("failed to publish. missing subjects")
	}

	go func() {
		mutex.RLock()
		defer mutex.RUnlock()

		for _, broker := range subscribers.Brokers(false, bId) {
			msg := Message{
				Command:  CmdPublish,
				Subjects: subjects,
				Payload:  payload,
			}
			go func(b *MemoryBroker, msg Message) {
				b.check <- msg
			}(broker, msg)
		}
	}()
	return nil
}
