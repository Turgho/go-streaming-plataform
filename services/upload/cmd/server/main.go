package main

import (
	"context"
	"log"
	"net"
	"os"
	"upload-service/internal/infra/database"
	"upload-service/internal/infra/interceptor"
	"upload-service/internal/repository"
	"upload-service/internal/server"
	"upload-service/pkg/pb"

	userpb "upload-service/pkg/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// banco
	db, err := database.NewMongoDatabase(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Client.Disconnect(ctx)

	// repository
	repo := repository.NewVideoRepository(db.Client)

	// user-service client
	userConn, err := grpc.NewClient("user-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)

	// gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(userClient)),
		grpc.StreamInterceptor(interceptor.AuthStreamInterceptor(userClient)),
	)

	pb.RegisterUploadServiceServer(grpcServer, server.NewServer(repo))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("upload-service rodando na porta :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
