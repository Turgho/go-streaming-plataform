package events

// Assunto NATS publicado após upload concluído com sucesso.
const VideoUploaded = "video.uploaded"

// VideoUploadedEvent é o contrato compartilhado entre upload e transcode.
type VideoUploadedEvent struct {
	VideoID  string `json:"video_id"`
	FilePath string `json:"file_path"`
	Mimetype string `json:"mimetype"`
	Width    int    `json:"width"`
}
