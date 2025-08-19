.PHONY: generate run clean

generate:
	@echo "Generating proto files..."
	@bash proto/generate.sh

generate-windows:
	@echo "Generating proto files for Windows..."
	protoc --go_out=. --go-grpc_out=. proto/grpcservice/hello.proto
	protoc --go_out=. --go-grpc_out=. proto/grpcservice/pushmessage.proto
	@echo "Proto generation completed!"

run:
	go run cmd/grpcservice/main.go

clean:
	rm -rf internal/grpcservice/protogenerate/*/