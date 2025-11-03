package postgres

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/yourorg/customer-mdm/internal/ports"
)

type TxManager struct {
    Pool *pgxpool.Pool
}

type tx struct {
    ctx context.Context
    tx  pgx.Tx
}

func (t *tx) Context() context.Context { return t.ctx }
func (t *tx) Commit() error            { return t.tx.Commit(t.ctx) }
func (t *tx) Rollback() error          { return t.tx.Rollback(t.ctx) }

func (m *TxManager) InTx(ctx context.Context, fn func(tx ports.Tx) error) error {
    pgxtx, err := m.Pool.Begin(ctx)
    if err != nil { return err }
    wrapper := &tx{ctx: ctx, tx: pgxtx}
    defer func() {
        // Rollback if not closed
        _ = pgxtx.Rollback(ctx)
    }()
    if err := fn(wrapper); err != nil {
        return err
    }
    return pgxtx.Commit(ctx)
}
