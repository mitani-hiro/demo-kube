package infrastructure

import (
	"fmt"
	"log"
	"net"

	"producer/internal/interface/handler"
	"producer/internal/usecase"
	"proto/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userUC := usecase.NewUserUsecase()
	userHandler := handler.NewUserHandler(userUC)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
