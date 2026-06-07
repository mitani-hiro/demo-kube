package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"common/ckafka"
	"consumer/internal/usecase"
)

const defaultBroker = "kafka-service:9092"

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = defaultBroker
	}

	// broker 不通のまま Reader を生成するとスタックするため fail-fast する
	if err := ckafka.Ping(broker); err != nil {
		log.Fatalf("kafka not reachable: %v", err)
	}

	reader := ckafka.NewReader(broker, "test-topic", "test-consumer-group")
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handler := usecase.NewMessageHandler()

	log.Println("start kafka consumer")
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			// ctx キャンセル時の戻りは graceful shutdown とみなし正常終了する。
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Printf("read message: %v", err)
			continue
		}

		if err := handler.HandleMessage(ctx, m.Value); err != nil {
			log.Printf("handle message: %v", err)
		}
	}
	log.Println("consumer stopped")
}
