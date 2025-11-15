// The purpose of this pakcage is to connect to the nats server and send all data to the nats
package publisher

import (
	"Paprika/models"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn

	mu sync.Mutex
}

func NewPublisher(serverUrl string) (*Publisher, error) {
	conn, err := nats.Connect(serverUrl)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn}, nil
}

// Publish - send a provided message to the connected NATS server with a provided subject
//
// TODO: add protobuf
func (p *Publisher) Publish(subject string, data interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := p.conn.Publish(subject, payload); err != nil {
		return nil
	}

	return nil
}

// PublishTickers - because the main services, that will publish messages are exhanges 
// (and exchanges just fetch a lot of data every second), its more convenient to pass an array of messages
//
// TODO: add protobuf
func (p *Publisher) PublishTickers(subject string, data *[]models.Ticker) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, elem := range *data {
		payload, err := json.Marshal(elem)
		if err != nil {
			return fmt.Errorf("failed to marshal the element in the provided object with a subject '%s': ", subject, err)
		}
		if err := p.conn.Publish(subject, payload); err != nil {
			return fmt.Errorf("failed to publish message with a subject '%s': %v", subject, err)
		}
	}

	return nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}
