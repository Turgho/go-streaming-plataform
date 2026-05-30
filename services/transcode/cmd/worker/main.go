package main

import (
	"log"
	"os"

	"transcode-service/internal/handler"
	grpcinterceptor "transcode-service/internal/infra/grpc"
	"transcode-service/internal/infra/ffmpeg"
	"transcode-service/internal/infra/message"
	uploadclient "transcode-service/internal/infra/upload"
	"transcode-service/internal/usecase"

	uploadpb "transcode-service/pkg/uploadpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	uploadConn, err := grpc.NewClient("upload-service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcinterceptor.ServiceKeyUnaryInterceptor()),
	)
	if err != nil {
		log.Fatalf("upload-service: %v", err)
	}
	defer uploadConn.Close()

	uploadGRPC := uploadclient.NewGRPCClient(uploadpb.NewUploadServiceClient(uploadConn))

	natsClient, err := message.NewNatsClient(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer natsClient.Close()

	transcoder := ffmpeg.New()
	transcodeUC := usecase.NewTranscodeUseCase(transcoder, uploadGRPC)
	w := handler.NewWorker(natsClient, transcodeUC)

	if err := w.Start(); err != nil {
		log.Fatalf("worker: %v", err)
	}

	select {}
}
