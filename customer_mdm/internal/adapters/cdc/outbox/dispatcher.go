package outbox

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/customer-mdm/internal/ports"
)

type Dispatcher struct {
    Pool      *pgxpool.Pool
    Publisher ports.EventPublisher
    BatchSize int
    PollEvery time.Duration
    Topic     string
}

type row struct {
    ID         string
    Payload    []byte
}

func (d *Dispatcher) Run(ctx context.Context) error {
    ticker := time.NewTicker(d.PollEvery)
    defer ticker.Stop()
    for {
        if err := d.dispatchOnce(ctx); err != nil {
            // TODO: log
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
        }
    }
}

func (d *Dispatcher) dispatchOnce(ctx context.Context) error {
    if d.BatchSize <= 0 { d.BatchSize = 100 }
    tx, err := d.Pool.Begin(ctx)
    if err != nil { return err }
    rows, err := tx.Query(ctx, `SELECT id, payload FROM outbox_events WHERE published_at IS NULL ORDER BY occurred_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED`, d.BatchSize)
    if err != nil {
        _ = tx.Rollback(ctx); return err
    }
    var batch []row
    for rows.Next() {
        var r row
        if err := rows.Scan(&r.ID, &r.Payload); err != nil {
            rows.Close(); _ = tx.Rollback(ctx); return err
        }
        batch = append(batch, r)
    }
    rows.Close()

    for _, r := range batch {
        if err := d.Publisher.Publish(ctx, ports.Event{
            Topic: d.Topic,
            Key:   nil,
            Value: r.Payload,
        }); err != nil {
            // increment attempts
            _, _ = tx.Exec(ctx, `UPDATE outbox_events SET attempts=attempts+1 WHERE id=$1`, r.ID)
            continue
        }
        _, _ = tx.Exec(ctx, `UPDATE outbox_events SET published_at=NOW() WHERE id=$1`, r.ID)
    }
    return tx.Commit(ctx)
}
