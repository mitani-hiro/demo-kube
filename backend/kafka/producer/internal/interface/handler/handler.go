package handler

import (
	"context"
	"producer/internal/usecase"
	"proto/pb"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) pb.UserServiceServer {
	return &userHandler{usecase: uc}
}

func (h *userHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.usecase.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}
