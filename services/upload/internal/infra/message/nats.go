package message

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}
	return &NatsClient{conn: conn}, nil
}

func (n *NatsClient) Publish(subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	return n.conn.Publish(subject, data)
}

func (n *NatsClient) Close() {
	n.conn.Close()
}
