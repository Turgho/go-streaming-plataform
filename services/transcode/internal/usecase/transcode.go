package usecase

import (
	"context"
	"log"

	"platform/events"
	domain "transcode-service/internal/domain/transcode"
	uploadpb "transcode-service/pkg/uploadpb"
)

type UploadClient interface {
	UpdateStatus(ctx context.Context, videoID string, status uploadpb.VideoStatus) error
}

type TranscodeUseCase struct {
	transcoder domain.Transcoder
	upload     UploadClient
}

func NewTranscodeUseCase(transcoder domain.Transcoder, upload UploadClient) *TranscodeUseCase {
	return &TranscodeUseCase{
		transcoder: transcoder,
		upload:     upload,
	}
}

func (uc *TranscodeUseCase) ProcessUploaded(ctx context.Context, event events.VideoUploadedEvent) {
	log.Printf("processando vídeo: %s", event.VideoID)

	resolutions := buildResolutions(event.Width)

	if err := uc.upload.UpdateStatus(ctx, event.VideoID, uploadpb.VideoStatus_TRANSCODING); err != nil {
		log.Printf("atualizar status TRANSCODING (%s): %v", event.VideoID, err)
	}

	if err := uc.transcoder.Transcode(ctx, event.VideoID, event.FilePath, resolutions); err != nil {
		log.Printf("transcodar %s: %v", event.VideoID, err)
		if updateErr := uc.upload.UpdateStatus(ctx, event.VideoID, uploadpb.VideoStatus_ERROR); updateErr != nil {
			log.Printf("atualizar status ERROR (%s): %v", event.VideoID, updateErr)
		}
		return
	}

	if err := uc.upload.UpdateStatus(ctx, event.VideoID, uploadpb.VideoStatus_READY); err != nil {
		log.Printf("atualizar status READY (%s): %v", event.VideoID, err)
		return
	}

	log.Printf("vídeo %s transcodificado com sucesso", event.VideoID)
}

func buildResolutions(width int) []string {
	switch {
	case width >= 3840:
		return []string{"4K", "1440p", "1080p", "720p", "480p", "360p", "240p"}

	case width >= 2560:
		return []string{"1440p", "1080p", "720p", "480p", "360p", "240p"}

	case width >= 1920:
		return []string{"1080p", "720p", "480p", "360p", "240p"}

	case width >= 1280:
		return []string{"720p", "480p", "360p", "240p"}

	case width >= 854:
		return []string{"480p", "360p", "240p"}

	case width >= 640:
		return []string{"360p", "240p"}

	default:
		return []string{"240p"}
	}
}
