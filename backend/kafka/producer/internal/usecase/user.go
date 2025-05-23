package usecase

import (
	"common/ckafka"
	"context"
	"log"
	"producer/internal/domain"

	"github.com/segmentio/kafka-go"
)

type UserUsecase interface {
	GetUser(id uint64) (*domain.User, error)
}

type userUsecase struct{}

func NewUserUsecase() UserUsecase {
	return &userUsecase{}
}

func (u *userUsecase) GetUser(id uint64) (*domain.User, error) {
	message := kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello Kafka from Go!"),
	}

	if err := ckafka.KafkaWriter.WriteMessages(context.Background(), message); err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// 仮データ返却
	return &domain.User{
		ID:   id,
		Name: "Taro Yamada",
	}, nil
}
