package usecase

import "platform/events"

// VideoPublisher publica eventos de vídeo (implementado pelo NATS na infra).
type VideoPublisher interface {
	PublishUploaded(event events.VideoUploadedEvent) error
}
