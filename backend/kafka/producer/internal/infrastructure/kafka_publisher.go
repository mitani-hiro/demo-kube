package infrastructure

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type kafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(writer *kafka.Writer) *kafkaPublisher {
	return &kafkaPublisher{writer: writer}
}

func (p *kafkaPublisher) Publish(ctx context.Context, key, value []byte) error {
	if err := p.writer.WriteMessages(ctx, kafka.Message{Key: key, Value: value}); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}
	return nil
}
