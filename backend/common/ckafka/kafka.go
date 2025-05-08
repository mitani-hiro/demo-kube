package ckafka

import (
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	topic   = "test-topic"
	groupID = "test-consumer-group"
)

var KafkaWriter *kafka.Writer

func InitKafka() error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	fmt.Printf("InitKafka kafkaBroker: %s\n", kafkaBroker)

	// トピックの存在確認
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		return fmt.Errorf("failed to connect to kafka: %v", err)
	}
	defer conn.Close()

	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return nil
}

func GetKafkaReader() *kafka.Reader {
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	fmt.Printf("InitKafka kafkaBroker: %s\n", kafkaBroker)

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  "",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

}
