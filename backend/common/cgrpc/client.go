package cgrpc

import (
	"google.golang.org/grpc"
)

func NewClient(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
