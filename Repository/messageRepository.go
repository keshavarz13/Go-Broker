package repository

import (
	"therealbroker/pkg/broker"
	"therealbroker/repository/model"
	"time"
)

type MessageRepository interface {
	Add(msg broker.Message) (int, error)
	Get(id int) (broker.Message, error)
}

func convertDtoToDataModel(inputDto broker.Message) model.Message {
	return model.Message{
		// Body:           inputDto.Body,
		Body:           "foo",
		ExpirationTime: time.Now().Add(inputDto.Expiration),
		Expiration:     inputDto.Expiration,
	}
}
