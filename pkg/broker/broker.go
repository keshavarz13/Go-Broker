package broker

import (
	"context"
	"io"
	"log"
	"time"
)

type Message struct {
	// This parameter is optional. If it's not provided,
	// the Message can't be accessible through Fetch()
	// id is unique per every subject
	id           int
	creationTime time.Time
	// Body of the message
	Body string
	// The time that message can be accessible through Fetch()
	// with the proper Message id
	// 0 when there is no need to keep message ( fire & forget mode )
	Expiration time.Duration
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
	// Publish returns an int as the id of message published.
	// It should preserve the order. So if we are publishing messages
	// A, B and C, all subscribers should get these messages as
	// A, B and C.
	Publish(ctx context.Context, subject string, msg Message) (int, error)

	// Subscribe listens to every publish, and returns the messages to all
	// subscribed clients ( channels ).
	// If the context is cancelled, you have to stop sending messages
	// to this subscriber. Do nothing on time-out
	Subscribe(ctx context.Context, subject string) (<-chan Message, error)

	// Fetch enables us to retrieve a message that is already published, if
	// it's not expired yet.
	Fetch(ctx context.Context, subject string, id int) (Message, error)
}
