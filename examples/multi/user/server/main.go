package main

import (
	"context"
	"google.golang.org/grpc"
	"multi/user"
	"net"
)

type UserService struct {
}

func (m *UserService) PingPong(ctx context.Context, ping *user.Ping) (*user.Pong, error) {
	panic("implement me")
}

type Err struct {
	msg string
}

func (e Err) Error() string {
	return e.msg
}

func main() {
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, user.NewUserServiceServerImpl(&UserService{}, func(s string) error {
		return Err{msg: s}
	}))
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}
