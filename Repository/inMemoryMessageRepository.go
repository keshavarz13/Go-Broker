package Repository

import (
	"sync"
	"therealbroker/pkg/broker"
)

type InMemoryMessageRepository struct {
	databaseLock sync.Mutex
	data         []broker.Message
}

func CreateInMemoryMessageRepository() MessageRepository {
	inMemRepo := InMemoryMessageRepository{data: make([]broker.Message, 0)}
	return &inMemRepo
}

func (imr *InMemoryMessageRepository) Add(msg broker.Message) (int, error) {
	imr.databaseLock.Lock()
	defer imr.databaseLock.Unlock()
	msg.SetCreationTime()
	imr.data = append(imr.data, msg)
	return len(imr.data) - 1, nil
}
func (imr *InMemoryMessageRepository) Get(id int) (broker.Message, error) {
	for index, message := range imr.data {
		if index == id {
			if message.IsExpired() {
				return broker.Message{}, broker.ErrExpiredID
			} else {
				return message, nil
			}
		}
	}
	return broker.Message{}, broker.ErrInvalidID
}
