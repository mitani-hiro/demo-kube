package main

import (
	"log"
	"net"

	"policy/infrastructure"
	"policy/internal/interface/handler"
	"policy/internal/interface/repository"
	"policy/internal/usecase"
	"proto/pb"

	"google.golang.org/grpc"
)

func main() {
	// Initialize database
	db := infrastructure.NewDatabase()

	// Initialize repository
	policyRepo := repository.NewPolicyRepository(db)

	// Initialize usecase
	policyUsecase := usecase.NewPolicyUsecase(policyRepo)

	// Initialize handler
	policyHandler := handler.NewPolicyHandler(policyUsecase)

	// Create gRPC server
	server := grpc.NewServer()
	pb.RegisterPolicyServiceServer(server, policyHandler)

	// Listen on port 50052
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Println("Policy service is running on port 50052")
	if err := server.Serve(listener); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}