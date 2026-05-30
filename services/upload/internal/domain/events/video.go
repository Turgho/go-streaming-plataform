package events

const VideoUploaded = "video.uploaded"

type VideoUploadedEvent struct {
	VideoID  string `json:"video_id"`
	FilePath string `json:"file_path"`
	Mimetype string `json:"mimetype"`
}
