package main

import (
	"context"
	"log"
	"net"
	"os"

	grpchandler "upload-service/internal/handler/grpc"
	"upload-service/internal/infra/database"
	"upload-service/internal/infra/interceptor"
	"upload-service/internal/infra/message"
	infraprobe "upload-service/internal/infra/probe"
	"upload-service/internal/repository"
	"upload-service/internal/usecase"
	"upload-service/pkg/pb"

	userpb "upload-service/pkg/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	db, err := database.NewMongoDatabase(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	defer db.Client.Disconnect(ctx)

	repo := repository.NewVideoRepository(db.Client, os.Getenv("MONGO_DB"))

	userConn, err := grpc.NewClient("user-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("user-service: %v", err)
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	serviceKey := os.Getenv("INTERNAL_SERVICE_KEY")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(userClient, serviceKey)),
		grpc.StreamInterceptor(interceptor.AuthStreamInterceptor(userClient)),
	)

	natsClient, err := message.NewNatsClient(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer natsClient.Close()

	probe := infraprobe.New()
	uploadUC := usecase.NewUploadUseCase(repo, probe, natsClient)
	videoUC := usecase.NewVideoUseCase(repo)

	pb.RegisterUploadServiceServer(grpcServer, grpchandler.NewServer(uploadUC, videoUC))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	log.Println("upload-service rodando na porta :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
