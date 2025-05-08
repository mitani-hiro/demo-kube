package main

import (
	"common/ckafka"
	"fmt"
	"producer/router"
)

func main() {
	if err := ckafka.InitKafka(); err != nil {
		fmt.Printf("Failed to initialize Kafka: %v\n", err)
		return
	}

	r := router.NewRouter()

	//r.Use(interceptor.Recovery())

	if err := r.Run(":8080"); err != nil {
		fmt.Printf("gin run error: %v\n", err)
		return
	}
}
