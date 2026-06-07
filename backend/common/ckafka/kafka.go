package ckafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// NewWriter は broker への疎通を確認した上で Writer を生成する。
// 起動時に broker 不通を検知し fail-fast させるため、ここで kafka.Dial する。
func NewWriter(broker, topic string) (*kafka.Writer, error) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return nil, fmt.Errorf("dial kafka broker %s: %w", broker, err)
	}
	defer conn.Close()

	return &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}, nil
}

// Ping は broker への疎通を確認する。
// broker 起動前に Reader を生成するとグループ参加が復旧不能な状態でスタックするため、
// consumer は起動時に Ping で fail-fast し K8s の再起動に委ねる。
func Ping(broker string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("dial kafka broker %s: %w", broker, err)
	}
	defer conn.Close()

	return nil
}

// NewReader は groupID が空の場合パーティション 0 を直接読む。
// グループ参加はコーディネーター未準備のタイミングで復旧不能にスタックするため、
// 単一 consumer のデモではグループを使わない方が堅牢。
func NewReader(broker, topic, groupID string) *kafka.Reader {
	cfg := kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	}
	if groupID == "" {
		cfg.Partition = 0
		cfg.StartOffset = kafka.FirstOffset
	} else {
		cfg.GroupID = groupID
	}

	return kafka.NewReader(cfg)
}
