package main

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"simple"
)

type DemoService struct {
}

func (m *DemoService) PingPong(ctx context.Context, ping *simple.Ping) (*simple.Pong, error) {
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
	simple.RegisterDemoServiceServer(grpcServer, simple.NewDemoServiceServerImpl(&DemoService{}, func(s string) error {
		return Err{msg: s}
	}))
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}
