package usecase

import (
	"context"
	"fmt"

	"producer/internal/domain"
)

type UserUsecase interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
}

type userUsecase struct {
	pub domain.MessagePublisher
}

func NewUserUsecase(pub domain.MessagePublisher) UserUsecase {
	return &userUsecase{pub: pub}
}

func (u *userUsecase) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	if err := u.pub.Publish(ctx, []byte("Key-A"), []byte("Hello Kafka from Go!")); err != nil {
		return nil, fmt.Errorf("publish message: %w", err)
	}

	return &domain.User{
		ID:   id,
		Name: "Taro Yamada",
	}, nil
}
