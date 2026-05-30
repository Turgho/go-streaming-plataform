package transcode

import "context"

type Transcoder interface {
	Transcode(ctx context.Context, videoID, inputPath string, resolutions []string) error
}
