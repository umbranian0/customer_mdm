package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

// Publisher is a thin wrapper around kafka-go that maintains per-topic writers.
type Publisher struct {
	client  *Client
	mu      sync.Mutex
	writers map[string]*kafka.Writer
}

func NewPublisher(client *Client) *Publisher {
	return &Publisher{
		client:  client,
		writers: make(map[string]*kafka.Writer),
	}
}

func (p *Publisher) Publish(ctx context.Context, ev ports.Event) error {
	w := p.getWriter(ev.Topic)
	if w == nil {
		return fmt.Errorf("no writer for topic %s", ev.Topic)
	}
	msg := kafka.Message{
		Key:   ev.Key,
		Value: ev.Value,
	}
	if len(ev.Headers) > 0 {
		msg.Headers = make([]kafka.Header, 0, len(ev.Headers))
		for k, v := range ev.Headers {
			msg.Headers = append(msg.Headers, kafka.Header{Key: k, Value: []byte(v)})
		}
	}
	return w.WriteMessages(ctx, msg)
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, w := range p.writers {
		_ = w.Close()
	}
}

func (p *Publisher) getWriter(topic string) *kafka.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()
	if w, ok := p.writers[topic]; ok {
		return w
	}
	w, err := p.client.NewWriter(topic)
	if err != nil {
		log.Printf("[kafka] ensure writer for topic %s failed: %v", topic, err)
		return nil
	}
	p.writers[topic] = w
	return w
}
