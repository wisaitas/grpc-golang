package service

import (
	"context"
	"fmt"

	"github.com/wisaitas/grpc-golang/internal/grpcservice/protogenerate/hello"
)

type HelloService interface {
	SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error)
}

type helloService struct {
	hello.UnimplementedHelloServiceServer
}

func NewHelloService() *helloService {
	return &helloService{}
}

func (s *helloService) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s! สวัสดี!", req.GetName())

	return &hello.HelloResponse{
		Message: message,
	}, nil
}
