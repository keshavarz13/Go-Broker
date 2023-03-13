package Repository

import (
	"therealbroker/pkg/broker"
)

type MessageRepository interface {
	Add(msg broker.Message) (int, error)
	Get(id int) (broker.Message, error)
}
