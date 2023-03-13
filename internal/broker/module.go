package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/repository"
)

type Module struct {
	queueLock   sync.RWMutex
	activeLock  sync.Mutex
	isClosed    bool
	subscribers map[string][]*chan broker.Message
	repository  repository.MessageRepository
}

func NewModule() broker.Broker {
	module := Module{
		isClosed:    false,
		subscribers: make(map[string][]*chan broker.Message),
		repository:  repository.CreatePostgresMessageRepository(),
	}
	return &module
}

func (m *Module) Close() error {
	m.activeLock.Lock()
	defer m.activeLock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	m.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	id := -1
	eErr := m.checkEnablity()
	if eErr != nil {
		return id, eErr
	}

	if msg.IsPersistable() {
		id, _ = m.repository.Add(msg)
	}
	m.sendToSubscribers(subject, msg)
	return id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	eErr := m.checkEnablity()
	if eErr != nil {
		return nil, eErr
	}
	m.queueLock.Lock()
	defer m.queueLock.Unlock()
	ch := make(chan broker.Message, 100001)
	m.subscribers[subject] = append(m.subscribers[subject], &ch)
	return ch, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	eErr := m.checkEnablity()
	if eErr != nil {
		return broker.Message{}, eErr
	}
	return m.repository.Get(id)
}

func (m *Module) checkEnablity() error {
	m.activeLock.Lock()
	defer m.activeLock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	return nil
}

func (m *Module) sendToSubscribers(subject string, msg broker.Message) {
	m.queueLock.RLock()
	for _, element := range m.subscribers[subject] {
		*element <- msg
	}
	m.queueLock.RUnlock()
}
