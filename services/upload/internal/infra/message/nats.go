package message

import (
	"encoding/json"
	"fmt"

	"platform/events"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	return &NatsClient{conn: conn}, nil
}

func (n *NatsClient) Publish(subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return n.conn.Publish(subject, data)
}

// PublishUploaded implementa usecase.VideoPublisher.
func (n *NatsClient) PublishUploaded(event events.VideoUploadedEvent) error {
	return n.Publish(events.VideoUploaded, event)
}

func (n *NatsClient) Close() {
	n.conn.Close()
}
