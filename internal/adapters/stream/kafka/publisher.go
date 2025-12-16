package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

// Publisher is a thin wrapper around kafka-go that maintains per-topic writers.
type Publisher struct {
	brokers []string
	mu      sync.Mutex
	writers map[string]*kafka.Writer
}

func NewPublisher(brokers string) *Publisher {
	return &Publisher{
		brokers: strings.Split(brokers, ","),
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
	if err := ensureTopic(p.brokers, topic, 3); err != nil {
		log.Printf("[kafka] ensure topic %s failed: %v", topic, err)
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10_000_000, // 10ms
		RequiredAcks: kafka.RequireOne,
	}
	p.writers[topic] = w
	return w
}
