package ports

import "context"

type Event struct {
    Topic   string
    Key     []byte
    Value   []byte
    Headers map[string]string
}

type EventPublisher interface {
    Publish(ctx context.Context, ev Event) error
}
