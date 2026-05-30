package message

import (
	"encoding/json"
	"fmt"
	"log"

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

func (n *NatsClient) Subscribe(subject string, handler func(event events.VideoUploadedEvent)) error {
	_, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		var event events.VideoUploadedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("unmarshal evento: %v", err)
			return
		}
		handler(event)
	})
	return err
}

func (n *NatsClient) Close() {
	n.conn.Close()
}
