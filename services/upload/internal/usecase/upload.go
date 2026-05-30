package usecase

import (
	"context"
	"fmt"

	"platform/events"
	"upload-service/internal/domain/entities"
	domain "upload-service/internal/domain/probe"
	domainrepo "upload-service/internal/domain/repositories"
)

type UploadUseCase struct {
	repo      domainrepo.VideoRepository
	probe     domain.Probe
	publisher VideoPublisher
}

func NewUploadUseCase(
	repo domainrepo.VideoRepository,
	probe domain.Probe,
	publisher VideoPublisher,
) *UploadUseCase {
	return &UploadUseCase{
		repo:      repo,
		probe:     probe,
		publisher: publisher,
	}
}

type CompleteUploadInput struct {
	VideoID     string
	UserID      string
	Title       string
	Description string
	Mimetype    string
	FilePath    string
}

// CompleteUpload valida o arquivo, persiste o vídeo e publica o evento de transcode.
func (uc *UploadUseCase) CompleteUpload(ctx context.Context, in CompleteUploadInput) (*entities.Video, error) {
	info, err := uc.probe.Inspect(ctx, in.FilePath)
	if err != nil {
		return nil, fmt.Errorf("inspeção do vídeo: %w", err)
	}

	video, err := entities.NewVideo(
		in.VideoID, in.UserID,
		in.Title, in.Description, in.Mimetype,
		in.FilePath, info.Size, info.Duration,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, video); err != nil {
		return nil, fmt.Errorf("persistir vídeo: %w", err)
	}

	event := events.VideoUploadedEvent{
		VideoID:  video.ID,
		FilePath: video.FilePath,
		Mimetype: string(video.Mimetype),
		Width:    info.Width,
	}

	if err := uc.publisher.PublishUploaded(event); err != nil {
		return nil, fmt.Errorf("publicar evento: %w", err)
	}

	return video, nil
}
