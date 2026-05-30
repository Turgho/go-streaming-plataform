package handler

import (
	"context"
	"log"

	"platform/events"
	"transcode-service/internal/infra/message"
	"transcode-service/internal/usecase"
)

// Worker escuta eventos NATS e delega o processamento ao use case.
type Worker struct {
	nats    *message.NatsClient
	usecase *usecase.TranscodeUseCase
}

func NewWorker(nats *message.NatsClient, uc *usecase.TranscodeUseCase) *Worker {
	return &Worker{nats: nats, usecase: uc}
}

func (w *Worker) Start() error {
	log.Println("worker iniciado, escutando", events.VideoUploaded)

	return w.nats.Subscribe(events.VideoUploaded, func(event events.VideoUploadedEvent) {
		w.usecase.ProcessUploaded(context.Background(), event)
	})
}
