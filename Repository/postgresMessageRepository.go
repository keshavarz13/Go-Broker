package repository

import (
	"strconv"
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/repository/model"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresMessageRepository struct {
	database     *gorm.DB
	messagePool  []ReadyToInsertMessage
	messageId    int
	databaseLock sync.Mutex
	poolLock     sync.Mutex
}

type ReadyToInsertMessage struct {
	message model.Message
	ch      *chan bool
}

func SetupDB() *gorm.DB {
	dsn := "host=localhost user=mohammadali password=gorm dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Migrator().DropTable(&model.Message{})
	db.Migrator().CreateTable(&model.Message{})

	return db
}

func CreatePostgresMessageRepository() MessageRepository {
	pgRepo := PostgresMessageRepository{}
	pgRepo.database = SetupDB()
	pgRepo.messageId = 1
	go pgRepo.flushMessagePooInTimePeriods()
	return &pgRepo
}

func (imr *PostgresMessageRepository) Add(msg broker.Message) (int, error) {
	databaseMessage := convertDtoToDataModel(msg)
	databaseMessage = imr.setIdToMessage(databaseMessage)
	ch := imr.sendMessageToMessagePool(databaseMessage)
	go imr.checkBulkCondistionAndTryToInsert(false)
	<-*ch
	return databaseMessage.ID, nil
}

func (imr *PostgresMessageRepository) Get(id int) (broker.Message, error) {
	var message = model.Message{MessageId: id}
	result := imr.database.First(&message)

	if result.Error != nil {
		return broker.Message{}, broker.ErrInvalidID
	}
	if message.IsExpired() {
		return broker.Message{}, broker.ErrExpiredID
	} else {
		return broker.Message{Body: message.Body, Expiration: message.Expiration}, nil
	}
}

func (imr *PostgresMessageRepository) setIdToMessage(msg model.Message) model.Message {
	imr.databaseLock.Lock()
	defer imr.databaseLock.Unlock()
	msg.MessageId = imr.messageId
	msg.Body = msg.Body + strconv.Itoa(imr.messageId)
	imr.messageId++
	return msg
}

func (imr *PostgresMessageRepository) sendMessageToMessagePool(msg model.Message) *chan bool {
	ch := make(chan bool, 1)
	imr.poolLock.Lock()
	defer imr.poolLock.Unlock()
	imr.messagePool = append(imr.messagePool, ReadyToInsertMessage{ch: &ch, message: msg})
	return &ch
}

func (imr *PostgresMessageRepository) checkBulkCondistionAndTryToInsert(checkSizeDisable bool) {
	imr.poolLock.Lock()

	if len(imr.messagePool) >= 100 || checkSizeDisable {
		messages := imr.exportCandidateListForInsertion()
		channels := imr.exportWatedChannels()
		imr.messagePool = []ReadyToInsertMessage{}
		imr.poolLock.Unlock()
		imr.database.CreateInBatches(messages, len(messages))
		for _, ch := range channels {
			*ch <- true
		}
	} else {
		imr.poolLock.Unlock()
	}
	return
}

func (imr *PostgresMessageRepository) exportCandidateListForInsertion() []model.Message {
	messages := []model.Message{}
	for _, messageInfo := range imr.messagePool {
		messages = append(messages, messageInfo.message)
	}
	return messages
}

func (imr *PostgresMessageRepository) exportWatedChannels() []*chan bool {
	channels := []*chan bool{}
	for _, messageInfo := range imr.messagePool {
		channels = append(channels, messageInfo.ch)
	}
	return channels
}

func (imr *PostgresMessageRepository) flushMessagePooInTimePeriods() {
	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			imr.checkBulkCondistionAndTryToInsert(true)
		}
	}

}
