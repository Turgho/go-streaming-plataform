package repositories

import (
	"context"
	"upload-service/internal/domain/entities"
)

type VideoRepository interface {
	Create(ctx context.Context, video *entities.Video) error
	GetByID(ctx context.Context, id string) (*entities.Video, error)
	ListByUserID(ctx context.Context, userID string) ([]*entities.Video, error)
	UpdateStatus(ctx context.Context, id string, status entities.VideoStatus) error
	UpdateMetadata(ctx context.Context, id string, duration float64, size int64) error
}
