package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"upload-service/internal/domain/entities"
	domainrepo "upload-service/internal/domain/repositories"
)

type VideoUseCase struct {
	repo domainrepo.VideoRepository
}

func NewVideoUseCase(repo domainrepo.VideoRepository) *VideoUseCase {
	return &VideoUseCase{repo: repo}
}

func (uc *VideoUseCase) GetByID(ctx context.Context, id, requesterUserID string) (*entities.Video, error) {
	video, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if video.UserID != requesterUserID {
		return nil, ErrForbidden
	}

	return video, nil
}

func (uc *VideoUseCase) ListByUserID(ctx context.Context, userID string) ([]*entities.Video, error) {
	return uc.repo.ListByUserID(ctx, userID)
}

func (uc *VideoUseCase) UpdateStatus(ctx context.Context, id string, status entities.VideoStatus) (*entities.Video, error) {
	if err := uc.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	video, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return video, nil
}

func (uc *VideoUseCase) UpdateMetadata(ctx context.Context, id string, duration float64, size int64) (*entities.Video, error) {
	if err := uc.repo.UpdateMetadata(ctx, id, duration, size); err != nil {
		return nil, err
	}

	video, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return video, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func WrapError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
