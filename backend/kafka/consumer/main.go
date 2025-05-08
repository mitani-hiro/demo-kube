package main

import (
	"consumer/interface/app"
	"context"
	"fmt"

	"common/ckafka"
)

func main() {
	fmt.Println("start kafka consumer")

	reader := ckafka.GetKafkaReader()
	defer reader.Close()

	for {
		fmt.Println("start kafka read message")

		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("failed to read message: %v\n", err)
			continue
		}

		app.Hoge(m.Value)
	}
}
