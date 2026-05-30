package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"upload-service/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("uso: client <token> <arquivo.mp4>")
	}

	token := os.Args[1]
	filePath := os.Args[2]

	conn, err := grpc.NewClient("localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUploadServiceClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"authorization", "Bearer "+token,
	)

	stream, err := client.UploadVideo(ctx)
	if err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}

	// abre o arquivo antes de enviar metadados
	// para pegar o tamanho real
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("failed to stat file: %v", err)
	}

	// envia metadados com size real do arquivo
	err = stream.Send(&pb.UploadVideoRequest{
		Data: &pb.UploadVideoRequest_Metadata{
			Metadata: &pb.VideoMetadata{
				Title:       "Meu video",
				Description: "teste",
				Mimetype:    "video/mp4",
				Size:        stat.Size(),
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to send metadata: %v", err)
	}

	// envia chunks de 64KB por vez
	buf := make([]byte, 64*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		if err := stream.Send(&pb.UploadVideoRequest{
			Data: &pb.UploadVideoRequest_Chunk{
				Chunk: buf[:n],
			},
		}); err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}

	fmt.Printf("✓ upload concluído: %s\n", res.Video.Id)
}
