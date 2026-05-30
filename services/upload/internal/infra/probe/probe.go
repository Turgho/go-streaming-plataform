package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	domain "upload-service/internal/domain/probe"
)

var ffprobeExec = func(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "ffprobe", args...)
}

type FFmpegProbe struct{}

func New() *FFmpegProbe {
	return &FFmpegProbe{}
}

// Inspect implementa domain.Probe usando ffprobe.
func (p *FFmpegProbe) Inspect(ctx context.Context, filePath string) (*domain.VideoInfo, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-show_format",
		filePath,
	}

	cmd := ffprobeExec(ctx, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var result struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int    `json:"width"`
		} `json:"streams"`
		Format struct {
			Duration string `json:"duration"`
			Size     string `json:"size"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse ffprobe output: %w", err)
	}

	if len(result.Streams) == 0 {
		return nil, fmt.Errorf("no streams found in file")
	}

	var duration float64
	fmt.Sscanf(result.Format.Duration, "%f", &duration)

	var size int64
	fmt.Sscanf(result.Format.Size, "%d", &size)

	var width int
	for _, s := range result.Streams {
		if s.CodecType == "video" {
			width = s.Width
			break
		}
	}

	return &domain.VideoInfo{
		Duration: duration,
		Size:     size,
		Width:    width,
	}, nil
}
