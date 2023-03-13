package repository

import (
	"sync"
	"therealbroker/pkg/broker"
)

type PostgresMessageRepository struct {
	databaseLock sync.Mutex
}

func CreatePostgresMessageRepository() MessageRepository {
	inMemRepo := PostgresMessageRepository{}
	return &inMemRepo
}

func (imr *PostgresMessageRepository) Add(msg broker.Message) (int, error) {
	imr.databaseLock.Lock()
	defer imr.databaseLock.Unlock()
	return -1, nil
}

func (imr *PostgresMessageRepository) Get(id int) (broker.Message, error) {
	return broker.Message{}, broker.ErrInvalidID
}
