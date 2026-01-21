package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umbranian0/customer-mdm/internal/ports"
)

type Dispatcher struct {
	Pool      *pgxpool.Pool
	Publisher ports.EventPublisher
	BatchSize int
	PollEvery time.Duration
	Topic     string
}

type row struct {
	ID          string
	AggregateID string
	Payload     []byte
	EventType   string
	Headers     map[string]string
}

func (d *Dispatcher) Run(ctx context.Context) error {
	log.Printf("[outbox] dispatcher started, topic=%s batch=%d every=%s\n", d.Topic, d.BatchSize, d.PollEvery)
	ticker := time.NewTicker(d.PollEvery)
	defer ticker.Stop()
	for {
		if err := d.dispatchOnce(ctx); err != nil {
			log.Println("[outbox] dispatch error:", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (d *Dispatcher) dispatchOnce(ctx context.Context) error {
	if d.BatchSize <= 0 {
		d.BatchSize = 100
	}
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	rows, err := tx.Query(ctx, `SELECT id, aggregate_id, payload, event_type, headers FROM outbox_events WHERE published_at IS NULL ORDER BY occurred_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED`, d.BatchSize)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	var batch []row
	for rows.Next() {
		var r row
		var headersBytes []byte
		if err := rows.Scan(&r.ID, &r.AggregateID, &r.Payload, &r.EventType, &headersBytes); err != nil {
			rows.Close()
			_ = tx.Rollback(ctx)
			return err
		}
		if len(headersBytes) > 0 {
			if err := json.Unmarshal(headersBytes, &r.Headers); err != nil {
				rows.Close()
				_ = tx.Rollback(ctx)
				return err
			}
		}
		if r.Headers == nil {
			r.Headers = map[string]string{}
		}
		batch = append(batch, r)
	}
	rows.Close()

	for _, r := range batch {
		if err := d.Publisher.Publish(ctx, ports.Event{
			Topic:   d.Topic,
			Key:     []byte(r.AggregateID),
			Value:   r.Payload,
			Headers: r.Headers,
		}); err != nil {
			// increment attempts
			_, _ = tx.Exec(ctx, `UPDATE outbox_events SET attempts=attempts+1 WHERE id=$1`, r.ID)
			log.Printf("[outbox] publish failed id=%s err=%v\n", r.ID, err)
			continue
		}
		_, _ = tx.Exec(ctx, `UPDATE outbox_events SET published_at=NOW() WHERE id=$1`, r.ID)
		log.Printf("[outbox] published id=%s event_type=%s key=%s\n", r.ID, r.EventType, r.AggregateID)
	}
	return tx.Commit(ctx)
}
