package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"common/ckafka"
	"producer/internal/infrastructure"
	"producer/internal/interface/handler"
	"producer/internal/usecase"
	"proto/pb"

	"google.golang.org/grpc"
)

const defaultBroker = "kafka-service:9092"

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = defaultBroker
	}

	writer, err := ckafka.NewWriter(broker, "test-topic")
	if err != nil {
		log.Fatalf("init kafka writer: %v", err)
	}
	defer writer.Close()

	publisher := infrastructure.NewKafkaPublisher(writer)
	userUC := usecase.NewUserUsecase(publisher)
	userHandler := handler.NewUserHandler(userUC)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("gracefully stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("server stopped")
}
