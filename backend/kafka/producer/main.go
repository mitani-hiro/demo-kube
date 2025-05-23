package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"common/ckafka"
	"producer/internal/interface/handler"
	"producer/internal/usecase"
	"proto/pb"

	"google.golang.org/grpc"
)

func main() {
	if err := ckafka.InitKafka(); err != nil {
		fmt.Printf("Failed to initialize Kafka: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userUC := usecase.NewUserUsecase()
	userHandler := handler.NewUserHandler(userUC)
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// サーバー起動を別 goroutine で
	go func() {
		fmt.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// シグナルをキャッチ
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Gracefully stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
