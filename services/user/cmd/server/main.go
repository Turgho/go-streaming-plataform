package main

import (
	"context"
	"log"
	"net"
	"os"
	"user-service/internal/infra/database"
	"user-service/internal/repository"
	"user-service/internal/server"
	"user-service/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	db, err := database.NewMongoDatabase(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	if err != nil {
		log.Fatalf("failed to init mongo database: %v", err)
	}

	userRepo := repository.NewUserRepository(db.Client)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, server.NewUserServer(userRepo))

	reflection.Register(grpcServer)

	log.Println("server running on port: 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
