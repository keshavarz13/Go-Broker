package repository

import (
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/repository/model"
)

type InMemoryMessageRepository struct {
	databaseLock sync.Mutex
	data         []model.Message
}

func CreateInMemoryMessageRepository() MessageRepository {
	inMemRepo := InMemoryMessageRepository{data: make([]model.Message, 0)}
	return &inMemRepo
}

func (imr *InMemoryMessageRepository) Add(msg broker.Message) (int, error) {
	imr.databaseLock.Lock()
	defer imr.databaseLock.Unlock()
	imr.data = append(imr.data, convertDtoToDataModel(msg))
	return len(imr.data) - 1, nil
}

func (imr *InMemoryMessageRepository) Get(id int) (broker.Message, error) {
	for index, message := range imr.data {
		if index == id {
			if message.IsExpired() {
				return broker.Message{}, broker.ErrExpiredID
			} else {
				return broker.Message{Body: message.Body, Expiration: message.Expiration}, nil
			}
		}
	}
	return broker.Message{}, broker.ErrInvalidID
}
