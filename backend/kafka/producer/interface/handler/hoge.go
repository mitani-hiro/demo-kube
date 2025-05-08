package handler

import (
	"common/ckafka"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func Hoge(c *gin.Context) {

	message := kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello Kafka from Go!"),
	}

	if err := ckafka.KafkaWriter.WriteMessages(context.Background(), message); err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
