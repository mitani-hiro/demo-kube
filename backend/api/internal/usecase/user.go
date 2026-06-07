package usecase

import (
	"context"
	"fmt"

	"proto/pb"

	"google.golang.org/grpc"
)

// UserClient は producer の gRPC クライアントを利用側で抽象化したもの。
// pb.UserServiceClient がこのインターフェースを満たす。
type UserClient interface {
	GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error)
}

type UserUsecase interface {
	CallProducer(ctx context.Context, id uint64) (*pb.GetUserResponse, error)
}

// connFactory は USE_SHARED_CONN=false のときリクエスト毎に新規コネクションを張るための関数。
// 毎回新規接続することで K8s ClusterIP (L4 LB) が Pod を分散させる挙動を再現する。
// nil の場合は sharedClient を使う（コネクション再利用）。
type userUsecase struct {
	sharedClient UserClient
	connFactory  func() (*grpc.ClientConn, error)
}

func NewUserUsecase(sharedClient UserClient, connFactory func() (*grpc.ClientConn, error)) UserUsecase {
	return &userUsecase{sharedClient: sharedClient, connFactory: connFactory}
}

func (u *userUsecase) CallProducer(ctx context.Context, id uint64) (*pb.GetUserResponse, error) {
	if u.connFactory != nil {
		conn, err := u.connFactory()
		if err != nil {
			return nil, fmt.Errorf("create grpc connection: %w", err)
		}
		defer conn.Close()
		return u.call(ctx, pb.NewUserServiceClient(conn), id)
	}

	return u.call(ctx, u.sharedClient, id)
}

func (u *userUsecase) call(ctx context.Context, client UserClient, id uint64) (*pb.GetUserResponse, error) {
	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("call producer GetUser: %w", err)
	}
	return res, nil
}
