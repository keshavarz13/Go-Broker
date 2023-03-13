package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/repository"
)

type Module struct {
	queueLock   sync.Mutex
	isClosed    bool
	subscribers map[string][]*chan broker.Message
	repository  repository.MessageRepository
}

func NewModule() broker.Broker {
	module := Module{
		isClosed:    false,
		subscribers: make(map[string][]*chan broker.Message),
		repository:  repository.CreateInMemoryMessageRepository(),
	}
	return &module
}

func (m *Module) Close() error {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	m.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	id := -1
	if msg.IsPersistable() {
		id, _ = m.repository.Add(msg)
	}
	m.queueLock.Lock()
	defer m.queueLock.Unlock()
	if m.isClosed {
		return id, broker.ErrUnavailable
	}
	for _, element := range m.subscribers[subject] {
		*element <- msg
	}
	return id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()
	if m.isClosed {
		return nil, broker.ErrUnavailable
	}
	ch := make(chan broker.Message, 70)
	m.subscribers[subject] = append(m.subscribers[subject], &ch)
	return ch, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.isClosed {
		return broker.Message{}, broker.ErrUnavailable
	}
	return m.repository.Get(id)
}
