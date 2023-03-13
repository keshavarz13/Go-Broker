package broker

import (
	"context"
	"io"
	"log"
	"time"
)

type Message struct {
	id           int
	creationTime time.Time
	Body         string
	Expiration   time.Duration
}

func (m *Message) IsPersistable() bool {
	if m.Expiration == 0 {
		return false
	}
	return true
}

func (m *Message) SetCreationTime() {
	log.Println("CreationTime:", time.Now().String())
	m.creationTime = time.Now()
}

func (m *Message) IsExpired() bool {
	if m.creationTime.Add(m.Expiration).After(time.Now()) {
		return false
	}
	log.Println("This message expired in ", m.creationTime.Add(m.Expiration).String(), " and current time is:", time.Now().String())
	return true
}

type Broker interface {
	io.Closer
	Publish(ctx context.Context, subject string, msg Message) (int, error)
	Subscribe(ctx context.Context, subject string) (<-chan Message, error)
	Fetch(ctx context.Context, subject string, id int) (Message, error)
}
