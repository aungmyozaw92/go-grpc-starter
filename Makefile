.PHONY: proto clean build build-client run test-client test-userlist test-crud deps setup-env

# Generate protobuf files
proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=module=github.com/aungmyozaw92/go-grpc-starter --go-grpc_opt=module=github.com/aungmyozaw92/go-grpc-starter proto/user.proto

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

# Build the test client
build-client:
	go build -o bin/client cmd/client/main.go

# Run the server with MySQL
run:
	@echo "Starting gRPC server with MySQL..."
	@if [ ! -f .env ]; then \
		echo "Warning: No .env file found. Run 'make setup-env' first."; \
	fi
	go run cmd/server/main.go

# Run the test client
test-client:
	@echo "Running gRPC test client..."
	go run cmd/client/main.go

# Test user list API
test-userlist:
	@echo "Testing user list API..."
	go run cmd/test_userlist/main.go

# Test CRUD operations
test-crud:
	@echo "Testing CRUD operations..."
	go run cmd/test_crud/main.go

# Test unique constraints
test-unique:
	@echo "Testing unique constraints..."
	go run cmd/test_unique/main.go

# Clean generated files
clean:
	rm -f proto/userpb/*.pb.go
	rm -rf bin/
