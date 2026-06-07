package usecase

import (
	"context"
	"log"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, value []byte) error {
	log.Printf("received message: %s", value)
	return nil
}
