package grpc

import (
	"context"
	"errors"
	"io"
	"strings"

	"upload-service/internal/infra/interceptor"
	"upload-service/internal/usecase"
	"upload-service/pkg/pb"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUploadServiceServer
	upload *usecase.UploadUseCase
	video  *usecase.VideoUseCase
}

func NewServer(upload *usecase.UploadUseCase, video *usecase.VideoUseCase) *Server {
	return &Server{
		upload: upload,
		video:  video,
	}
}

func (s *Server) UploadVideo(stream pb.UploadService_UploadVideoServer) error {
	first, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "receber metadata: %v", err)
	}

	meta := first.GetMetadata()
	if meta == nil {
		return status.Error(codes.InvalidArgument, "primeira mensagem deve conter metadata")
	}

	userID, ok := stream.Context().Value(interceptor.UserIDKey).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "usuário não autenticado")
	}

	videoID, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 21)
	if err != nil {
		return status.Error(codes.Internal, "gerar id do vídeo")
	}

	filePath, err := saveStreamToDisk(videoID, meta, func() (*pb.UploadVideoRequest, error) {
		msg, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil, io.EOF
		}
		return msg, recvErr
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return status.Error(codes.Internal, "stream vazio")
		}
		if strings.Contains(err.Error(), "mimetype inválido") {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		return status.Errorf(codes.Internal, "salvar arquivo: %v", err)
	}

	video, err := s.upload.CompleteUpload(stream.Context(), usecase.CompleteUploadInput{
		VideoID:     videoID,
		UserID:      userID,
		Title:       meta.Title,
		Description: meta.Description,
		Mimetype:    meta.Mimetype,
		FilePath:    filePath,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "%v", err)
	}

	return stream.SendAndClose(&pb.UploadVideoResponse{Video: toProtoVideo(video)})
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.VideoMetadata, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id é obrigatório")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	video, err := s.video.GetByID(ctx, req.Id, userID)
	if err != nil {
		return nil, mapUsecaseError(err)
	}

	return toProtoVideo(video), nil
}

func (s *Server) ListByUserID(ctx context.Context, req *pb.ListByUserIDRequest) (*pb.ListByUserIDResponse, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Lista apenas os vídeos do usuário autenticado (ignora user_id arbitrário do request).
	videos, err := s.video.ListByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoVideos := make([]*pb.VideoMetadata, 0, len(videos))
	for _, v := range videos {
		protoVideos = append(protoVideos, toProtoVideo(v))
	}

	return &pb.ListByUserIDResponse{Videos: protoVideos}, nil
}

// UpdateStatus é chamado pelo transcode-service (autenticação via x-service-key no interceptor).
func (s *Server) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.VideoMetadata, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id é obrigatório")
	}
	if _, ok := pb.VideoStatus_name[int32(req.Status)]; !ok {
		return nil, status.Error(codes.InvalidArgument, "status inválido")
	}

	video, err := s.video.UpdateStatus(ctx, req.Id, fromProtoStatus(req.Status))
	if err != nil {
		return nil, mapUsecaseError(err)
	}

	return toProtoVideo(video), nil
}

func (s *Server) UpdateMetadata(ctx context.Context, req *pb.UpdateMetadataRequest) (*pb.VideoMetadata, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id é obrigatório")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.video.GetByID(ctx, req.Id, userID)
	if err != nil {
		return nil, mapUsecaseError(err)
	}

	video, err := s.video.UpdateMetadata(ctx, existing.ID, req.Duration, req.Size)
	if err != nil {
		return nil, mapUsecaseError(err)
	}

	return toProtoVideo(video), nil
}

func userIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(interceptor.UserIDKey).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "usuário não autenticado")
	}
	return userID, nil
}

func mapUsecaseError(err error) error {
	switch {
	case usecase.IsNotFound(err):
		return status.Error(codes.NotFound, err.Error())
	case usecase.IsForbidden(err):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
