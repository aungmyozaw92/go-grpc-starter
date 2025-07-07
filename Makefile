.PHONY: proto clean build run deps setup-env

# Generate protobuf files
proto:
	protoc --go_out=. --go-grpc_out=. proto/user.proto
	mkdir -p proto/userpb
	if [ -d "./github.com" ]; then \
		mv ./github.com/aungmyozaw92/go-grpc-starter/proto/userpb/*.go proto/userpb/ && \
		rm -rf ./github.com; \
	fi

# Install dependencies
deps:
	go mod tidy
	go get gorm.io/gorm
	go get gorm.io/driver/mysql
	go get github.com/joho/godotenv

# Setup .env file from template
setup-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		cp config/env.template .env; \
		echo ".env file created! Please edit it with your database credentials."; \
	else \
		echo ".env file already exists."; \
	fi

# Build the application  
build:
	go build -o bin/server cmd/server/main.go

# Run the server with MySQL
run:
	@echo "Starting gRPC server with MySQL..."
	@if [ ! -f .env ]; then \
		echo "Warning: No .env file found. Run 'make setup-env' first."; \
	fi
	go run cmd/server/main.go

# Clean generated files
clean:
	rm -f proto/userpb/*.pb.go
	rm -rf bin/