package domain

import "context"

type User struct {
	ID   uint64
	Name string
}

type MessagePublisher interface {
	Publish(ctx context.Context, key, value []byte) error
}
