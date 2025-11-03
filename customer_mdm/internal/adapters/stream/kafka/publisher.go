package kafka

import (
    "context"
    "github.com/umbranian0/customer-mdm/internal/ports"
    "github.com/segmentio/kafka-go"
)

type Publisher struct {
    p *ck.Producer
}

func NewPublisher(brokers string, extra map[string]string) (*Publisher, error) {
    cfg := &ck.ConfigMap{"bootstrap.servers": brokers}
    for k, v in (extra) {
        (*cfg)[k] = v
    }
    p, err := ck.NewProducer(cfg)
    if err != nil { return nil, err }
    return &Publisher{p: p}, nil
}

func (k *Publisher) Publish(ctx context.Context, ev ports.Event) error {
    msg := &ck.Message{
        TopicPartition: ck.TopicPartition{Topic: &ev.Topic, Partition: ck.PartitionAny},
        Key:            ev.Key,
        Value:          ev.Value,
    }
    return k.p.Produce(msg, nil)
}

func (k *Publisher) Close() { k.p.Close() }
