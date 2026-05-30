package grpc

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"upload-service/pkg/pb"
)

const videosBasePath = "/data/videos"

// saveStreamToDisk recebe os chunks do gRPC e grava o arquivo original no disco.
func saveStreamToDisk(videoID string, meta *pb.VideoMetadata, recv func() (*pb.UploadVideoRequest, error)) (string, error) {
	exts, _ := mime.ExtensionsByType(meta.Mimetype)
	if len(exts) == 0 {
		return "", fmt.Errorf("mimetype inválido: %s", meta.Mimetype)
	}

	dirPath := filepath.Join(videosBasePath, videoID)
	filePath := filepath.Join(dirPath, "original"+exts[0])

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return "", fmt.Errorf("criar diretório: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("criar arquivo: %w", err)
	}
	defer file.Close()

	for {
		msg, err := recv()
		if errors.Is(err, io.EOF) {
			return filePath, nil
		}
		if err != nil {
			return "", err
		}
		if _, err := file.Write(msg.GetChunk()); err != nil {
			return "", fmt.Errorf("gravar chunk: %w", err)
		}
	}
}
