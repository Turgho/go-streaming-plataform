package entities

import (
	"fmt"
	"strings"
	"time"
)

const (
	MaxVideoSize int64 = 5 * 1024 * 1024 * 1024 // 5GB
	MinVideoSize int64 = 1                      // 1 byte
)

type VideoStatus string

const (
	StatusUploaded    VideoStatus = "UPLOADED"
	StatusTranscoding VideoStatus = "TRANSCODING"
	StatusReady       VideoStatus = "READY"
	StatusError       VideoStatus = "ERROR"
)

type VideoMimetype string

const (
	MimetypeMP4       VideoMimetype = "video/mp4"
	MimetypeQuicktime VideoMimetype = "video/quicktime"
	MimetypeMKV       VideoMimetype = "video/x-matroska"
	MimetypeWebM      VideoMimetype = "video/webm"
	MimetypeAVI       VideoMimetype = "video/x-msvideo"
)

type Video struct {
	ID          string        `bson:"_id"`
	UserID      string        `bson:"user_id"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Mimetype    VideoMimetype `bson:"mimetype"`
	Size        int64         `bson:"size"`
	Duration    float64       `bson:"duration"`
	Status      VideoStatus   `bson:"status"`
	FilePath    string        `bson:"file_path"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

func NewVideo(
	id, userID, title string,
	description, mimetype, filepath string,
	size int64, duration float64,
) (*Video, error) {
	var errs []string

	if id == "" {
		errs = append(errs, "id é obrigatório")
	}
	if userID == "" {
		errs = append(errs, "userID é obrigatório")
	}
	if title == "" {
		errs = append(errs, "title é obrigatório")
	}
	if mimetype == "" {
		errs = append(errs, "mimetype é obrigatório")
	} else if !VideoMimetype(mimetype).IsValid() {
		errs = append(errs, "mimetype inválido")
	}
	if filepath == "" {
		errs = append(errs, "filepath é obrigatório")
	}
	if size <= MinVideoSize {
		errs = append(errs, "arquivo muito pequeno")
	} else if size > MaxVideoSize {
		errs = append(errs, "arquivo muito grande, máximo 5GB")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("video inválido: %s", strings.Join(errs, ", "))
	}

	now := time.Now()
	return &Video{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Mimetype:    VideoMimetype(mimetype),
		Size:        size,
		Duration:    duration,
		Status:      StatusUploaded,
		FilePath:    filepath,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (m VideoMimetype) IsValid() bool {
	switch m {
	case MimetypeMP4, MimetypeQuicktime, MimetypeMKV, MimetypeWebM, MimetypeAVI:
		return true
	}
	return false
}
