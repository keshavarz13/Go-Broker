package model

import "time"

type Message struct {
	ID             int           `gorm:"primary_key"`
	MessageId      int           `gorm:"not null"`
	Body           string        `gorm:"not null"`
	ExpirationTime time.Time     `gorm:"not null"`
	Expiration     time.Duration `gorm:"not null"`
}

func (message *Message) IsExpired() bool {
	return message.ExpirationTime.Before(time.Now())
}
