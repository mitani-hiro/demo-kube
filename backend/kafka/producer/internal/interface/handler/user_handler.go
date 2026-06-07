package handler

import (
	"context"
	"os"

	"producer/internal/usecase"
	"proto/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) pb.UserServiceServer {
	return &userHandler{usecase: uc}
}

func (h *userHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.usecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user: %v", err)
	}

	// 応答元 Pod を特定できるよう Hostname を載せる（LB 分散の検証用）。
	hostname, _ := os.Hostname()
	return &pb.GetUserResponse{
		Id:       user.ID,
		Name:     user.Name,
		Hostname: hostname,
	}, nil
}
