package grpc

import (
	"upload-service/internal/domain/entities"
	"upload-service/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
