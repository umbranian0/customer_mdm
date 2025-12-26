package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Consumer wraps kafka-go's Reader so callers can focus on handling messages.
type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(client *Client, topic, groupID string) (*Consumer, error) {
	reader, err := client.NewReader(topic, groupID)
	if err != nil {
		return nil, err
	}
	return &Consumer{reader: reader}, nil
}

func (c *Consumer) Fetch(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
