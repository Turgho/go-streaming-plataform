package domain

import "context"

// VideoInfo contém metadados extraídos do arquivo de vídeo.
type VideoInfo struct {
	Duration float64
	Size     int64
	Width    int
}

// Probe é a porta de inspeção de vídeo (implementada por ffprobe na infra).
type Probe interface {
	Inspect(ctx context.Context, filePath string) (*VideoInfo, error)
}
