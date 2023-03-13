package model

import "time"

type Message struct {
	ID             int
	Body           string
	ExpirationTime time.Time
	Expiration     time.Duration
}

func (message *Message) IsExpired() bool {
	return message.ExpirationTime.Before(time.Now())
}
