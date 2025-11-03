package ports

import "context"

type Tx interface {
    Context() context.Context
    Commit() error
    Rollback() error
}

type TxManager interface {
    InTx(ctx context.Context, fn func(tx Tx) error) error
}
