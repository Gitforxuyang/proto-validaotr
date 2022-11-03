package main

import (
	"context"
	"google.golang.org/grpc"
	"multi/post"
	"net"
)

type PostService struct {
}

func (m *PostService) PingPong(ctx context.Context, ping *post.Ping) (*post.Pong, error) {
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
	post.RegisterPostServiceServer(grpcServer, post.NewPostServiceServerImpl(&PostService{}, func(s string) error {
		return Err{msg: s}
	}))
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}
