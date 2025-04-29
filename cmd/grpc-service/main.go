package main

import (
	"log"
	"os"

	grpcservice "github.com/wisaitas/grpc-golang/internal/grpc-service"
)

func main() {
	port := "50051"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	if err := grpcservice.RunServer(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
