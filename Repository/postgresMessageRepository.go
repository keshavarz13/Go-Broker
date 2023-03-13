package repository

import (
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/repository/model"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type PostgresMessageRepository struct {
	database     *gorm.DB
	databaseLock sync.Mutex
}

func SetupDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=mysecretpassword")
	if err != nil {
		panic(err)
	}
	db.DropTable(&model.Message{})
	db.CreateTable(&model.Message{})
	return db
}

func CreatePostgresMessageRepository() MessageRepository {
	inMemRepo := PostgresMessageRepository{}
	inMemRepo.database = SetupDB()
	return &inMemRepo
}

func (imr *PostgresMessageRepository) Add(msg broker.Message) (int, error) {
	imr.databaseLock.Lock()
	defer imr.databaseLock.Unlock()

	databaseMessage := convertDtoToDataModel(msg)

	result := imr.database.Create(&databaseMessage) // pass pointer of data to Create

	if result.Error != nil {
		return -1, result.Error
	}
	return databaseMessage.ID, nil
}

func (imr *PostgresMessageRepository) Get(id int) (broker.Message, error) {
	var message = model.Message{ID: id}
	result := imr.database.First(&message)
	if result.Error != nil {
		return broker.Message{}, broker.ErrInvalidID
	}
	return broker.Message{Body: message.Body, Expiration: message.Expiration}, nil
}
