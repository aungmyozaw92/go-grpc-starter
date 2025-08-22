package main

import (
	"log"
	"net"

	"github.com/aungmyozaw92/go-grpc-starter/config"
	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"github.com/aungmyozaw92/go-grpc-starter/internal/infrastructure"
	grpcHandler "github.com/aungmyozaw92/go-grpc-starter/internal/interface/grpc"
	"github.com/aungmyozaw92/go-grpc-starter/internal/usecase"
	"github.com/aungmyozaw92/go-grpc-starter/proto/userpb"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := cfg.ConnectDatabase()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	repo := infrastructure.NewUserRepository(db)
	uc := usecase.NewUserUseCase(repo)
	handler := grpcHandler.NewUserHandler(uc)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, handler)

	log.Printf("gRPC server running on %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
