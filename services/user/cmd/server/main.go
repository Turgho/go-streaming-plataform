package main

import (
	"context"
	"log"
	"net"
	"os"

	grpchandler "user-service/internal/handler/grpc"
	"user-service/internal/infra/database"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	"user-service/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	db, err := database.NewMongoDatabase(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}

	userRepo := repository.NewUserRepository(db.Client, os.Getenv("MONGO_DB"))
	userUC := usecase.NewUserUseCase(userRepo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpchandler.NewServer(userUC))
	reflection.Register(grpcServer)

	log.Println("user-service rodando na porta :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
