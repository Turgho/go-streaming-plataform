package server

import (
	"context"
	"io"
	"mime"
	"os"
	"path/filepath"
	"upload-service/internal/domain/entities"
	domain "upload-service/internal/domain/repositories"
	"upload-service/internal/infra/interceptor"
	"upload-service/pkg/pb"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedUploadServiceServer
	repo domain.VideoRepository
}

func NewServer(repo domain.VideoRepository) *Server {
	return &Server{repo: repo}
}

func (s *Server) UploadVideo(stream pb.UploadService_UploadVideoServer) error {
	first, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive metadata: %v", err)
	}

	meta := first.GetMetadata()
	if meta == nil {
		return status.Error(codes.InvalidArgument, "first message must contain metadata")
	}

	// userID do contexto
	userID, ok := stream.Context().Value(interceptor.UserIDKey).(string)
	if !ok {
		return status.Error(codes.Unauthenticated, "user not found in context")
	}

	videoID, err := gonanoid.New()
	if err != nil {
		return status.Error(codes.Internal, "failed to generate video id")
	}

	dirPath := filepath.Join("/data/videos", videoID)
	exts, _ := mime.ExtensionsByType(meta.Mimetype)
	if len(exts) == 0 {
		return status.Errorf(codes.InvalidArgument, "mimetype inválido: %s", meta.Mimetype)
	}

	filePath := filepath.Join(dirPath, "original"+exts[0])
	os.MkdirAll(dirPath, 0755)

	file, err := os.Create(filePath)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create file: %v", err)
	}
	defer file.Close()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive chunk: %v", err)
		}
		if _, err := file.Write(msg.GetChunk()); err != nil {
			return status.Errorf(codes.Internal, "failed to write chunk: %v", err)
		}
	}

	video, err := entities.NewVideo(
		videoID, userID,
		meta.Title, meta.Description, meta.Mimetype,
		filePath, meta.Size, 0, // <- valor é atualizado no ffprobe
	)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "video inválido: %v", err)
	}

	if err := s.repo.Create(stream.Context(), video); err != nil {
		return status.Errorf(codes.Internal, "failed to save video: %v", err)
	}

	// NATS vai aqui depois

	return stream.SendAndClose(&pb.UploadVideoResponse{
		Video: toProtoVideo(video),
	})
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.VideoMetadata, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id é obrigatório")
	}

	video, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoVideo(video), nil
}

func (s *Server) ListByUserID(ctx context.Context, req *pb.ListByUserIDRequest) (*pb.ListByUserIDResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id é obrigatório")
	}

	var protoVideos []*pb.VideoMetadata
	videos, err := s.repo.ListByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, video := range videos {
		protoVideos = append(protoVideos, toProtoVideo(video))
	}

	return &pb.ListByUserIDResponse{
		Videos: protoVideos,
	}, nil
}

func (s *Server) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.VideoMetadata, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id é obrigatório")
	}
	if _, ok := pb.VideoStatus_name[int32(req.Status)]; !ok {
		return nil, status.Error(codes.InvalidArgument, "status inválido")
	}

	if err := s.repo.UpdateStatus(ctx, req.Id, fromProtoStatus(req.Status)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	video, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoVideo(video), nil
}

// Helpers
func toProtoVideo(video *entities.Video) *pb.VideoMetadata {
	return &pb.VideoMetadata{
		Id:          video.ID,
		UserId:      video.UserID,
		Title:       video.Title,
		Description: video.Description,
		Mimetype:    string(video.Mimetype),
		Size:        video.Size,
		Duration:    video.Duration,
		Status:      toProtoStatus(video.Status),
		FilePath:    video.FilePath,
		CreatedAt:   timestamppb.New(video.CreatedAt),
		UpdatedAt:   timestamppb.New(video.UpdatedAt),
	}
}

func toProtoStatus(s entities.VideoStatus) pb.VideoStatus {
	switch s {
	case entities.StatusUploaded:
		return pb.VideoStatus_UPLOADED
	case entities.StatusTranscoding:
		return pb.VideoStatus_TRANSCODING
	case entities.StatusReady:
		return pb.VideoStatus_READY
	case entities.StatusError:
		return pb.VideoStatus_ERROR
	default:
		return pb.VideoStatus_UPLOADED
	}
}

func fromProtoStatus(s pb.VideoStatus) entities.VideoStatus {
	switch s {
	case pb.VideoStatus_UPLOADED:
		return entities.StatusUploaded
	case pb.VideoStatus_TRANSCODING:
		return entities.StatusTranscoding
	case pb.VideoStatus_READY:
		return entities.StatusReady
	case pb.VideoStatus_ERROR:
		return entities.StatusError
	default:
		return entities.StatusUploaded
	}
}
