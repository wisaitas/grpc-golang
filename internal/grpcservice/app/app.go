package app

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/caarlos0/env/v11"
	"github.com/wisaitas/grpc-golang/internal/grpcservice"
	"github.com/wisaitas/grpc-golang/internal/grpcservice/protogenerate/hello"
	"github.com/wisaitas/grpc-golang/internal/grpcservice/protogenerate/pushmessage"
	"github.com/wisaitas/grpc-golang/internal/grpcservice/service"
	"github.com/wisaitas/grpc-golang/pkg"
)

func init() {
	if err := env.Parse(&grpcservice.ENV); err != nil {
		log.Fatalf("failed to parse env: %v", pkg.Error(err))
	}
}

func RunServer() {
	listener, err := net.Listen("tcp", ":"+grpcservice.ENV.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", pkg.Error(err))
	}

	s := grpc.NewServer()

	// Register Hello Service
	helloService := service.NewHelloService()
	hello.RegisterHelloServiceServer(s, helloService)

	// Register Push Message Service
	pushMessageService := service.NewPushMessageService()
	pushmessage.RegisterPushMessageServiceServer(s, pushMessageService)

	reflection.Register(s)

	log.Printf("gRPC server listening on port %s", grpcservice.ENV.Server.Port)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", pkg.Error(err))
	}
}
