package kafka

import (
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Client centralizes common Kafka setup (brokers, dialer) so publishers and
// consumers can share the same connection settings.
type Client struct {
	brokers []string
	dialer  *kafka.Dialer
}

func NewClient(brokers string) *Client {
	return &Client{
		brokers: splitAndTrim(brokers),
		dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	}
}

func (c *Client) Brokers() []string {
	return append([]string(nil), c.brokers...)
}

func (c *Client) EnsureTopic(topic string, partitions int) error {
	return ensureTopic(c.brokers, topic, partitions)
}

func (c *Client) NewWriter(topic string) (*kafka.Writer, error) {
	if err := c.EnsureTopic(topic, 3); err != nil {
		log.Printf("[kafka] ensure topic %s skipped: %v", topic, err)
	}
	return &kafka.Writer{
		Addr:         kafka.TCP(c.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}, nil
}

func (c *Client) NewReader(topic, groupID string) (*kafka.Reader, error) {
	if err := c.EnsureTopic(topic, 1); err != nil {
		return nil, err
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		GroupID:  groupID,
		Topic:    topic,
		Dialer:   c.dialer,
		MinBytes: 1e4,
		MaxBytes: 10e6,
	}), nil
}

func splitAndTrim(brokers string) []string {
	parts := strings.Split(brokers, ",")
	var res []string
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			res = append(res, v)
		}
	}
	return res
}
